package session

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/chenxinqun/ginWarp/configs"
	"github.com/chenxinqun/ginWarpPkg/cryptox/token"
	"github.com/chenxinqun/ginWarpPkg/datax/redisx"
)

const (
	TokenField   = "token"
	DataField    = "data"
	PrefixField  = "prefix"
	SecretField  = "secret"
	ExpireField  = "expire"
	JWTField     = "JWT"
	AuthKeyField = "authKey"
)

func GetSessionKey(userId int64, prefixs ...string) string {
	var prefix string
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	return fmt.Sprintf("%sjwt:token:%d", prefix, userId)
}

type Session interface {
	Set(key string, val interface{})
	Get(key string) (val interface{})
	LoadData(j json.RawMessage) error
	CreateToken() (string, error)
	ParseToken() (*token.Claims, error)
	SetToken(token string) Session
	GetToken() string
	SetExpire(duration time.Duration) Session
	Renewal() bool
	Read() Session
	Save() error
	Remove() bool
}

func New(userId int64, userName string, r redisx.Repo) Session {
	s := new(session)
	s.userId = userId
	s.userName = userName
	s.ctx = context.Background()
	s.key = GetSessionKey(userId)
	s.r = r
	s.data = make(map[string]interface{})
	jwt := configs.Default().Handler.GetStringMap(JWTField)
	expire := time.Duration(30)
	secretStr := ""
	if jwt != nil {
		if e, ok := jwt[ExpireField]; ok {
			switch e.(type) {
			case float64:
				expire = time.Duration(e.(float64))
			case int:
				expire = time.Duration(e.(int))
			case int64:
				expire = time.Duration(e.(int64))
			case int32:
				expire = time.Duration(e.(int32))
			}

		}
		if sS, ok := jwt[SecretField]; ok {
			switch sS.(type) {
			case string:
				secretStr = sS.(string)
			}
		}
	}
	s.expire = time.Minute * expire

	s.jwt = token.New(secretStr)
	if pf, ok := jwt[PrefixField]; ok {
		switch pf.(type) {
		case string:
			s.prefix = pf.(string)
		}
	}

	return s
}

type session struct {
	mt       sync.Mutex
	ctx      context.Context
	jwt      token.Token
	userId   int64
	userName string
	prefix   string
	key      string
	data     map[string]interface{}
	token    string
	expire   time.Duration
	r        redisx.Repo
}

func GetAuthKey() string {
	mp := configs.Default().Handler.GetStringMap(JWTField)
	return mp[AuthKeyField].(string)
}

func (s *session) CreateToken() (string, error) {
	tokenStr, err := s.jwt.JwtSign(s.userId, s.userName, s.expire)
	return tokenStr, err
}

func (s *session) GetToken() string {
	return s.token
}

// ParseToken 如果报错就是token无效或者已经过期
func (s *session) ParseToken() (*token.Claims, error) {
	ret, err := s.jwt.JwtParse(s.token)
	return ret, err
}

// LoadData 从json加载key值
func (s *session) LoadData(j json.RawMessage) error {
	s.mt.Lock()
	defer s.mt.Unlock()
	buf := make(map[string]interface{})
	err := json.Unmarshal(j, &buf)
	for key, val := range buf {
		s.data[key] = val
	}
	return err
}

// Set 设置某个Key值
func (s *session) Set(key string, val interface{}) {
	s.mt.Lock()
	defer s.mt.Unlock()
	s.data[key] = val
}

// Get 获取某个key值
func (s *session) Get(key string) (val interface{}) {
	s.mt.Lock()
	defer s.mt.Unlock()
	return s.data[key]
}

// SetExpire 设置过期时间
func (s *session) SetExpire(duration time.Duration) Session {
	s.expire = duration
	return s
}

// SetToken 设置token
func (s *session) SetToken(token string) Session {
	s.token = token
	return s
}

// Renewal token续期
func (s *session) Renewal() bool {
	return s.r.Expire(s.key, s.expire)
}

// Read 读取所有数据
func (s *session) Read() Session {
	ret, err := s.r.HGetAll(s.key).Result()
	if err != nil {
		return nil
	}
	err = json.Unmarshal([]byte(ret[DataField]), &s.data)
	if err != nil {
		return nil
	}
	s.token = ret[TokenField]
	return s
}

// Save 保存所有数据
func (s *session) Save() error {
	lockKey := s.key + ":lock:save"
	lock, err := s.r.SetNX(lockKey, 1, s.expire).Result()
	defer s.r.Del(lockKey)
	if err != nil {
		return err
	}
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	if lock {
		err = s.r.HSet(s.key, map[string]interface{}{TokenField: s.token, DataField: string(data)}).Err()
	}
	return err
}

// Remove 保存所有数据
func (s *session) Remove() bool {
	lockKey := s.key + ":lock:remove"
	lock, err := s.r.SetNX(lockKey, 1, s.expire).Result()
	defer s.r.Del(lockKey)
	if err != nil {
		return false
	}

	if lock {
		return s.r.Del(s.key)
	}
	return false
}
