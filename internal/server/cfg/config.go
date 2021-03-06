package serverConf

import (
	"fmt"
	_fileUtils "github.com/aaronchen2k/tester/internal/pkg/libs/file"
	"github.com/aaronchen2k/tester/internal/pkg/utils"
	serverConst "github.com/aaronchen2k/tester/internal/server/utils/const"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/jinzhu/configor"
)

var Config = struct {
	LogLevel string `yaml:"logLevel" env:"LogLevel" default:"info"`
	Debug    bool   `yaml:"debug" env:"Debug" default:"false"`
	BinData  bool   `default:"true" env:"BinData"`
	Https    bool   `default:"false" env:"Https"`
	CertPath string `default:"" env:"CertPath"`
	CertKey  string `default:"" env:"CertKey"`
	Port     int    `default:"8085" env:"Port"`
	Host     string `default:"127.0.0.1" env:"HostId"`
	Adapter  struct {
		VmPlatform        serverConst.VmPlatform        `yaml:"vmPlatform" env:"VmPlatform" default:"pve"`
		ContainerPlatform serverConst.ContainerPlatform `yaml:"containerPlatform" env:"ContainerPlatform" default:"portainer"`
	} `yaml:"adapter,flow"`
	Admin struct {
		UserName        string `env:"AdminUserName" default:"admin"`
		Name            string `env:"AdminName" default:"admin"`
		Password        string `env:"AdminPassword" default:"P2ssw0rd"`
		RoleName        string `env:"AdminRoleName" default:"admin"`
		RoleDisplayName string `env:"RoleDisplayName" default:"超级管理员"`
	} `yaml:"admin,flow"`
	DB    DBConfig `yaml:"db,flow"`
	Redis struct {
		Enable bool   `env:"RedisDisable" default:"false"`
		Host   string `env:"RedisHost" default:"localhost"`
		Port   string `env:"RedisPort" default:"6379"`
		Pwd    string `env:"RedisPwd" default:""`
	} `yaml:"redis,flow"`

	Limit struct {
		Disable bool    `env:"LimitDisable" default:"true"`
		Limit   float64 `env:"LimitLimit" default:"1"`
		Burst   int     `env:"LimitBurst" default:"5"`
	}
	Qiniu struct {
		Enable    bool   `env:"QiniuEnable" default:"false"`
		Host      string `env:"QiniuHost" default:""`
		Accesskey string `env:"QiniuAccesskey" default:""`
		Secretkey string `env:"QiniuSecretkey" default:""`
		Bucket    string `env:"QiniuBucket" default:""`
	}
	Options struct {
		UploadMaxSize int64 `env:"uploadMaxSize" default:"100"`
	}
}{}

type DBConfig struct {
	Prefix   string `yaml:"prefix" env:"DBPrefix" default:"tester_"`
	Name     string `yaml:"name" env:"DBName" default:"tester"`
	Adapter  string `yaml:"adapter" env:"DBAdapter" default:"sqlite3"`
	Host     string `yaml:"host" env:"DBHost" default:"localhost"`
	Port     string `yaml:"port" env:"DBPort" default:"3306"`
	User     string `yaml:"user" env:"DBUser" default:"root"`
	Password string `yaml:"password" env:"DBPassword" default:"P2ssw0rd"`
}

func Init() {
	// in current dir
	exeDir := _utils.GetExeDir()
	configPath := filepath.Join(exeDir, "server.yml")
	if !_fileUtils.FileExist(configPath) { // debug mode
		configPath = filepath.Join(exeDir, "cmd", "server", "server.yml")
	}

	fmt.Println(fmt.Sprintf("配置YML文件路径：%v", configPath))
	if err := configor.Load(&Config, configPath); err != nil {
		logger.Println(fmt.Sprintf("Config Path:%s ,Error:%s", configPath, err.Error()))
		return
	}

	if Config.Debug {
		fmt.Println(fmt.Sprintf("配置项：%+v", Config))
	}
}

func GetRedisUris() []string {
	addrs := make([]string, 0, 0)
	hosts := strings.Split(Config.Redis.Host, ";")
	ports := strings.Split(Config.Redis.Port, ";")
	for _, h := range hosts {
		for _, p := range ports {
			addrs = append(addrs, fmt.Sprintf("%s:%s", h, p))
		}
	}
	return addrs
}
