package ip

import (
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	apollo "main/app/common/config"
	"strings"
)

var domain string

func GetIpLocFromApi(ip string) (loc string) {
	var err error
	if domain == "" {
		domain, err = apollo.GetMainDomain()
		if err != nil {
			return "未知"
		}
	}
	apiAddr := "http://ip." + domain + "/api/parse?ip=" + ip
	res, err := req.NewRequest().Get(apiAddr)
	j := gjson.Parse(res.String())
	if j.Get("ok").Bool() == false {
		return "未知"
	}
	locStr := j.Get("location").String()
	output := strings.Split(locStr, "|")
	if output[0] != "中国" {
		return output[0]
	}
	strings.Trim(output[2], "省")
	return output[2]
}
