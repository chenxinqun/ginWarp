package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/chenxinqun/ginWarpPkg/convert"
	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	"github.com/chenxinqun/ginWarpPkg/identify"
	"github.com/chenxinqun/ginWarpPkg/sysx/environment"
	"github.com/spf13/cast"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	confWatcher etcdx.WatcherFunc
	confSet     etcdx.SetConfFunc
)

func init() {
	SetConfWatcher(ConfigWatcher)
	SetConfSet(ConfigSet)
}

func GetConfWatcher() etcdx.WatcherFunc {
	return confWatcher
}

func GetConfSet() etcdx.SetConfFunc {
	return confSet
}

func SetConfWatcher(f etcdx.WatcherFunc) { confWatcher = f }

func SetConfSet(f etcdx.SetConfFunc) { confSet = f }

//nolint:nolintlint,funlen
func ConfigMatching(er etcdx.Repo, confKey string, etcdValue []byte, addLog bool, initLog bool,
	confTypeOf reflect.Type, confValueOf reflect.Value) {
	key := confKey
	if addLog {
		addLog = true
		er.GetRepo().Logger.Info("配置新增", zap.String(key, string(etcdValue)))
	}
	var (
		oldVal interface{}
	)
	vTp, exist := confTypeOf.FieldByName(key)

	// 结构体中有的key这样处理
	if exist { //nolint:nolintlint,nestif
		oldVal = confValueOf.FieldByName(key).Interface()
		vTag := vTp.Tag.Get("json")
		if vTag != "" {
			val := confValueOf.FieldByName(key)

			switch val.Type().Kind() {
			// 如果是结构体, 这样处理
			case reflect.Struct:
				v := reflect.New(val.Type()).Interface()
				e := json.Unmarshal(etcdValue, &v)
				if e == nil {
					t := reflect.ValueOf(v)
					val.Set(t.Elem())
				} else {
					er.GetRepo().Logger.Error("通过配置中心修改配置, 序列化json时出错", zap.Error(e), zap.String(key, string(etcdValue)))
				}
			// 如果不是结构体, 这样处理
			default:
				sv := string(etcdValue)
				switch val.Interface().(type) {
				case []etcdx.ServiceInfo:
					v := make([]etcdx.ServiceInfo, 0)
					e := json.Unmarshal(etcdValue, &v)
					if e == nil {
						val.Set(reflect.ValueOf(v))
					} else {
						er.GetRepo().Logger.Error("通过配置中心修改配置, 序列化json时出错", zap.Error(e), zap.String(key, string(etcdValue)))
					}
				case bool:
					v := cast.ToBool(sv)
					val.Set(reflect.ValueOf(v))
				case string:
					v := cast.ToString(sv)
					val.Set(reflect.ValueOf(v))
				case int32, int16, int8, int:
					v := cast.ToInt(sv)
					val.Set(reflect.ValueOf(v))
				case uint:
					v := cast.ToUint(sv)
					val.Set(reflect.ValueOf(v))
				case uint32:
					v := cast.ToUint32(sv)
					val.Set(reflect.ValueOf(v))
				case uint64:
					v := cast.ToUint64(sv)
					val.Set(reflect.ValueOf(v))
				case int64:
					v := cast.ToInt64(sv)
					val.Set(reflect.ValueOf(v))
				case float64, float32:
					v := cast.ToFloat64(sv)
					val.Set(reflect.ValueOf(v))
				case time.Time:
					v := cast.ToTime(sv)
					val.Set(reflect.ValueOf(v))
				case time.Duration:
					v := cast.ToDuration(sv)
					val.Set(reflect.ValueOf(v))
				case []string:
					v := make([]string, 0)
					e := json.Unmarshal(etcdValue, &v)
					if e != nil {
						er.GetRepo().Logger.Error(fmt.Sprintf("通过配置中心修改配置, 序列化json时出错, err: %s, key: %s, value: %s", e, key, etcdValue))
					}
					val.Set(reflect.ValueOf(v))
				case []int:
					v := make([]int, 0)
					e := json.Unmarshal(etcdValue, &v)
					if e != nil {
						er.GetRepo().Logger.Error(fmt.Sprintf("通过配置中心修改配置, 序列化json时出错, err: %s, key: %s, value: %s", e, key, etcdValue))
					}
					val.Set(reflect.ValueOf(v))
				}
			}
		}
	} else {
		// 结构体重没有的key, 这样处理
		oldVal = Default().Handler.Kvs()[key]
		er.SetConfig(key, etcdValue, Default().Handler.Kvs())
	}
	if initLog {
		er.GetRepo().Logger.Info("从配置中心重新初始化配置", zap.String("key", key), zap.Any("配置文件内容", oldVal), zap.String("配置中心内容", string(etcdValue)))
	} else {
		if !addLog {
			// 改动配置
			er.GetRepo().Logger.Info("配置改动", zap.String("key", key), zap.Any("原配置", oldVal), zap.String("新配置", string(etcdValue)))
		}
	}
}

// ConfigWatcher 监听配置中心配置变动的函数, 可以根据实际情况, 调用SetConfWatcher覆写这个函数
func ConfigWatcher(er etcdx.Repo, wresp clientv3.WatchResponse) {
	confTypeOf := reflect.TypeOf(Default())
	confValueOf := reflect.ValueOf(writeDefault()).Elem()
	for _, ev := range wresp.Events {
		key := strings.TrimPrefix(string(ev.Kv.Key), er.GetRepo().ConfigInfo.Prefix)
		switch ev.Type {
		case mvccpb.PUT: //修改或者新增
			ConfigMatching(er, key, ev.Kv.Value, ev.Kv.Version == 1, false, confTypeOf, confValueOf)
		}
	}
}

func GetFileName(fileName string, flag int) string {
	suffix := ".json"
	if _, exist := identify.IsExists(fileName + suffix); exist {
		flag += 1
		return GetFileName(fmt.Sprintf("%s%d", fileName, flag), flag)
	}
	return fileName + suffix
}

func BakEtcdConf(etcdBak map[string]interface{}) error {
	jsonBak, err := json.Marshal(etcdBak)
	timeLayout := time.Now().Format("20060102150405")
	fileName := DefaultBakPath + "/" + DefaultEtcdBakName + timeLayout
	_ = os.MkdirAll(DefaultBakPath, 0755)
	fileName = GetFileName(fileName, 0)
	f, err := os.Create(fileName)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err == nil {
		_, _ = f.Write(jsonBak)
	}
	return err
}

// ConfigSet 从配置中心初始化配置的函数, 可以根据实际情况, 调用SetConfSet覆写这个函数
//nolint:nolintlint,gocognit
func ConfigSet(er etcdx.Repo, resp clientv3.GetResponse) {
	confTypeOf := reflect.TypeOf(Default())
	confValueOf := reflect.ValueOf(writeDefault()).Elem()
	etcdMap := make(map[string]string)
	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), er.GetRepo().ConfigInfo.Prefix)
		etcdMap[key] = string(kv.Key)
		ConfigMatching(er, key, kv.Value, false, true, confTypeOf, confValueOf)
	}
	configList, _ := convert.StructToMap(Default())
	etcdBak := make(map[string]interface{})
	for k, v := range configList {
		// 判断是否存在于etcd中, 如果不存在, 则往etcd中写配置
		_, exist := etcdMap[k]
		vTp, _ := confTypeOf.FieldByName(k)
		vtag := vTp.Tag.Get("json")
		key := er.GetRepo().ConfigInfo.Prefix + k
		if !exist && vtag != "" {
			var (
				err     error
				jsonVal []byte
				etcdVal string
			)
			king := vTp.Type.Kind()
			// 不是结构体, 数组, 切片 和 map的话, 就直接用字符串格式化. 否则就用json序列化.
			if king != reflect.Struct && king != reflect.Array &&
				king != reflect.Slice && king != reflect.Map {
				etcdVal = fmt.Sprintf("%v", v)
			} else {
				jsonVal, err = json.Marshal(v)
				etcdVal = string(jsonVal)
			}
			if err == nil && len(etcdVal) > 0 {
				timeout := time.Second * time.Duration(er.GetRepo().ConfigInfo.TTL)
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				_, e := er.GetConn().KV.Put(ctx, key, etcdVal)
				cancel()
				if e != nil {
					er.GetRepo().Logger.Info("从本地导出配置到etcd时错误", zap.Error(e), zap.String("config key", k),
						zap.String("etcd key", key), zap.String("value", string(jsonVal)))
					continue
				}
				er.GetRepo().Logger.Info("从本地导出配置到etcd", zap.String("config key", k),
					zap.String("etcd key", key), zap.String("value", string(jsonVal)))
			}
		}
		if vtag != "" {
			etcdBak[key] = v
		}
	}
	for k, v := range Default().Handler.Kvs() {
		key := er.GetRepo().ConfigInfo.Prefix + k
		etcdBak[key] = v
	}
	// 如果不是生产环境和预发布环境的话
	if !environment.Active().IsPro() && !environment.Active().IsPre() {
		e := BakEtcdConf(etcdBak)
		if e != nil {
			er.GetRepo().Logger.Error("保存etcd配置文件副本时错误", zap.Error(e))
		}
	}
}
