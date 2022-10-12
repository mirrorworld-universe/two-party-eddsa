package settings

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Base *BaseCfg
	Log  *LogCfg
	//Rpc        *RpcCfg
	DB *DBCfg
	//Mongo      *MongoCfg
	//Img        *Img
	//Kafka      *KafkaConf
	//Redis      *RedisConf
	//Node       *NodeCfg // solana node
	//Syncer     *SyncerCfg
	//Middleware *Middleware
	//Sso        *SSO
}

type BaseCfg struct {
	Port     string
	Model    string
	Env      string
	SSO      string
	User     string
	Password string
	P0Url    string
	P1Url    string
}

type LogCfg struct {
	Path    string
	MaxSize int
	MaxAge  int
	Backups int
}

type DBCfg struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  int
}

var cfg = Config{}

func InitConfig(path string) *Config {
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		panic("load config error: " + err.Error() + path)
	}
	return &cfg
}
