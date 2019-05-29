package fone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/utils"
)

var (
	ErrDomainParams   = errors.New("params error: domain error")
	ErrNSParams       = errors.New("params error: name server")
	ErrMissNSParams   = errors.New("params error: miss authority")
	ErrDNSIDParams    = errors.New("params error: dnsid error")
	ErrMissInfoParams = errors.New("params error: miss host or user or passwd")
	ErrIDParams       = errors.New("params error: dnsid or record error")
)

const (
	FONETokenKey = "fone-token-const"
)

type Result struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Err    string      `json:"err"`
	Data   interface{} `json:"data"`
}

type FOneDNS struct {
	username string
	password string
	domain   string
	// token    string
	// tokenExpireTime time.Time // token 过期时间
	client *http.Client
}

// New new foneDNS
func New(host, username, password string, timeout int) (sdk.DNSSDK, error) {

	f := &FOneDNS{}

	if host == "" || username == "" || password == "" {
		return f, ErrMissInfoParams
	}

	if timeout <= 0 || timeout > 99 {
		timeout = 5
	}

	f.username = username
	f.password = password

	f.domain = strings.TrimSuffix(host, "/")
	f.domain = fmt.Sprintf("%s/api/v1", f.domain)

	f.client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	return f, nil
}

// Login login to edge system
func (f *FOneDNS) Login() (token string, err error) {

	p := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{
		UserName: f.username,
		Password: f.password,
	}

	var tmp struct {
		Token string `json:"token"`
	}

	err = f.httpDo("POST", "/user/login", p, &tmp)
	if err != nil {
		utils.Error("[fone] login error:", err)
		return
	}

	c.Set(FONETokenKey, tmp.Token, 50*time.Minute)

	utils.Infof("[fone] login token ==> %s", tmp.Token)
	fmt.Println("[login token] =>", tmp.Token)

	return tmp.Token, nil
}

// 如果接口返回成功，则返回 data 中的信息
func (f *FOneDNS) httpDo(method, uri string, request interface{}, response interface{}) (err error) {

	if response != nil && reflect.TypeOf(response).Kind() != reflect.Ptr {
		err = fmt.Errorf("response must be ptr")
		return
	}

	method = strings.ToUpper(method)

	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}

	uri = fmt.Sprintf("%s%s", f.domain, uri)

	token := ""
	if !strings.Contains(uri, "/user/login") {
		if tmpToken, ok := c.Get(FONETokenKey); !ok {
			// 说明是普通请求
			token, err = f.Login()
			if err != nil {
				utils.Error("[http] login error:", err)
				return
			}
		} else {
			token = tmpToken.(string)
		}
	}

	utils.Infof("[fone] http: %v, token: %v", uri, token)

	var req *http.Request
	if request == nil {
		req, err = http.NewRequest(method, uri, nil)
	} else {

		reqBody, err := json.Marshal(request)
		if err != nil {
			return err
		}

		req, err = http.NewRequest(method, uri, bytes.NewReader(reqBody))
	}

	if err != nil {
		utils.Error("[http] new request error:", err)
		return
	}

	req.Header.Add("Auth", token)

	resp, err := f.client.Do(req)
	if err != nil {
		utils.Error("[http] request do error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.Error("[http] status code:", resp.StatusCode)
		err = fmt.Errorf("status code error: %d", resp.StatusCode)
		return
	}

	var r Result

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		utils.Error("[http] json error:", err)
		return
	}

	if r.Status != 0 {
		utils.Errorf("[http] %v %v %v %v", method, r.Msg, r.Err, uri)
		return fmt.Errorf("response error: %v %v %v %v", method, uri, r.Msg, r.Err)
	}

	// fmt.Println((r.Data))

	if response == nil || r.Data == nil {
		return nil
	}

	tmpData, err := json.Marshal(r.Data)
	if err != nil {
		return
	}

	err = json.Unmarshal(tmpData, response)
	if err != nil {
		return err
	}

	return nil
}
