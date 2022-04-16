package service

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/asynccnu/card_service_v2/pkg/errno"
)

// 和网页有关函数（基础）

type accountRequestParams struct {
	lt         string
	execution  string
	eventId    string
	submit     string
	jsessionid string
}

var TIMEOUT = time.Duration(30 * time.Second)

// 预处理，打开 account.ccnu.edu.cn 获取模拟登陆需要的表单字段
func makeAccountPreflightRequest() (*accountRequestParams, error) {
	var (
		JSESSIONID string
		lt         string
		execution  string
		_eventId   string
	)

	params := &accountRequestParams{}

	// 初始化 http client
	client := http.Client{
		Timeout: TIMEOUT,
	}

	// 初始化 http request
	request, err := http.NewRequest("GET", "https://account.ccnu.edu.cn/cas/login", nil)
	if err != nil {
		log.Println(err)
		return params, err
	}

	// 发起请求
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return params, err
	}

	// 读取 Body
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return params, err
	}

	// 获取 Cookie 中的 JSESSIONID
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			JSESSIONID = cookie.Value
		}
	}

	if JSESSIONID == "" {
		log.Println("Can not get JSESSIONID")
		return params, errors.New("Can not get JSESSIONID")
	}

	// 正则匹配 HTML 返回的表单字段
	ltReg := regexp.MustCompile("name=\"lt\".+value=\"(.+)\"")
	executionReg := regexp.MustCompile("name=\"execution\".+value=\"(.+)\"")
	_eventIdReg := regexp.MustCompile("name=\"_eventId\".+value=\"(.+)\"")

	bodyStr := string(body)

	ltArr := ltReg.FindStringSubmatch(bodyStr)
	if len(ltArr) != 2 {
		log.Println("Can not get form paramater: lt")
		return params, errors.New("Can not get form paramater: lt")
	}
	lt = ltArr[1]

	execArr := executionReg.FindStringSubmatch(bodyStr)
	if len(execArr) != 2 {
		log.Println("Can not get form paramater: execution")
		return params, errors.New("Can not get form paramater: execution")
	}
	execution = execArr[1]

	_eventIdArr := _eventIdReg.FindStringSubmatch(bodyStr)
	if len(_eventIdArr) != 2 {
		log.Println("Can not get form paramater: _eventId")
		return params, errors.New("Can not get form paramater: _eventId")
	}
	_eventId = _eventIdArr[1]

	params.lt = lt
	params.execution = execution
	params.eventId = _eventId
	params.submit = "LOGIN"
	params.jsessionid = JSESSIONID

	return params, nil
}

// 模拟登陆并且获取 cookie
func makeAccountRequest(sid, password string, params *accountRequestParams, client *http.Client) (string, error) {
	v := url.Values{}
	v.Set("username", sid)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("_eventId", params.eventId)
	v.Set("submit", params.submit)

	request, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.jsessionid, strings.NewReader(v.Encode()))
	if err != nil {
		log.Println(err)
		return "", err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 验证是否登录成功
	reg := regexp.MustCompile("class=\"success\"")
	matched := reg.MatchString(string(body))
	if !matched {
		log.Println("Wrong sid or pwd")
		return "", errno.ErrPasswordIncorrect
	}

	// 获取 cookie
	var sessionID, routeportal string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			sessionID = cookie.Value
		} else if cookie.Name == "routeportal" {
			routeportal = cookie.Value
		}
	}

	// 获取 token
	token, err := GetPortalTokenFrom(sessionID, routeportal, client)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return token, err
}

// 请求获取 token
func GetPortalTokenFrom(sessionID, routeportal string, client *http.Client) (string, error) {
	request, err := http.NewRequest("GET", "http://one.ccnu.edu.cn/cas/login_portal;jsessionid="+sessionID, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	chanResp := make(chan *http.Response)
	f := func() {
		resp, err := client.Do(request)
		if err != nil {
			resp, err = client.Do(request)
		}
		chanResp <- resp
	}
	for i := 0; i < 3; i++ {
		go f()
	}

	go func() {
		time.Sleep(time.Second * 10)
		chanResp <- nil
	}()

	resp := <-chanResp
	if resp == nil {
		return "", errors.New("get http://one.ccnu.edu.cn wrong")
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PORTAL_TOKEN" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("get token failed")
}
