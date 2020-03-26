package httpkit

import (
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

type HttpClient struct {
	client *http.Client
	bb     *blackboardkit.BlackBoradKit
}

func NewHttpClient(timeout uint, maxIdle int) *HttpClient {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: maxIdle,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(timeout) * time.Second, //建立连接的等待时间
				KeepAlive: 1000 * time.Second,
			}).Dial,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	hc := &HttpClient{}
	hc.client = client
	hc.bb = blackboardkit.NewBlockBorad("httpkit", "httpclient", "http客户端记录")
	return hc
}

//简化的post
func (hc *HttpClient) Post(url string, body string) (string, error) {
	var buf io.Reader
	buf = strings.NewReader(body)

	t := hc.bb.Start("http post: " + url)
	res, err := hc.client.Post(url, "application/x-www-form-urlencoded;charset=utf-8", buf)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http post error: ", " url="+url, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http post read error: ", " url="+url, "err=", err)
		return "", err
	}
	hc.bb.Log("http post result: ", " url="+url, "body="+body, " resp="+string(content))
	return string(content), nil
}

func (hc *HttpClient) PostForm(url string, data url.Values) (string, error) {
	t := hc.bb.Start("http post: " + url)
	res, err := hc.client.PostForm(url, data)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http post error: ", " url="+url, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http post read error: ", " url="+url, "err=", err)
		return "", err
	}
	hc.bb.Log("http post result:", " url="+url, "data="+fmt.Sprint(data), " resp: "+string(body))
	return string(body), nil
}

func (hc *HttpClient) Get(url string) (string, error) {
	t := hc.bb.Start("http get: " + url)
	res, err := hc.client.Get(url)
	hc.bb.End(t)
	if err != nil {
		hc.bb.Err("http get error:", "  url="+url, "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http get read error:", " url="+url, "err=", err)
		return "", err
	}
	hc.bb.Log("http get result:", "url="+url, "resp: "+string(body))
	return string(body), nil
}

func (hc *HttpClient) Do(req *http.Request) (string, error) {
	res, err := hc.client.Do(req)
	if err != nil {
		hc.bb.Err("http   error:", "  url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(req.Body), "err=", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		hc.bb.Err("http   error: ", " url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(req.Body), "err=", err)

		return "", err
	}
	hc.bb.Log("http  result:", "url="+req.URL.String(), "header="+fmt.Sprint(req.Header), "body="+fmt.Sprint(req.Body), " resp: "+string(body))
	return string(body), nil
}
