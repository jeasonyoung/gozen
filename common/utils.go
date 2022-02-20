package common

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/gcharset"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

//EncryptPassword 密码加密
func EncryptPassword(password, salt string) string {
	return gmd5.MustEncryptString(password + salt)
}

//GetClientIp 获取客户端IP
func GetClientIp(r *ghttp.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.GetClientIp()
	}
	return ip
}

func GetGeoByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "[::1]" || ip == "127.0.0.1" {
		return "内网IP"
	}
	url := "https://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := g.Client().GetBytes(url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return ""
	}
	return json.GetString("addr")
}

// CurrencyLong 货币转为为分
func CurrencyLong(currency interface{}) int64 {
	strArr := gstr.Split(gconv.String(currency), ".")
	switch len(strArr) {
	case 1:
		return gconv.Int64(strArr[0]) * 100
	case 2:
		if len(strArr[1]) == 1 {
			strArr[1] += "0"
		} else if len(strArr[1]) > 2 {
			strArr[1] = gstr.SubStr(strArr[1], 0, 2)
		}
		return gconv.Int64(strArr[0])*100 + gconv.Int64(strArr[1])
	}
	return 0
}
