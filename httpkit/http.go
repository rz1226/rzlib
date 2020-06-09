package httpkit

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rz1226/rzlib/blackboardkit"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HTTPClient struct {
	client *http.Client
	bb     *blackboardkit.BlackBoradKit
}

func NewHTTPClient(timeout uint, maxIdle int) *HTTPClient {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: maxIdle,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(timeout) * time.Second, //  建立连接的等待时间
				KeepAlive: 1000 * time.Second,
			}).Dial,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	hc := &HTTPClient{}
	hc.client = client
	hc.bb = blackboardkit.NewBlockBorad("httpkit", "httpclient", "http客户端记录")
	return hc
}

//  简化的post
func (hc *HTTPClient) Post(urlStr, body string) (string, error) {

	buf := strings.NewReader(body)

	t := hc.bb.Start("http post: " + urlStr)
	res, err := hc.client.Post(urlStr, "application/x-www-form-urlencoded;charset=utf-8", buf)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http post error: ", " url="+urlStr, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http post read error: ", " url="+urlStr, "err=", err)
		return "", err
	}
	hc.bb.Log("http post result: ", " url="+urlStr, "body="+body, "resp="+string(content))
	return string(content), nil
}


func (hc *HTTPClient) Post2(url string, bodyType string, body io.Reader) (string, error) {
	t := hc.bb.Start( "http post2: "+url)
	res, err := hc.client.Post(url, bodyType, body)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err( "http post error: url="+url, "err=",err)
		return "", err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err( "http post read error: url="+url, "err=",err)
		return "", err
	}
	hc.bb.Log( "http post2 result:"+url, " resp: "+string(content))
	return string(content), nil
}


func (hc *HTTPClient) PostForm(urlStr string, data url.Values) (string, error) {
	t := hc.bb.Start("http post: " + urlStr)
	res, err := hc.client.PostForm(urlStr, data)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http post error: ", " url="+urlStr, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http post read error: ", " url="+urlStr, "err=", err)
		return "", err
	}
	hc.bb.Log("http post result:", " url="+urlStr, "data="+fmt.Sprint(data), "resp: "+string(body))
	return string(body), nil
}

func (hc *HTTPClient) Get(urlStr string) (string, error) {
	t := hc.bb.Start("http get: " + urlStr)
	res, err := hc.client.Get(urlStr)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http get error:", "  url="+urlStr, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http get read error:", " url="+urlStr, "err=", err)
		return "", err
	}
	hc.bb.Log("http get result:", "url="+urlStr, "resp: "+string(body))
	return string(body), nil
}


//copy request.Body
func copyBody(r *http.Request) (string , error) {
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		// 把刚刚读出来的再写进去
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes), err
	}
	return "", errors.New("empty body")
}


func (hc *HTTPClient) Do(req *http.Request) (string, error) {
	// 复制body ,日志用
	bodyForLog, err := copyBody( req )
	if err != nil {
		bodyForLog = ""
	}

	res, err := hc.client.Do(req)

	if err != nil {
		hc.bb.Err("http   error:", "  url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(bodyForLog), "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http   error: ", " url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(bodyForLog), "err=", err)

		return "", err
	}
	hc.bb.Log("http  result:", "url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(bodyForLog), "resp: "+string(body))
	return string(body), nil
}
