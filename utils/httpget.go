//基本的GET请求
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Responseinfo struct {
	Url				string
	Content     	string
	Headers			string
	Server			string
	Title			string
	Cert			string
}

// 获取证书内容，参考byro07/fwhatweb
func getCerts(resp *http.Response) []byte {
	var certs []byte
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		var str string
		if js, err := json.Marshal(cert); err == nil {
			certs = js
		}
		str = string(certs) + cert.Issuer.String() + cert.Subject.String()
		certs = []byte(str)
	}
	return certs
}

func GetRespon(url string, timeout int) (*http.Response, error){
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ParseRespon(resp *http.Response) Responseinfo{
	defer resp.Body.Close()
	var responseinfo Responseinfo
	body, _ := ioutil.ReadAll(resp.Body)
	responseinfo.Content = string(body)
	for k, v := range resp.Header {
		responseinfo.Headers += fmt.Sprintf("%v: %v\n", k, v[0])
	}
	responseinfo.Server = fmt.Sprintf("Server : %v\n", resp.Header["Server"])
	responseinfo.Title 	= RegexpTitle(string(body))
	responseinfo.Cert 	= string(getCerts(resp))
	return responseinfo
}

//todo: http和https同时存在时，有时会是两个不同系统
func Requrl(url string, timeout int) (Responseinfo, error) {
	var responseinfo Responseinfo
	httpurl := "http://" + url
	resp, err := GetRespon(httpurl, timeout)
	if err != nil {
		httpsurl := "https://" + url
		resp, err = GetRespon(httpsurl, timeout)
		if err != nil {
			return responseinfo, err
		}
	}
	responseinfo = ParseRespon(resp)

	return responseinfo, nil
}