package conf

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"

	"github.com/gopherty/wings/utils"
)

var instance *Configuration

// Register configuration register.
type Register struct {
}

// Configuration 配置模块
type Configuration struct {
	DB            *DB            `json:"DB" yaml:"DB"  jsonnet:"DB"`
	Server        *Server        `json:"Server" yaml:"Server" jsonnet:"Server"`
	GatewayServer *GatewayServer `json:"GatewayServer" yaml:"GatewayServer" jsonnet:"GatewayServer"`
	Logger        *Logger        `json:"Logger" yaml:"Logger" jsonnet:"Logger"`
	User          *User          `json:"User" yaml:"User" jsonnet:"User"`
	Secret        *Secret        `json:"Secret" yaml:"Secret" jsonnet:"Secret"`
}

// DB 数据库配置
type DB struct {
	// 数据库驱动
	Driver string `json:"Driver" yaml:"Driver" jsonnet:"Driver"`
	// 连接字符串
	Source string `json:"Source" yaml:"Source" jsonnet:"Source"`

	// 是否显示 SQL 语句
	ShowSQL bool `json:"ShowSQL" yaml:"ShowSQL" jsonnet:"ShowSQL"`
	// 数据库连接池数量
	MaxOpenConns int `json:"MaxOpenConns" yaml:"MaxOpenConns" jsonnet:"MaxOpenConns"`
	// 数据库连接最大空闲数
	MaxIdleConns int `json:"MaxIdleConns" yaml:"MaxIdleConns" jsonnet:"MaxIdleConns"`
	// 缓存大小
	Cached int `json:"Cached" yaml:"Cached" jsonnet:"Cached"`
}

// Server grpc server
type Server struct {
	Address string `json:"Address" yaml:"Address" jsonnet:"Address"`
	//证书验证文件
	CertFile string `json:"CertFile" yaml:"CertFile" jsonnet:"CertFile"`
	// 证书
	KeyFile string `json:"KeyFile" yaml:"KeyFile" jsonnet:"KeyFile"`
	// 是否为发布版本
	Release bool `json:"Release" yaml:"Release" jsonnet:"Release"`
}

// GatewayServer grpc gateway server
type GatewayServer struct {
	Address string `json:"Address" yaml:"Address" jsonnet:"Address"`
}

// Logger 日志配置
type Logger struct {
	// 日志等级: debug,info,warn,error,dpanic,panic,fatal
	Level string `json:"Level" yaml:"Level" jsonnet:"Level"`
	// 日志时间输出格式
	TimeFormat string `json:"TimeFormat" yaml:"TimeFormat" jsonnet:"TimeFormat"`
	// 是否开启开发模式
	Development bool `json:"Development" yaml:"Development" jsonnet:"Development"`
	// 日志输出格式
	Encoding string `json:"Encoding" yaml:"Encoding" jsonnet:"Encoding"`
	// 日志文件输出位置
	LogsPath string `json:"LogsPath" yaml:"LogsPath" jsonnet:"LogsPath"`
}

// User system username and password
type User struct {
	Name   string `json:"Name"  yaml:"Name" jsonnet:"Name"`
	Passwd string `json:"Passwd"  yaml:"Passwd" jsonnet:"Passwd"`
}

// Secret JWT 密钥
type Secret struct {
	Access  string `json:"Access"  yaml:"Access" jsonnet:"Access"`
	Refresh string `json:"Refresh"  yaml:"Refresh" jsonnet:"Refresh"`
}

// Name .
func (Register) Name() string {
	return "Common.Configuration"
}

// Regist 注册配置模块
func (Register) Regist() (err error) {
	instance = &Configuration{}
	path, err := utils.BasePath()
	if err != nil {
		return
	}
	path = filepath.Join(path, "config.yml")
	_, err = instance.loader(path)
	if err != nil {
		return
	}
	if instance == nil {
		err = errors.New("configuration load failed")
	}
	return
}

// 加载配置文件
func (c *Configuration) loader(path string) (b []byte, err error) {
	path = strings.ToLower(strings.TrimSpace(path))
	switch filepath.Ext(path) {
	case ".json":
		b, err = utils.ReadFile(path)
		if err != nil {
			return
		}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(b, c)
		if err != nil {
			return
		}
	case ".yml":
		b, err = utils.ReadFile(path)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(b, c)
		if err != nil {
			return
		}
	case ".jsonnet":
		b, err = utils.ReadFile(path)
		if err != nil {
			return
		}
		vm := jsonnet.MakeVM()
		var jsonStr string
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		jsonStr, err = vm.EvaluateSnippet(path, string(b))
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(jsonStr), c)
		if err != nil {
			return
		}
	default:
		err = errors.New("format not support")
	}
	return
}

// Instance 获取配置对象
func Instance() *Configuration {
	return instance
}
