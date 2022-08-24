package config

import (
	"bytes"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/constant"
	agolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"main/app/common/log"
	"main/app/utils"
)

const (
	appID   = "10001" // apollo appid
	cluster = "dev"   // apollo 集群类型
	ip      = ""      // apollo agolloClient service ip
	secret  = ""      // apollo app secret

	// app下的所有命名空间
	namespaces = "gohu.yaml,database-dsn,oauth.yaml,user.yaml,NSQ"

	// app default namespace
	// 有关服务器的配置
	mainNamespace = "gohu.yaml"

	// app database namespace
	// 有关数据库的配置
	databaseNamespace = "database-dsn"

	// app mq namespace
	// 有关消息队列的配置
	NSQNamespace = "NSQ"
)

type Agollo struct {
	client agollo.Client
	vipers map[string]*viper.Viper
}

var agolloClient *Agollo

// NewConfigClient 获取Agollo客户端
func NewConfigClient() (c *Agollo, err error) {
	c = new(Agollo)
	vipers := make(map[string]*viper.Viper)
	c.vipers = vipers
	appConfig := &agolloConfig.AppConfig{
		AppID:            appID,
		Cluster:          cluster,
		IP:               ip,
		NamespaceName:    namespaces,
		IsBackupConfig:   true,       // 是否在本地备份
		BackupConfigPath: "./config", // 备份文件路径
		Secret:           secret,
		MustStart:        true,
	}
	// 客户端不解析 content, viper 来解析
	extension.AddFormatParser(constant.Properties, &emptyParser{})
	extension.AddFormatParser(constant.YAML, &emptyParser{})

	// 设置 apollo 的日志器
	logger := log.GetSugaredLogger()
	agollo.SetLogger(logger)

	client, err := agollo.StartWithConfig(func() (*agolloConfig.AppConfig, error) {
		return appConfig, nil
	})

	// 设置 配置监听功能(此项目暂时用不到,先不写)
	//client.AddChangeListener(&CustomChangeListener{})

	c.client = client
	agolloClient = c

	return
}

func GetConfigClient() (*Agollo, error) {
	if agolloClient == nil {
		return NewConfigClient()
	}
	return agolloClient, nil
}

func (c *Agollo) GetViper(namespace string) (*viper.Viper, error) {
	if v, ok := c.vipers[namespace]; ok {
		return v, nil
	} else {
		v := viper.New()
		namespaceType := utils.GetNamespaceType(namespace)
		v.SetConfigType(namespaceType)
		buffer := bytes.NewBufferString(c.client.GetConfig(namespace).GetValue("content"))
		if namespaceType == "properties" {
			buffer = bytes.NewBufferString(c.client.GetConfig(namespace).GetContent())
		}
		err := v.ReadConfig(buffer)
		if err != nil {
			return nil, fmt.Errorf("get viper failed, err: %v", err)
		}
		return v, nil
	}
}

func (c *Agollo) UnmarshalServiceConfig(namespace, serviceType, serviceName string, dst interface{}) (err error) {
	if serviceType == "api" {
		return c.UnmarshalKey(namespace, serviceType, dst)
	}
	return c.UnmarshalKey(namespace, serviceType+"."+serviceName, dst)
}

func (c *Agollo) UnmarshalKey(namespace, key string, dst interface{}) (err error) {
	v, err := c.GetViper(namespace)
	if err != nil {
		return
	}
	err = v.UnmarshalKey(key, dst, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.Squash = true
	})
	return
}

// GetMysqlDsn 返回 mysql DSN
func (c *Agollo) GetMysqlDsn(namespace string) string {
	config := c.client.GetConfig(namespace)
	serverNum := config.GetIntValue("database.mysql.serverNum", 0) // 使用 mysql 几号服务器(考虑以后将数据库分布式化,使用 TiDB )

	databaseConfig := c.client.GetConfig(databaseNamespace)

	// 拼接dsn字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Asia%%2FShanghai",
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.username", serverNum)), // 数据库用户名
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.password", serverNum)), // 数据库密码
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.address", serverNum)),  // 数据库地址
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.port", serverNum)),     // 数据库端口
		config.GetValue("database.mysql.databaseName"),                             // mysql 的数据库名字
		config.GetValue("database.mysql.databaseCharset"),                          // mysql 的数据库使用的字符集
	)
	return dsn
}

// NewRedisOptions 返回 *redisOptions
func (c *Agollo) NewRedisOptions(namespace string) *redis.Options {
	config := c.client.GetConfig(namespace)
	serverNum := config.GetIntValue("database.redis.serverNum", 0) // 使用 redis 几号服务器

	databaseConfig := c.client.GetConfig(databaseNamespace)

	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			databaseConfig.GetValue(fmt.Sprintf("redis.server%d.address", serverNum)),
			databaseConfig.GetValue(fmt.Sprintf("redis.server%d.port", serverNum)),
		),
		Password: databaseConfig.GetValue(fmt.Sprintf("redis.server%d.password", serverNum)),
		DB:       config.GetIntValue("database.redis.databaseNum", 0), // 使用 redis 几号数据库
	}
}

func (c *Agollo) GetClientDetails() (clientAuths map[string]interface{}) {
	v, _ := c.GetViper("oauth.yaml")
	return v.GetStringMap("Client")
}
