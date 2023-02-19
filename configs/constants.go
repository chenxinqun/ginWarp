package configs

var UI = `
█╗║╔╚═╝╗║╔█╔╚═╝╗█╗║╔╚═╝╗║╔█
----------ginWarp----------
█╗║╔╚═╝╗║╔█╔╚═╝╗█╗║╔╚═╝╗║╔█
`

func SetUI(ui string) {
	UI = ui
}

const (
	DefaultBakPath        = "./bak"
	DefaultConfigsPath    = "./configs"
	DefaultEtcdBakName    = "etcd_conf_bak"
	DefaultConfigFileType = "toml"
	DefaultConfigSuffix   = "Configs"
	EtcdTLSCaFile         = DefaultConfigsPath + "/etcdtls/ca.pem"
	EtcdTLSCertFile       = DefaultConfigsPath + "/etcdtls/etcd-client.pem"
	EtcdTLSCertKeyFile    = DefaultConfigsPath + "/etcdtls/etcd-client-key.pem"
	// DefaultProjectListen 默认监听地址
	DefaultProjectListen      = ":8888"
	DefaultRegisterTTL        = 10
	DefaultRegisterScheme     = "http"
	DefaultRegisterHealthKey  = "Status"
	DefaultRegisterHealthVal  = "ok"
	DefaultRegisterHealthPath = "/system/health"
)
