package service

import(
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
	
)
//和网页有关函数（基础）

//结构体
type accountReqeustParams struct {
	lt         string
	execution  string
	eventId   string
	submit     string
	JSESSIONID string
}


//确认模拟登陆是否成功
func ConfirmUser(sid string, pwd string) bool {
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return false
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return false
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	//fmt.Println(params)
	result := makeAccountRequest(sid, pwd, params, &client)

	return result
}

// 预处理，打开 account.ccnu.edu.cn 获取模拟登陆需要的表单字段
func makeAccountPreflightRequest() (*accountReqeustParams, error) {
	var JSESSIONID string
	var lt string
	var execution string
	var eventId string

	params := &accountReqeustParams{}

	// 初始化 http client
	client := http.Client{
		//Timeout: TIMEOUT,
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
		//fmt.Println(cookie.Value)
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
	eventIdReg := regexp.MustCompile("name=\"eventId\".+value=\"(.+)\"")

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

	eventIdArr := eventIdReg.FindStringSubmatch(bodyStr)
	if len(eventIdArr) != 2 {
		log.Println("Can not get form paramater: eventId")
		return params, errors.New("Can not get form paramater: eventId")
	}
	eventId = eventIdArr[1]

	params.lt = lt
	params.execution = execution
	params.eventId = eventId
	params.submit = "LOGIN"
	params.JSESSIONID = JSESSIONID

	return params, nil
}

// 进行模拟登陆
func makeAccountRequest(sid, password string, params *accountReqeustParams, client *http.Client) bool {
	v := url.Values{}
	v.Set("username", sid)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("eventId", params.eventId)
	v.Set("submit", params.submit)

	request, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.JSESSIONID, strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	reg := regexp.MustCompile("class=\"success\"")
	matched := reg.MatchString(string(body))
	if !matched {
		log.Println("Wrong sid or pwd")
		return false
	}
	
	return true
}

//模拟登陆并且获取cookie
func makeAccountRequest2(sid, password string, params *accountReqeustParams, client *http.Client)( w string) {
	v := url.Values{}
	v.Set("username", sid)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("eventId", params.eventId)
	v.Set("submit", params.submit)

	request, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.JSESSIONID, strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
	}

	var s1, s2 string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			s1 = cookie.Value
		} else if cookie.Name == "routeportal" {
			s2 = cookie.Value
		}
	}
	if err,w = GetToken(s1, s2, client);err!=nil{
		log.Println(err)
		return ""
	}

	return w
}

//获取token
func GetToken(sessionid, routeportal string, client *http.Client) (error,string ){
	request, err := http.NewRequest("GET", "http://one.ccnu.edu.cn/cas/login_portal;jsessionid=" + sessionid, nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		return err,""
	}
	
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PORTAL_TOKEN" {
			return nil,cookie.Value
		}
	}
	return err,""
}

