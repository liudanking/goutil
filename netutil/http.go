package netutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	UA_SAFARI = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14"
	UA_CHROME = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36"
	UA_IOS    = "Mozilla/5.0 (iPhone; CPU iPhone OS 9_2_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Mobile/13D15 Weixinheadline/Ios/1.1.0"
)

type HttpClient struct {
	client  *http.Client
	request *http.Request
	err     error
	cloned  bool
}

var defaultHttpClient *HttpClient

func init() {
	tr := &http.Transport{
		Dial: func(network, addr string) (conn net.Conn, err error) {
			return net.DialTimeout(network, addr, 5*time.Second)
		},
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   45 * time.Second,
	}

	defaultHttpClient = NewHttpClient(client)
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{
		client: client,
	}
}

func DefaultHttpClient() *HttpClient {
	return defaultHttpClient
}

func (hc *HttpClient) clone() *HttpClient {
	if hc.cloned {
		return hc
	}
	return &HttpClient{
		client:  hc.client,
		request: hc.request,
		err:     hc.err,
		cloned:  true,
	}
}

func (hc *HttpClient) valid() error {
	if hc.err != nil {
		return hc.err
	}
	if hc.request == nil {
		return errors.New("request is empty")
	}
	return nil
}

func (hc *HttpClient) RequestForm(method, addr string, params map[string]interface{}) *HttpClient {
	var request *http.Request
	var err error
	var addrUrl *url.URL
	hc = hc.clone()
	values := hc.map2Values(params)
	switch method {
	case "GET", "DELETE":
		addrUrl, err = url.Parse(addr)
		if err == nil {
			newValues := hc.mergeValues(values, addrUrl.Query())
			addrUrl.RawQuery = newValues.Encode()
			request, err = http.NewRequest(method, addrUrl.String(), nil)
		}
	case "PUT", "POST":
		request, err = http.NewRequest(method, addr, strings.NewReader(values.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	default:
		err = errors.New("method not supported")
	}

	hc.request = request
	hc.err = err
	return hc
}

func (hc *HttpClient) JSONBody(method, addr string, body io.Reader) *HttpClient {
	request, err := http.NewRequest(method, addr, body)
	if err == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	hc.request = request
	hc.err = err
	return hc
}

func (hc *HttpClient) mergeValues(v1, v2 url.Values) url.Values {
	v3 := url.Values{}
	for k, v := range v1 {
		for _, s := range v {
			v3.Add(k, s)
		}
	}
	for k, v := range v2 {
		for _, s := range v {
			v3.Add(k, s)
		}
	}
	return v3
}

func (hc *HttpClient) map2Values(m map[string]interface{}) url.Values {
	values := url.Values{}
	for k, v := range m {
		vType := reflect.TypeOf(v)
		switch vType.Kind() {
		case reflect.Array, reflect.Slice:
			vValue := reflect.ValueOf(v)
			for i := 0; i < vValue.Len(); i++ {
				values.Add(k, fmt.Sprintf("%v", vValue.Index(i)))
			}
		default:
			values.Add(k, fmt.Sprintf("%v", v))
		}
	}

	return values
}

func (hc *HttpClient) DoByte() ([]byte, int, error) {

	if err := hc.valid(); err != nil {
		return nil, 0, err
	}

	rsp, err := hc.client.Do(hc.request)
	if err != nil {
		return nil, 0, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	return data, rsp.StatusCode, err
}

func (hc *HttpClient) DoJSONRaw(rsp interface{}) ([]byte, error) {
	data, _, err := hc.DoByte()
	if err != nil {
		return nil, err
	}

	return data, json.Unmarshal(data, rsp)
}

func (hc *HttpClient) DoJSON(rsp interface{}) ([]byte, error) {
	data, _, err := hc.DoByte()
	if err != nil {
		return nil, err
	}
	return data, json.Unmarshal(data, rsp)
}

func (hc *HttpClient) Header(header map[string]string) *HttpClient {
	hcc := hc.clone()
	if err := hcc.valid(); err != nil {
		return hcc
	}

	for k, v := range header {
		hcc.request.Header.Set(k, v)
	}

	// fmt.Printf("header:%+v\n", hcc.request.Header)

	return hcc
}

func (hc *HttpClient) UserAgent(ua string) *HttpClient {
	hcc := hc.clone()
	if err := hcc.valid(); err != nil {
		return hcc
	}

	hcc.request.Header.Set("User-Agent", ua)
	return hcc
}

// Proxy always deep clone a new http client
func (hc *HttpClient) Proxy(proxy func(*http.Request) (*url.URL, error)) *HttpClient {
	hcc := hc.clone()
	if tr, ok := hc.client.Transport.(*http.Transport); ok {
		// copy transport
		trNew := *tr
		// set proxy
		trNew.Proxy = proxy

		// copy client
		hcc.client = &http.Client{
			Transport: &trNew,
			Timeout:   hcc.client.Timeout,
		}

	} else {
		hcc.err = errors.New("assert transport type failed")
	}
	return hcc
}
