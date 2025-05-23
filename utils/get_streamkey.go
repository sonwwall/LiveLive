package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LiveGoResponse struct {
	ChannelKey string `json:"data"`
}

func GetStreamKey(room string) (string, error) {
	url := fmt.Sprintf("http://localhost:8090/control/get?room=%s", room)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求 LiveGo 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LiveGo 返回错误状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result LiveGoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return result.ChannelKey, nil
}
