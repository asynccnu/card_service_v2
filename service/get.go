package service

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

//两个和api有关函数

//获取校园卡流水
func GetConsumeList(userId, password, limit, page, start, end string) (string, error) {
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return "", err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return "", err
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	token, err := makeAccountRequest2(userId, password, params, &client)
	v := url.Values{}

	v.Set("limit", limit)
	v.Set("page", page)
	v.Set("start", start)
	v.Set("end", end)
	v.Set("tranType", "")

	req, err := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/query_trans", strings.NewReader(v.Encode()))
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return "", nil
	}

	result := string(bstr)
	return result, nil
}

//获取校园卡余额以及状态，及在用或者丢失
func GetCardInfo(User_id, Password string) (string, error) {
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return "", err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return "", err
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	token, err := makeAccountRequest2(User_id, Password, params, &client)
	if err != nil {
		log.Println(err)
		return "", err
	}

	v := url.Values{}

	req, err := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/get_info", nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(string(v.Encode()))
		log.Println(err)
		return "", err
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	result := string(bstr)
	return result, nil
}
