package service

import(
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	"strings"
)
//两个和api有关函数

//获取校园卡流水
func DoList(User_id, Password, Limit, Page, Start, End string) string{
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return "false"
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return "false"
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	token := makeAccountRequest2(User_id, Password, params, &client)
	v := url.Values{}
	
	v.Set("limit", Limit)
	v.Set("page", Page)
	v.Set("start", Start)
	v.Set("end", End)
	v.Set("tranType", "")
	
	req, err := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/query_trans", strings.NewReader(v.Encode()))
	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return ""
	}

	result := string(bstr)
	return result
}

//获取校园卡余额以及状态，及在用或者丢失
func DoStatus(User_id, Password string) string{
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return ""
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return ""
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	token := makeAccountRequest2(User_id, Password, params, &client)
	v := url.Values{}
	
	req, _ := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/get_info", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(string(v.Encode()))
		log.Println("====")
		log.Println(err)
		return ""
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return ""
	}

	result := string(bstr)
	return result
}