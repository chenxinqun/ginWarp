package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	clientV3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"strings"
)

var (
	env  string
	addr string
	pla  string
	user string
	pwd  string
)

func init() {
	flag.StringVar(&env, "e", "pro", "目标环境: dev 开发, test 测试, pre 预发布, pro 生产")
	flag.StringVar(&pla, "pla", "", "业务平台")
	flag.StringVar(&addr, "t", "127.0.0.1:2379", "目标地址, 多个使用逗号分割")
	flag.StringVar(&user, "u", "", "账号")
	flag.StringVar(&pwd, "p", "", "密码")
	flag.Parse()
}

type ServiceConf struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

type Service struct {
	Name  string        `toml:"name"`
	Confs []ServiceConf `toml:"confs"`
}

type Cfg struct {
	Env              string
	BusinessPlatform string    // 业务平台
	ConfigPrefix     string    `toml:"configPrefix"`
	Services         []Service `toml:"services"`
}

func main() {
	cfgByte, err := ioutil.ReadFile(fmt.Sprintf("%sCfg.toml", env))
	if err != nil {
		log.Fatal("读取配置文件报错: ", err)
	}
	cfg := &Cfg{}
	err = toml.Unmarshal(cfgByte, cfg)
	if len(cfg.ConfigPrefix) == 0 {
		cfg.ConfigPrefix = "/Config"
	}
	if len(pla) > 0 {
		// 业务平台
		cfg.BusinessPlatform = pla
	}
	if len(cfg.BusinessPlatform) > 0 {
		cfg.ConfigPrefix = fmt.Sprintf("/%s%s", cfg.BusinessPlatform, cfg.ConfigPrefix)
	}
	cfg.Env = env
	if len(cfg.Env) == 0 {
		cfg.Env = "pro"
	}
	if err != nil {
		log.Fatal("解析配置文件内容报错: ", err)
	}
	endpoints := strings.Split(strings.TrimSpace(addr), ",")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	etcdClient, err := clientV3.New(clientV3.Config{Endpoints: endpoints, Username: user, Password: pwd, Context: ctx})
	fmt.Println("etcd 已连接", etcdClient.Endpoints())
	if err != nil {
		cancel()
		log.Fatal("etcd连接报错: ", err)
	}
	for _, service := range cfg.Services {
		k := fmt.Sprintf("%s/%s/%s", cfg.ConfigPrefix, cfg.Env, service.Name)
		log.Println("删除旧配置", k)
		_, e := etcdClient.Delete(ctx, k, clientV3.WithPrefix())
		if e != nil {
			log.Println("删除旧配置失败: ", k)
		}
		for _, conf := range service.Confs {
			key := fmt.Sprintf("%s/%s/%s/%s", cfg.ConfigPrefix, cfg.Env, service.Name, conf.Key)
			fmts := make([]interface{}, 0)
			for i := 0; i < strings.Count(conf.Value, `%s`); i++ {
				fmts = append(fmts, cfg.Env)
			}
			conf.Value = fmt.Sprintf(conf.Value, fmts...)
			log.Println("添加新配置", key, conf.Value)
			resp, e := etcdClient.Put(ctx, key, conf.Value)
			if e != nil {
				log.Println("添加新配置失败", *resp, key, conf.Value)
			} else {
				log.Println("添加新配置成功", *resp, key, conf.Value)
			}
		}
	}
}
