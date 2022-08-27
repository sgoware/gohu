package ip

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/cast"
	"github.com/thedevsaddam/gojsonq/v2"
	apollo "main/app/common/config"
	"net"
	"strings"
)

var domain string

func GetIpLocFromApi(ip string) (loc string, err error) {
	if domain == "" {
		domain, err = apollo.GetDomain()
		if err != nil {
			return "", fmt.Errorf("get domain failed, %v", err)
		}
	}
	apiAddr := "http://ip." + domain + "/api/parse?ip=" + ip
	res, err := req.NewRequest().Get(apiAddr)
	j := gojsonq.New().FromString(res.String())
	if cast.ToBool(j.Find("ok")) == false {
		return "", errors.New("query api failed")
	}
	j.Reset()
	locStr := j.Find("location").(string)
	output := strings.Split(locStr, "|")
	if output[0] != "中国" {
		return output[0], nil
	}
	strings.Trim(output[2], "省")
	return output[2], nil
}

func GetIPLocFromLocal(rawIp string) (loc string, err error) {

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return "", errors.New("open db failed")
	}
	defer db.Close()

	ip := net.ParseIP(rawIp)
	record, err := db.City(ip)
	if err != nil {
		return "", errors.New("query failed")
	}
	fmt.Println(record.Country.Names["zh-CN"])
	fmt.Println(record.City.Names["zh-CN"])
	//loc = ToChineseLocation(record.Country.Names["zh-CN"],record.Location)
	return loc, nil
}
