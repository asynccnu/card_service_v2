package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type OriginCardInfo struct {
	CardInfo CardInfo `json:"cardInfo"`
}

// 用于解析网页json格式数据
type CardInfo struct {
	No           string  `json:"no"`
	DeptName     string  `json:"deptName"`
	StatusDesc   string  `json:"statusDesc"`
	Balance      float32 `json:"balance"`
	Xm           string  `json:"xm"`
	ValidityDate string  `json:"validityDate"`
	Status       string  `json:"status"`
	Username     string  `json:"username"`
}

// 一卡通信息请求
func MakeCardInfoRequest(sid, password string) (*CardInfo, error) {
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
		Jar:     jar,
	}

	token, err := makeAccountRequest(sid, password, params, &client)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	v := url.Values{}

	req, err := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/get_info", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(string(v.Encode()))
		log.Println(err)
		return nil, err
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result OriginCardInfo

	if err := json.Unmarshal([]byte(bstr), &result); err != nil {
		return nil, err
	}

	return &result.CardInfo, nil
}

type ConsumeBodyData struct {
	Result DealResult `json:"result"`
}

type DealResult struct {
	Rows []*DealRow `json:"rows"`
}

type DealRow struct {
	DealName   string  `json:"dealName"`   // 交易类型，消费/圈存机充值
	OrgName    string  `json:"orgName"`    // 消费地点，交易商方
	TransMoney float32 `json:"transMoney"` // 交易金额
	DealDate   string  `json:"dealDate"`   // 时间，格式：2020-01-18 07:25:32
	OutMoney   float32 `json:"outMoney"`   // 剩余余额
}

// 消费流水请求
func MakeConsumesRequest(sid, password, limit, page, start, end string) ([]*DealRow, error) {
	params, err := makeAccountPreflightRequest()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}

	token, err := makeAccountRequest(sid, password, params, &client)
	v := url.Values{}

	v.Set("limit", limit)
	v.Set("page", page)
	v.Set("start", start)
	v.Set("end", end)
	v.Set("tranType", "")

	req, err := http.NewRequest("POST", "http://one.ccnu.edu.cn/ecard_portal/query_trans", strings.NewReader(v.Encode()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	austr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", austr)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bstr, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	var result ConsumeBodyData

	if err = json.Unmarshal(bstr, &result); err != nil {
		return nil, err
	}

	return result.Result.Rows, nil
}
