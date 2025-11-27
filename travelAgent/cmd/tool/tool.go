package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"travelAgent/cmd/model"
)

// getCity 模糊查询城市（内部函数）
func getCity(city string) *model.LocationResponse {
	url := "https://qd2tunyyy8.re.qweatherapi.com/geo/v2/city/lookup?location=" + city
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("X-QW-Api-Key", "80e8a0f704154b5896c7949d9029b4ac")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("查询城市编码出错", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("查询城市异常")
		return nil
	}

	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	var locationResp = &model.LocationResponse{}
	json.Unmarshal(bytes, locationResp)
	return locationResp
}

// GetWeather 通过调用和风天气API查询真实的天气信息
func GetWeather(city string) string {
	locationResp := getCity(city)
	if locationResp == nil {
		return ""
	}
	id := locationResp.Location[0].ID

	url := "https://qd2tunyyy8.re.qweatherapi.com/v7/weather/now?location=" + id
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("X-QW-Api-Key", "80e8a0f704154b5896c7949d9029b4ac")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("查询天气出错", err)
		return "错误:查询天气时遇到网络问题"
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("查询天气异常")
		return "查询天气异常"
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	var weather = model.WeatherResponse{}
	json.Unmarshal(bytes, &weather)
	temp := weather.Now.Temp
	text := weather.Now.Text
	return fmt.Sprintf("%s当前天气：%s，气温%s摄氏度", locationResp.Location[0].Name, text, temp)
}

// GetAttraction 根据城市和天气，使用Tavily Search API搜索并返回优化后的景点推荐
func GetAttraction(city, weather string) string {
	query := fmt.Sprintf("%s在%s天气下最值得去的旅游景点推荐及理由", city, weather)

	url := "https://api.tavily.com/search"
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	m := map[string]string{
		"query":          query,
		"include_answer": "basic",
	}
	marshal, _ := json.Marshal(m)
	request, _ := http.NewRequest("POST", url, bytes.NewBufferString(string(marshal)))
	request.Header.Set("Authorization", "Bearer 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")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return "错误:执行Tavily搜索时出现问题" + err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		return "错误:执行Tavily搜索报错"
	}

	defer resp.Body.Close()
	searchResponse := &model.SearchResponse{}
	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, searchResponse)

	// 如果有回答摘要 直接返回
	if searchResponse.Answer != nil {
		return *searchResponse.Answer
	}

	// 如果没有 需要格式化原始结果
	resultList := make([]string, 0)
	for _, result := range searchResponse.Results {
		resultList = append(resultList, fmt.Sprintf("%s:%s", result.Title, result.Content))
	}
	if len(resultList) == 0 {
		return "抱歉，没有找到相关的旅游景点推荐。"
	}

	return "根据搜索，为您找到以下信息:\n" + "\n" + strings.Join(resultList, "\n") + "\n"
}

// AvailableTools 可用工具列表（已废弃，请使用 agent 包中的 getAvailableTools）
var AvailableTools = map[string]interface{}{
	"get_weather":    GetWeather,
	"get_attraction": GetAttraction,
}
