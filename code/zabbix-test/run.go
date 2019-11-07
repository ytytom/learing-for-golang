package zabbix

import (
	"testing"
)

func get_client() *Client {
	return New("http://172.20.13.203/zabbix", "Admin", "zabbix")
}

func Test_downloadImage(t *testing.T) {
	client := get_client()
	_, err := client.SaveImage("69618", "")
	if err != nil {
		t.Fatal(err)
	}
}
