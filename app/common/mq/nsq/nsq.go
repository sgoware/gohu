package nsq

import (
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/reiver/go-telnet"
	apollo "main/app/common/config"
	"math/rand"
	"strings"
	"time"
)

func GetConfig() (*nsq.Config, error) {
	config := nsq.NewConfig()

	configClient, err := apollo.GetConfigClient()
	if err != nil {
		return nil, err
	}

	v, err := configClient.GetViper(apollo.NSQNamespace)
	if err != nil {
		return nil, err
	}

	config.AuthSecret = v.GetString("nsq-auth.secret")

	return config, nil

}

func MustGetNSQDAddr() string {
	addr, err := GetNSQDAddr()
	if err != nil {
		return ""
	}
	return addr
}

func GetNSQDAddr() (string, error) {
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		return "", err
	}
	v, err := configClient.GetViper(apollo.NSQNamespace)
	if err != nil {
		return "", err
	}

	addrsSrc := v.GetString("nsq-nsqd.tcp")
	addrs := strings.Split(addrsSrc, ";")
	addrNum := len(addrs)

	// 简单的负载均衡
	rand.Seed(time.Now().Unix())
	serverNum := rand.Intn(addrNum - 1)

	cnt := 0
	for i := serverNum; cnt != addrNum-1; i++ {
		addr := addrs[i]
		// 检测端口连通性
		_, err := telnet.DialTo(addr)
		if err == nil {
			return addr, nil
		}

		if i == addrNum {
			i = 0
		}
		cnt++
	}
	return "", errors.New("no servers available")
}

func MustGetNSQLookupAddrs() []string {
	addrs, err := GetNSQLookupdAddrs()
	if err != nil {
		return nil
	}
	return addrs
}

func GetNSQLookupdAddrs() ([]string, error) {
	configClient, err := apollo.GetConfigClient()
	if err != nil {
		return nil, err
	}

	v, err := configClient.GetViper(apollo.NSQNamespace)
	if err != nil {
		return nil, err
	}

	addrsSrc := v.GetString("nsq-lookupd.http")
	addrs := strings.Split(addrsSrc, ";")

	return addrs, nil
}
