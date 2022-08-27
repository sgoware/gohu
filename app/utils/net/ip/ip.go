package ip

import (
	"github.com/imroc/req/v3"
	"github.com/spf13/cast"
	"github.com/thedevsaddam/gojsonq/v2"
	apollo "main/app/common/config"
	"strings"
)

var domain string

func GetIpLocFromApi(ip string) (loc string) {
	var err error
	if domain == "" {
		domain, err = apollo.GetDomain()
		if err != nil {
			return "未知"
		}
	}
	apiAddr := "http://ip." + domain + "/api/parse?ip=" + ip
	res, err := req.NewRequest().Get(apiAddr)
	j := gojsonq.New().FromString(res.String())
	if cast.ToBool(j.Find("ok")) == false {
		return "未知"
	}
	j.Reset()
	locStr := j.Find("location").(string)
	output := strings.Split(locStr, "|")
	if output[0] != "中国" {
		return output[0]
	}
	strings.Trim(output[2], "省")
	return output[2]
}
