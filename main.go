package main

import (
	"net/http"

	"io/ioutil"

	"fmt"

	"net"

	"strings"

	"github.com/emacsist/alfred3-work-flow/utils"
	alfred "github.com/ruedap/go-alfred"
)

const (
	URL       string = "http://ip.cn/"
	URL_QUERY string = "http://ip.cn/index.php?ip=%s"
)

func main() {
	query := utils.GetQuery()
	response := alfred.NewResponse()
	if len(strings.TrimSpace(query)) == 0 {
		body := get(URL)
		dealResponse(body, response)
	} else {
		ip := net.ParseIP(query)
		if ip.To4() == nil {
			utils.AdItem(response, fmt.Sprintf("%v 不是有效的IP地址", query))
		} else {
			body := get(fmt.Sprintf(URL_QUERY, query))
			dealResponse(body, response)
		}
	}
	utils.AlfredOutput(response)
}

func get(url string) string {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "curl/7.51.0")
	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err.Error()
	} else {
		body, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return err.Error()
		} else {
			return string(body)
		}
	}
}

func dealResponse(httpBody string, alfredResponse *alfred.Response) {
	utils.AdItem(alfredResponse, "查询的结果是: "+httpBody)
	title := httpBody
	var ipString string
	var location string
	var suffx string
	isCurrentIP := true
	_, err := fmt.Sscanf(string(title), "当前 IP：%s 来自：%s %s", &ipString, &location, &suffx)
	if err != nil {
		_, err := fmt.Sscanf(string(title), "IP：%s 来自：%s %s", &ipString, &location, &suffx)
		if err != nil {
			utils.AdItem(alfredResponse, string(err.Error()))
			return
		}
		isCurrentIP = false
	}
	// fmt.Printf("%v => %v\n", ipString, location)
	if isCurrentIP {
		utils.AdItem(alfredResponse, "当前IP: "+ipString)
	} else {
		utils.AdItem(alfredResponse, "IP: "+ipString)
	}
	utils.AdItem(alfredResponse, "来自: "+location+" "+suffx)
}
