package testing

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/chenxinqun/ginWarp/configs"

	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.etcd.io/etcd/client/v3"
)

func TestDiscover(t *testing.T) {
	flag.Parse()
	// 初始化 logger
	loggers, _ := logger.NewJSONLogger(
		logger.WithConsole(true),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", "zh_example", "fat")),
		logger.WithTimeLayout("2006-01-02 15:04:05"),
		logger.WithLogFile("etcd_test.log", 100, 5, 30),
	)
	cfg := configs.Default()
	ser, err := etcd.New(cfg.ETCD, loggers)
	if err != nil {
		t.Fatal(err)
	}

	err = ser.Register(&cfg.Register)
	if err != nil {
		t.Fatal(err)
	}
	err = ser.Discover(cfg.Register)
	if err != nil {
		t.Fatal(err)
	}
	serviceList := ser.GetServiceListing()
	if len(serviceList) == 0 {
		t.Fatal("服务列表为空, 测试不通过", serviceList)
	}
	select {
	case <-time.After(11 * time.Second):
		_ = ser.Close()
	}
	t.Log("测试通过", fmt.Sprintf("%v", serviceList))
}

func TestEtcdQuery(t *testing.T) {
	flag.Parse()
	cfg := configs.Default()
	conf := clientv3.Config{Endpoints: cfg.ETCD.Endpoints, DialTimeout: time.Duration(cfg.ETCD.DialTimeout) * time.Second}
	cli, err := clientv3.New(conf)
	if err != nil {
		t.Fatal("连接etcd报错, 测试失败 ", err)
	}
	resp, err := cli.Get(context.Background(), "/config/mysql", clientv3.WithPrevKV())
	if err != nil {
		t.Fatal("查询出错, 测试失败", err)
	}
	t.Log("返回值", resp)
	if resp != nil {
		for _, kv := range resp.Kvs {
			t.Log(string(kv.Key), string(kv.Value))
		}
	}

}
