package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main() {

	query := `sum by (namespace, pod, container) (rate(container_cpu_usage_seconds_total{namespace="kube-system", container="calico-kube-controllers"}[10m]))`
	encodedQuery := url.QueryEscape(query)
	now := time.Now()

	// 十分钟前
	tenMinutesAgo := now.Add(-5 * time.Hour)
	url := fmt.Sprintf("http://10.10.6.20:30759/api/v1/query_range?query=%s&start=%d&end=%d&step=300",
		encodedQuery, tenMinutesAgo.Unix(), now.Unix())
	// 发起 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 输出结果
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	return

}
