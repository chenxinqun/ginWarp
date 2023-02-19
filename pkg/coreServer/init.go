package coreServer

import (
	"context"
	"errors"
	"fmt"
	"github.com/chenxinqun/ginWarpPkg/datax/clickhousex"
	"github.com/chenxinqun/ginWarpPkg/errno"
	"reflect"

	"github.com/chenxinqun/ginWarp/configs"

	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	"github.com/chenxinqun/ginWarpPkg/datax/kafkax"
	"github.com/chenxinqun/ginWarpPkg/datax/mongox"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/datax/redisx"
	"github.com/chenxinqun/ginWarpPkg/idGen"
	"github.com/chenxinqun/ginWarpPkg/sysx/environment"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// 注册进来的自定义初始化函数, 是不允许修改配置项的, 因此传给他的配置都是只读的
type InitFunc func(cfg configs.Config) error

var initBeforeHandlers = make([]InitFunc, 0)
var initAfterHandlers = make([]InitFunc, 0)

var initHandlers = make([]InitFunc, 0)

func RegisterInitFunc(initFunc InitFunc) {
	initHandlers = append(initHandlers, initFunc)
}

func init() {
	// 按次序依次初始化

	// 初始化 MySQL 如果相关配置为空, 默认不进行初始化
	initBeforeHandlers = append(initBeforeHandlers, InitMySQLBefore)
	// 初始化 Clickhouse 如果相关配置为空, 默认不进行初始化
	initBeforeHandlers = append(initBeforeHandlers, InitClickhouseBefore)
	// 初始化MongoDB 如果相关配置为空, 默认不进行初始化
	initBeforeHandlers = append(initBeforeHandlers, InitMongoBefore)
	//初始化 RedisRepo 如果相关配置为空, 默认不进行初始化
	initBeforeHandlers = append(initBeforeHandlers, InitRedisBefore)
	// 初始化卡夫卡连接
	initBeforeHandlers = append(initBeforeHandlers, InitKafkaBefore)
	// 初始化ID生成器
	initBeforeHandlers = append(initBeforeHandlers, InitidGenBefore)
}

func init() {
	// 按次序依次初始化

	// 注册自己服务到注册中心
	initAfterHandlers = append(initAfterHandlers, InitRegisterAfter)
	// 从注册中心发现别的服务
	initAfterHandlers = append(initAfterHandlers, InitDiscoverAfter)
}

func InitHandlerExec(cfg configs.Config, logger *zap.Logger, handlers ...InitFunc) {
	for _, f := range handlers {
		err := f(cfg)
		if err != nil {
			logger.Fatal("初始化失败", zap.Error(err))
		}
	}
}

func InitMySQLBefore(cfg configs.Config) (err error) {
	if cfg.MySQL.Write != nil && len(cfg.MySQL.Write) > 0 && cfg.MySQL.Write[0].Addr != "" {
		// 初始化etcd连接
		_, err = mysqlx.New(cfg.MySQL)
		if err != nil {
			return errno.Errorf("new mysql err: %v", err)
		}
	}
	return nil
}

func InitClickhouseBefore(cfg configs.Config) (err error) {
	if cfg.Clickhouse.Addr != "" {
		// 初始化etcd连接
		_, err = clickhousex.New(cfg.Clickhouse)
		if err != nil {
			return errno.Errorf("new clickhouse err: %v", err)
		}
	}
	return nil
}

func InitRedisBefore(cfg configs.Config) (err error) {
	if cfg.Redis.Addr != "" || (len(cfg.Redis.Sentinel.SentinelAddrs) > 0 && cfg.Redis.Sentinel.SentinelAddrs[0] != "") {
		// 初始化etcd连接
		_, err = redisx.New(cfg.Redis)
		if err != nil {
			return errno.Errorf("new redis err: %v", err)
		}
	}
	return nil
}

func InitMongoBefore(cfg configs.Config) (err error) {
	if cfg.Mongo.Addrs != nil && len(cfg.Mongo.Addrs) > 0 && cfg.Mongo.Addrs[0] != "" {
		// 初始化etcd连接
		_, err = mongox.New(cfg.Mongo)
		if err != nil {
			return errno.Errorf("new mongo err: %v", err)
		}
	}
	return nil
}

func InitKafkaBefore(cfg configs.Config) (err error) {
	if cfg.Kafka.BrokerList != nil && len(cfg.Kafka.BrokerList) > 0 && cfg.Kafka.BrokerList[0] != "" {
		// 初始化卡夫卡连接
		_, err = kafkax.NewConsumerGroup(cfg.Kafka)
		if err != nil {
			return errno.Errorf("new kafka ConsumerGroup err: %v", err)
		}
		_, err = kafkax.NewProducer(cfg.Kafka)
		if err != nil && !errors.Is(err, kafkax.ProducerEnableErr) {
			return errno.Errorf("new kafka producer err: %v", err)
		}
		_, err = kafkax.NewAsyncProducer(cfg.Kafka)
		if err != nil && !errors.Is(err, kafkax.AsyncProducerEnableErr) {
			return errno.Errorf("new kafka AsyncProducer err: %v", err)
		}
	}
	return nil
}

func InitRegisterAfter(cfg configs.Config) (err error) {
	if !cfg.MicroserviceModel {
		return nil
	}
	etcdRepo := etcdx.Default()
	// 服务注册
	err = etcdRepo.Register(&cfg.Register)
	if err != nil {
		return errno.Errorf("register service http err: %v", err)
	}
	// grpc服务注册
	if !reflect.DeepEqual(cfg.GrpcRegister, etcdx.ServiceInfo{}) && cfg.GrpcRegister.Scheme == "grpc" {
		err = etcdRepo.Register(&cfg.GrpcRegister)
		if err != nil {
			return errno.Errorf("register service grpc err: %v", err)
		}
	}
	return nil
}

func InitDiscoverAfter(cfg configs.Config) (err error) {
	if !cfg.MicroserviceModel {
		return nil
	}
	etcdRepo := etcdx.Default()
	// 服务发现, 获取服务列表
	for _, info := range cfg.Discover {
		err = etcdRepo.Discover(info)
		if err != nil {
			return errno.Errorf("discover service err: %v", err)
		}
	}
	return nil
}

// InitidGenBefore 从etcd读取配置, 初始化ID生成器
func InitidGenBefore(cfg configs.Config) (err error) {
	DID := cfg.DatacenterID
	WID := cfg.WorkerID
	if cfg.MicroserviceModel {
		etcdRepo := etcdx.Default()
		breakT := false
		maxDID := idGen.GetDataCenterIDMax()
		maxWID := idGen.GetWorkerIDMax()
		value := etcdRepo.GetRepo().Service.Key
		keyPrefix := fmt.Sprintf("/DataCenterAndWorkerID/%s", cfg.Env)
		conn := etcdRepo.GetConn()
		ttl := int(etcdRepo.GetRepo().Service.Val.TTL)
		ctx, cancel := etcdRepo.TimeOutCtx(ttl)
		resp, _ := conn.Grant(ctx, int64(ttl))
		cancel()
		for i := int64(1); i < maxDID+1; i++ {
			keyP := keyPrefix + fmt.Sprintf("/%v", i)
			for j := int64(1); j < maxWID+1; j++ {
				key := keyP + fmt.Sprintf("/%v", j)
				kvc := clientv3.NewKV(conn)
				// 创建修改 = 0, 即 key 不存在
				compare := clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
				// 则按租约创建一个 key .
				opPut := clientv3.OpPut(key, value, clientv3.WithLease(resp.ID))
				ctx, cancel = etcdRepo.TimeOutCtx(ttl)
				// 开启一个事务. (并发时只有一个能写成功).
				ret, err := kvc.Txn(ctx).If(compare).Then(opPut).Commit()
				cancel()
				// 事务不成功则跳过
				if err != nil || !ret.Succeeded {
					continue
				}

				// 事务成功则跳出循环
				DID = i
				WID = j
				breakT = true
				break
			}
			if breakT {
				keepAlive, _ := conn.KeepAlive(context.Background(), resp.ID)
				go func() {
					for res := range keepAlive {
						etcdRepo.GetRepo().Logger.Debug("续约数据datacenterID与workerID成功", zap.Any("resp", res))
						continue
					}
				}()
				break
			}
		}
		cfg.DatacenterID = DID
		cfg.WorkerID = WID
	}
	if DID > idGen.GetDataCenterIDMax() {
		return errno.Errorf("idGen 的 DatacenterID 不能大于 %d", idGen.GetDataCenterIDMax())
	}
	if WID > idGen.GetWorkerIDMax() {
		return errno.Errorf("idGen 的 WorkerID 不能大于 %d", idGen.GetWorkerIDMax())
	}
	// 只有配置有效才会初始化
	if (DID > 0 && DID <= idGen.GetDataCenterIDMax()) && (WID > 0 && WID <= idGen.GetWorkerIDMax()) {
		_, _ = idGen.NewSnowflake(DID, WID)
	}
	return nil
}

func InitServerManual(env, serviceName, addr, buildDate, buildVersion string, closeFileLog bool) {
	// 初始化环境变量
	environment.InitEnv(env)
	confHandles := make([]configs.OptionHandler, 0)
	if !closeFileLog {
		confHandles = append(confHandles, configs.WithFile)
	}
	// 初始化配置文件
	configs.InitConf(addr, serviceName, environment.Active().Value(), buildDate, buildVersion, confHandles...)
	// 从配置文件, 读取业务平台信息
	environment.InitBusinessPlatform(configs.Default().BusinessPlatform)
}

func InitEtcdManual(cfg configs.Config, logger *zap.Logger) (err error) {
	if !cfg.MicroserviceModel {
		logger.Info("没有开启微服务模式, 跳过ETCD初始化")
		return nil
	}
	if cfg.ETCD.Endpoints != nil && len(cfg.ETCD.Endpoints) > 0 && cfg.ETCD.Endpoints[0] != "" {
		// 初始化etcd连接
		_, err = etcdx.New(cfg.ETCD, logger)
		if err != nil {
			return errno.Errorf("new etcd err: %v", err)
		}
	}
	return nil
}

func InitConfiguringManual(cfg configs.Config, etcdRepo etcdx.Repo) (err error) {
	if !cfg.MicroserviceModel {
		return nil
	}
	// 从配置中心获取配置
	err = etcdRepo.Configuring(cfg.Configuring, configs.GetConfSet(), configs.GetConfWatcher())
	if err != nil {
		return errno.Errorf("init config from etcd err: %v", err)
	}
	return nil
}
