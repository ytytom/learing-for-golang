package zabbix

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

var (
	HTTPRequestTimeOut     = time.Duration(5) * time.Second
	HTTPRequestRetryPeriod = time.Duration(5) * time.Second
)

//获取timestamp
func getTimeStr() string {
	timeDelta, _ := time.ParseDuration("-3h")
	stime := time.Now().Add(timeDelta)
	timeStr := stime.Format("20060102150405")
	return timeStr
}

func generateErr(frefix string, errs []error) error {
	if errs == nil {
		return nil
	}
	temp := []string{frefix}
	for _, err := range errs {
		temp = append(temp, err.Error())
	}
	return fmt.Errorf(strings.Join(temp, ":"))
}

type Client struct {
	reqAgent *gorequest.SuperAgent
	username string
	password string
	url      string
}

func New(url, username, password string) *Client {
	reqAgent := gorequest.New().
		SetDoNotClearSuperAgent(true).
		Timeout(HTTPRequestTimeOut).
		Retry(3, HTTPRequestRetryPeriod, http.StatusBadRequest, http.StatusInternalServerError)

	return &Client{
		reqAgent: reqAgent,
		username: username,
		password: password,
		url:      url,
	}
}

func (client *Client) indexUrl() string {
	return fmt.Sprintf("%s/index.php", client.url)
}

func (client *Client) login() error {
	res, _, errs := client.reqAgent.Get(client.indexUrl()).
		Query("autologin=1&enter=Sign%20in").
		Query("name="+client.username).
		Query("password="+client.password).
		Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0").
		AppendHeader("Referer", client.indexUrl()).End()
	if errs != nil {
		return generateErr("", errs)
	}
	fmt.Println(res.Status)
	return nil
}

func (client *Client) SaveImage(itemid, fileName string) (string, error) {

	if err := client.login(); err != nil {
		return "", fmt.Errorf("ZabbixClient:SaveImage:GetCookies:%v", err)

	}
	timeStr := getTimeStr()
	_, byteBody, errs := client.reqAgent.Get(client.url + "/chart.php").
		Query("period=7200&width=600").
		Query("time+" + timeStr).
		Query("itemids=" + itemid).
		EndBytes()
	if errs != nil {
		return "", generateErr("ZabbixClient:SaveImage:RequestImage", errs)
	}
	if fileName == "" {
		fileName = filepath.Join("tmp", timeStr+".jpg")
	} else {
		fileName = filepath.Join("tmp", fileName+".jpg")
	}

	out, err := os.Create(fileName)
	if err != nil {
		return fileName, fmt.Errorf("ZabbixClient:SaveImage:CreateImage:%v", err)
	}
	defer out.Close()
	_, err = out.Write(byteBody)
	if err != nil {
		return fileName, fmt.Errorf("ZabbixClient:SaveImage:SaveImage:%v", err)
	}
	return fileName, nil
}
