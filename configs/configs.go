package configs

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/chenxinqun/ginWarpPkg/datax/clickhousex"
	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	"github.com/chenxinqun/ginWarpPkg/datax/kafkax"
	"github.com/chenxinqun/ginWarpPkg/datax/mongox"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/datax/redisx"
	"github.com/chenxinqun/ginWarpPkg/ipTools"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"github.com/chenxinqun/ginWarpPkg/sysx/environment"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var configObj = new(Config)

// writeDefault 这个配置对象返回指针, 可写, 只允许内部使用
func writeDefault() *Config {
	return configObj
}

// Default 这个配置对象返回的是copy的变量, 修改不会改变真是配置, 请用在只读场合
func Default() Config {
	return *configObj
}

type BuildInfo struct {
	BuildDate    string // 编译时,由CI/CD工具复写此变量, 并在main函数中传进来
	BuildVersion string // 编译时,由CI/CD工具复写此变量, 并在main函数中传进来
}

type Config struct {
	ConfFile          string // 配置文件
	Env               string // 环境变量
	DatacenterID      int64  // ID生成器所用到的数据中心ID
	WorkerID          int64  // ID生成器所用到的微服务进程ID
	MicroserviceModel bool   // 微服务模式, 如果为True 则开启ETCD配置注册中心
	ServiceName       string
	BuildInfo         BuildInfo
	BusinessPlatform  string `toml:"BusinessPlatform"` // 业务平台
	// 获取k8s域名
	K8sDomain string `toml:"K8sDomain" json:"K8sDomain"`
	// ETCD 不需要从ETCD获取, 因此不加json标签.
	ETCD etcdx.Info `toml:"ETCD"`
	// Configuring 不需要从ETCD获取, 因此不加json标签.
	Configuring etcdx.ConfigInfo `toml:"Configuring"`
	// Register 不需要从ETCD获取, 因此不加json标签.
	Register etcdx.ServiceInfo `toml:"Register"`
	// GrpcRegister 不需要从ETCD获取, 因此不加json标签.
	GrpcRegister etcdx.ServiceInfo `toml:"GrpcRegister"`
	// Discover 不需要从ETCD获取, 因此不加json标签.
	Discover []etcdx.ServiceInfo `toml:"Discover"`

	Loggers    logger.LoggersInfo `toml:"Loggers" json:"Loggers"`
	Kafka      kafkax.Info        `toml:"Kafka" json:"Kafka"`
	MySQL      mysqlx.Info        `toml:"MySQL" json:"MySQL"`
	Mongo      mongox.Info        `toml:"Mongo" json:"Mongo"`
	Redis      redisx.Info        `toml:"Redis" json:"Redis"`
	Clickhouse clickhousex.Info   `toml:"Clickhouse" json:"Clickhouse"`

	Language string `toml:"Language" json:"Language"`

	// 从配置中心获取的配置清单存放在这里
	Handler MapHandel
}

type option struct {
	LoadedFile bool
}

type OptionHandler func(opt *option)

func WithFile(opt *option) {
	opt.LoadedFile = true
}

func InitConf(addr, serviceName, env, buildDate, buildVersion string, handlers ...OptionHandler) {
	if env == "" {
		env = environment.Active().Value()
	}
	opt := new(option)
	for _, handle := range handlers {
		handle(opt)
	}
	if opt.LoadedFile {
		// 加载配置文件
		initConfigFile(env)
	}

	// 初始化服务监听地址
	initListen(addr)

	// 初始化服务注册与服务发现相关配置
	initRegisterDiscover(serviceName, env)

	// 初始化卡夫卡配置(一定要在初始化服务注册配置的后面)
	initKafka()

	// 设置日志领域
	InitDomain(env)

	// 这个参数很多地方用得到, 保险补全
	if Default().ServiceName == "" {
		writeDefault().ServiceName = serviceName
	}
	InitBuildInfo(buildDate, buildVersion)

}

func initConfigFile(env string) {
	// 记录日志文件路径
	configName := env + DefaultConfigSuffix
	cwd, _ := os.Getwd()
	dirname := path.Join(cwd, DefaultConfigsPath)
	writeDefault().ConfFile = path.Join(dirname, configName+"."+DefaultConfigFileType)
	fmt.Println("加载配置文件", Default().ConfFile)

	viper.SetConfigName(configName)
	viper.SetConfigType(DefaultConfigFileType)
	viper.AddConfigPath(DefaultConfigsPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 读取配置文件
	if err := viper.Unmarshal(writeDefault()); err != nil {
		panic(err)
	}

	// 监听变化
	viper.WatchConfig()
	// 如果发生了变化, 回调
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(writeDefault()); err != nil {
			panic(err)
		}
	})
	// 设置不限字段的配置对象, 用来读取业务代码中的自定义配置.
	writeDefault().Handler = NewMapHandel(make(map[string]interface{}))

	// 记录环境变量
	writeDefault().Env = env
}

func InitDomain(env string) {
	if Default().BusinessPlatform != "" {
		env = fmt.Sprintf("%s:%s", Default().BusinessPlatform, env)
	}
	domain := fmt.Sprintf("[%s]%s-%s-%s", env, Default().Register.ProjectName, Default().Register.ProjectVersion, Default().Register.Name)
	writeDefault().Loggers.SetDomain(domain)
}

func InitBuildInfo(buildDate, buildVersion string) {
	if buildDate != "" {
		writeDefault().BuildInfo.BuildDate = buildDate
	}
	if buildVersion != "" {
		writeDefault().BuildInfo.BuildVersion = buildVersion
	}
}

func initListen(addr string) {
	// 初始化监听地址
	if Default().Register.Port > 0 {
		writeDefault().Register.Listen = fmt.Sprintf(":%v", Default().Register.Port)
	}
	if addr != "" {
		writeDefault().Register.Listen = addr
	}
	// 如果命令行与配置文件中都没有, 设置默认值
	if Default().Register.Listen == "" {
		writeDefault().Register.Listen = DefaultProjectListen
		fmt.Println("警告: 配置文件中没有配置 Register.Listen, 命令行也没有输入-addr, 设置监听地址为默认值:", Default().Register.Listen)
	}
}

func initRegisterDiscover(serviceName, env string) {
	if len(Default().K8sDomain) > 0 {
		writeDefault().Register.K8sDomain = Default().K8sDomain
	}
	// 最终用监听端口替换配置的端口
	ports := strings.Split(Default().Register.Listen, ":")
	writeDefault().Register.Port, _ = strconv.Atoi(ports[len(ports)-1])
	// 服务注册心跳时间
	if Default().Register.TTL <= 0 {
		writeDefault().Register.TTL = DefaultRegisterTTL
	}

	// 设置注册中心命名空间(一般用来区分环境)
	if env != "" {
		writeDefault().Register.Namespace = env
	}

	if Default().Register.Namespace == "" {
		writeDefault().Register.Namespace = environment.Active().Value()
	}
	// 设置service name
	if serviceName != "" {
		writeDefault().Register.Name = serviceName
	}
	// 服务协议
	if Default().Register.Scheme == "" {
		writeDefault().Register.Scheme = DefaultRegisterScheme
	}
	// 健康检查返回值
	if Default().Register.HealthVerdict == nil {
		writeDefault().Register.HealthVerdict = map[string]string{DefaultRegisterHealthKey: DefaultRegisterHealthVal}
	}
	// 健康检查路径
	if Default().Register.HealthPath == "" {
		writeDefault().Register.HealthPath = DefaultRegisterHealthPath
	}
	// 服务IP地址
	if Default().Register.Addr == "" {
		if len(Default().K8sDomain) > 0 {
			writeDefault().Register.Addr = fmt.Sprintf("%s.%s", Default().Register.Name, strings.TrimPrefix(Default().K8sDomain, "."))
		} else {
			writeDefault().Register.Addr = ipTools.LocalIP()
		}
	}
	// 设置服务ID,
	_ = writeDefault().Register.SetID()
	if Default().Discover == nil {
		writeDefault().Discover = make([]etcdx.ServiceInfo, 0)
	}

	// 配置中心
	if reflect.DeepEqual(Default().Configuring, etcdx.ConfigInfo{}) {
		writeDefault().Configuring = etcdx.ConfigInfo(Default().Register)
		writeDefault().Configuring.Prefix = ""
	}
	// 设置微服务前缀(用来在etcd中归纳服务与区分)
	writeDefault().Register.SetPrefix()
	// 设置配置中心的前缀
	writeDefault().Configuring.SetPrefix()
	// 设置服务名称到外层, 方便调用
	writeDefault().ServiceName = Default().Register.Name
	writeDefault().Discover = append(Default().Discover, Default().Register)
}

func initKafka() {
	if Default().Kafka.Group.GroupID == "" || !strings.HasPrefix(Default().Kafka.Group.GroupID, "zh-") {
		writeDefault().Kafka.Group.GroupID = Default().Register.Name
	}
}
