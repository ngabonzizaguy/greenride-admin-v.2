package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	JSON_HEADER = "application/json; charset=utf-8"
	FORM_HEADER = "application/x-www-form-urlencoded; charset=utf-8"
)

var (
	ConnectionTimeout = 60 * time.Second
	HandshakeTimeout  = 60 * time.Second
	ReadWriteTimeout  = 60 * time.Second
)

// GlobalHttpClientPool 使用 sync.Pool 存储 http.Client 实例
var GlobalHttpClientPool = sync.Pool{
	New: func() any {
		return NewHttpClient()
	},
}

// GetHttpClient 从池中获取一个 http.Client
func GetHttpClient() *http.Client {
	return GlobalHttpClientPool.Get().(*http.Client)
}

// PutHttpClient 将 http.Client 放回池中
func PutHttpClient(cli *http.Client) {
	GlobalHttpClientPool.Put(cli)
}

func NewHttpClient() *http.Client {
	// pre-create client connection pool
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: ConnectionTimeout,
		}).DialContext,
		TLSHandshakeTimeout: HandshakeTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // comment this line to verify remote server certificate
		},
	}
	cli := &http.Client{
		Transport: transport,
		Timeout:   ConnectionTimeout,
	}
	return cli
}

func NewHttpClientWithProxy(proxy string) *http.Client {
	// pre-create client connection pool
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: ConnectionTimeout,
		}).DialContext,
		TLSHandshakeTimeout: HandshakeTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // comment this line to verify remote server certificate
		},
	}

	if proxy != "" {
		url, err := url.Parse(proxy)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(url)
	}
	cli := &http.Client{
		Transport: transport,
		Timeout:   ConnectionTimeout,
	}
	return cli
}

func PostWithHeader(url string, data []byte, headers map[string]string) (string, *http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	cli := GetHttpClient()
	resp, err := cli.Do(req)
	PutHttpClient(cli)
	if err != nil {
		return "", resp, err
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	return string(respData), resp, nil
}

func GetWithHeader(url string, headers map[string]string) (string, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	cli := GetHttpClient()
	resp, err := cli.Do(req)
	PutHttpClient(cli)
	if err != nil {
		return "", resp, err
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	return string(respData), resp, nil
}
func PostClientJsonDataWithHeader(cli *http.Client, url string, data []byte, headers map[string]string) (string, *http.Response, error) {
	if len(headers) == 0 {
		headers = map[string]string{}
	}
	headers["content-type"] = JSON_HEADER
	return PostClientDataWithHeader(cli, url, data, headers)
}

func PostClientDataWithHeader(cli *http.Client, url string, data []byte, headers map[string]string) (string, *http.Response, error) {
	reqData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqData))
	if err != nil {
		return "", nil, err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	//req.Header.Set("content-type", "application/json; charset=utf-8")
	resp, err := cli.Do(req)
	if err != nil {
		return "", resp, err
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	return string(respData), resp, nil
}

func PostClientWithHeader(cli *http.Client, url string, data any, headers map[string]string) (string, error, *http.Response) {
	reqData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqData))
	if err != nil {
		return "", err, nil
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	//req.Header.Set("content-type", "application/json; charset=utf-8")
	resp, err := cli.Do(req)
	if err != nil {
		return "", err, resp
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	return string(respData), nil, resp
}

func GetClientWithHeader(cli *http.Client, url string, headers map[string]string) (string, error, *http.Response) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err, nil
	}
	if len(headers) > 0 {
		for key, header := range headers {
			req.Header.Set(key, header)
		}
	}
	resp, err := cli.Do(req)
	if err != nil {
		return "", err, resp
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	return string(respData), nil, resp
}
func PostJson(url string, data map[string]any) (string, *http.Response, error) {
	headers := map[string]string{
		"content-type": JSON_HEADER,
	}
	return PostWithHeader(url, ToJsonByte(data), headers)
}

func PostForm(url string, data map[string]any) (string, *http.Response, error) {
	headers := map[string]string{
		"content-type": FORM_HEADER,
	}
	bytes := []byte(ToQueryUrl(data))
	return PostWithHeader(url, bytes, headers)
}
func PostJsonDataWithHeader(url string, data []byte, headers map[string]string) (string, *http.Response, error) {
	if len(headers) == 0 {
		headers = map[string]string{}
	}
	headers["content-type"] = JSON_HEADER
	return PostWithHeader(url, data, headers)
}
func PostJsonWithHeader(url string, data map[string]any, headers map[string]string) (string, *http.Response, error) {
	if len(headers) == 0 {
		headers = map[string]string{}
	}
	headers["content-type"] = JSON_HEADER
	body := ToJsonByte(data)
	return PostWithHeader(url, body, headers)
}

func Get(url string) (string, *http.Response, error) {
	return GetWithHeader(url, nil)
}
