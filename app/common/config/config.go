package config

import (
	"bytes"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/constant"
	agolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"main/app/common/log"
	"main/app/utils"
	"os"
)

const (
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

type agolloConnConfig struct {
	AppID       string // apollo appid
	ClusterName string // apollo 集群类型
	IP          string // apollo agolloClient service ip
	Secret      string // apollo app secret
}

var agolloClient *Agollo

// NewConfigClient 获取Agollo客户端
func NewConfigClient() (c *Agollo, err error) {
	logger := log.GetSugaredLogger()
	connConfig := agolloConnConfig{
		AppID:       os.Getenv("APOLLO_APP_ID"),
		ClusterName: os.Getenv("APOLLO_CLUSTER_NAME"),
		IP:          os.Getenv("APOLLO_IP"),
		Secret:      os.Getenv("APOLLO_SECRET"),
	}
	logger.Debugf("connConfig: \n%v", connConfig)
	c = new(Agollo)
	vipers := make(map[string]*viper.Viper)
	c.vipers = vipers
	appConfig := &agolloConfig.AppConfig{
		AppID:            connConfig.AppID,
		Cluster:          connConfig.ClusterName,
		IP:               connConfig.IP,
		NamespaceName:    namespaces,
		IsBackupConfig:   true,       // 是否在本地备份
		BackupConfigPath: "./config", // 备份文件路径
		Secret:           connConfig.Secret,
		MustStart:        true,
	}
	// 客户端不解析 content, viper 来解析
	extension.AddFormatParser(constant.Properties, &emptyParser{})
	extension.AddFormatParser(constant.YAML, &emptyParser{})

	// 设置 apollo 的日志器
	agollo.SetLogger(logger)

	client, err := agollo.StartWithConfig(func() (*agolloConfig.AppConfig, error) {
		return appConfig, nil
	})

	// 设置 配置监听功能(此项目暂时用不到,先不写)
	//client.AddChangeListener(&CustomChangeListener{})

	c.client = client
	agolloClient = c

	logger.Info("Initialize config successfully!")

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

func (c *Agollo) GetClientDetails() (clientAuths map[string]interface{}) {
	v, _ := c.GetViper("oauth.yaml")
	return v.GetStringMap("Client")
}
