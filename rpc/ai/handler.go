package main

import (
	ai "LiveLive/kitex_gen/livelive/ai"
	"LiveLive/kitex_gen/livelive/base"
	"LiveLive/kitex_gen/livelive/websocket/websocketservice"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"context"
)

// AIServiceImpl implements the last service interface defined in the IDL.
type AIServiceImpl struct {
	wsClient websocketservice.Client
}

// AnalyzeAudio implements the AIServiceImpl interface.
func (s *AIServiceImpl) AnalyzeAudio(ctx context.Context, req *ai.AnalyzeAudioReq) (resp *ai.AnalyzeAudioResp, err error) {
	//上传到阿里云oss
	url, err := UploadAudio(req.WavPath)
	if err != nil {
		res := &ai.AnalyzeAudioResp{
			Summary: "",
			BaseResp: &base.BaseResp{
				Code: 60001,
				Msg:  "上传失败：" + err.Error(),
			},
		}
		return res, nil
	}
	// 2. 提交任务
	taskId, err := SubmitToVolc(url)
	if err != nil {
		return &ai.AnalyzeAudioResp{
			BaseResp: &base.BaseResp{
				Code: 60002,
				Msg:  "提交任务失败：" + err.Error(),
			},
		}, nil
	}

	// 3. 查询结果（加上 retry 逻辑）
	var text string
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		text, err = QueryVolcResult(taskId)
		if err == nil {
			break
		}
	}
	if err != nil {
		return &ai.AnalyzeAudioResp{
			BaseResp: &base.BaseResp{
				Code: 60003,
				Msg:  "查询结果失败：" + err.Error(),
			},
		}, nil
	}

	//4. 进行ai总结

	summary, err := GenerateSummaryFromText(ctx, text)
	if err != nil {
		log.Printf("生成总结失败: %v", err)
		return &ai.AnalyzeAudioResp{
			BaseResp: &base.BaseResp{
				Code: 60004,
				Msg:  "生成总结失败: " + err.Error(),
			},
		}, nil
	}

	return &ai.AnalyzeAudioResp{
		Summary: summary,
		BaseResp: &base.BaseResp{
			Code: 0,
			Msg:  "OK",
		},
	}, nil
}

// 配置阿里云oss的一些配置
var (
	Endpoint        = "oss-cn-beijing.aliyuncs.com" // 根据你创建的 Bucket 所在地域调整
	AccessKeyID     = "LTAI5tReroSL953zyeKeCdhA"
	AccessKeySecret = "BQcKmnkzLPAD1JYu42f6bxTgfYxYdv"
	BucketName      = "sonwwall-livelive"
	OssBaseURL      = "https://sonwwall-livelive.oss-cn-beijing.aliyuncs.com/"
)

// UploadAudio uploads a local audio file to OSS and returns the public URL.
func UploadAudio(localFilePath string) (string, error) {
	// 创建 OSS 客户端
	client, err := oss.New(Endpoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		return "", fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	// 获取 Bucket 实例
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return "", fmt.Errorf("获取Bucket失败: %w", err)
	}

	// 获取文件名作为对象名称
	objectName := filepath.Base(localFilePath)

	// 上传本地文件到 OSS
	err = bucket.PutObjectFromFile(objectName, localFilePath)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	// 设置文件为公共可读
	err = bucket.SetObjectACL(objectName, oss.ACLPublicRead)
	if err != nil {
		return "", fmt.Errorf("设置公开读失败: %w", err)
	}

	// 拼接公开可访问链接
	publicURL := OssBaseURL + objectName
	return publicURL, nil
}

//封装提交任务函数

const (
	apiAppKey    = "6367434862"
	apiAccessKey = "0kVVXn6KkBUNNkTAVry7MFx5P4_rTQgE"
	resourceId   = "volc.bigasr.auc"
	submitURL    = "https://openspeech.bytedance.com/api/v3/auc/bigmodel/submit"
	queryURL     = "https://openspeech.bytedance.com/api/v3/auc/bigmodel/query"
)

type SubmitRequest struct {
	User struct {
		Uid string `json:"uid"`
	} `json:"user"`
	Audio struct {
		Format string `json:"format"`
		URL    string `json:"url"`
		Rate   int    `json:"rate,omitempty"`
	} `json:"audio"`
	Request struct {
		ModelName  string `json:"model_name"`
		EnableItn  bool   `json:"enable_itn"`
		EnablePunc bool   `json:"enable_punc"`
		ShowUtt    bool   `json:"show_utterances"`
	} `json:"request"`
}

func SubmitToVolc(url string) (string, error) {
	reqBody := SubmitRequest{}
	reqBody.User.Uid = "123456"
	reqBody.Audio.URL = url
	reqBody.Audio.Format = "wav"
	reqBody.Audio.Rate = 16000
	reqBody.Request.ModelName = "bigmodel"
	reqBody.Request.EnableItn = true
	reqBody.Request.EnablePunc = true
	reqBody.Request.ShowUtt = false

	bodyBytes, _ := json.Marshal(reqBody)

	requestId := uuid.New().String()
	req, _ := http.NewRequest("POST", submitURL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-App-Key", apiAppKey)
	req.Header.Set("X-Api-Access-Key", apiAccessKey)
	req.Header.Set("X-Api-Resource-Id", resourceId)
	req.Header.Set("X-Api-Request-Id", requestId)
	req.Header.Set("X-Api-Sequence", "-1")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Api-Status-Code") != "20000000" {
		msg := resp.Header.Get("X-Api-Message")
		return "", fmt.Errorf("提交失败: %s", msg)
	}

	return requestId, nil
}

//封装查询任务函数

type Response struct {
	Result struct {
		Text       string `json:"text"`
		Utterances []struct {
			Text      string `json:"text"`
			StartTime int    `json:"start_time"`
			EndTime   int    `json:"end_time"`
			Definite  bool   `json:"definite"`
			Words     []struct {
				Text          string `json:"text"`
				StartTime     int    `json:"start_time"`
				EndTime       int    `json:"end_time"`
				BlankDuration int    `json:"blank_duration"`
			} `json:"words"`
		} `json:"utterances"`
	} `json:"result"`

	AudioInfo struct {
		Duration int `json:"duration"`
	} `json:"audio_info"`
}

func QueryVolcResult(requestId string) (string, error) {
	req, _ := http.NewRequest("POST", queryURL, bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-App-Key", apiAppKey)
	req.Header.Set("X-Api-Access-Key", apiAccessKey)
	req.Header.Set("X-Api-Resource-Id", resourceId)
	req.Header.Set("X-Api-Request-Id", requestId)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Api-Status-Code") != "20000000" {
		msg := resp.Header.Get("X-Api-Message")
		return "", fmt.Errorf("查询失败: %s", msg)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("解析失败：", err)
		log.Println("响应内容：", string(body))
		return "", err
	}

	log.Println("识别结果文本：", res.Result.Text)
	return res.Result.Text, nil
}

// 创建模型
func createHuoShanChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Model:   "ep-20250528160851-c9czd",
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Region:  "cn-beijing",
		APIKey:  "8f51e4bd-de97-4ad2-87a6-a63f3f397a03", // 注意：生产环境请勿硬编码 APIKey
	})
	if err != nil {
		log.Fatalf("create chat model failed: %v", err)
	}
	return chatModel
}

// GenerateSummaryFromText 接收识别出的文本，返回总结内容
func GenerateSummaryFromText(ctx context.Context, text string) (string, error) {
	// 创建聊天模型
	llm := createHuoShanChatModel(ctx)

	// 构造摘要请求模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个经验丰富的教学助理，擅长总结课堂内容，帮助学生高效复习。请根据以下转写内容，提取课堂的主要内容、知识点与结论。"),
		schema.UserMessage("课堂内容如下：{transcript}\n\n请生成总结："),
	)

	// 构造 prompt 消息
	messages, err := template.Format(ctx, map[string]any{
		"transcript": text,
	})
	if err != nil {
		return "", fmt.Errorf("格式化 prompt 失败: %w", err)
	}

	// 调用模型生成总结
	result, err := llm.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("模型生成失败: %w", err)
	}

	return result.Content, nil
}

// ChatWithAI implements the AIServiceImpl interface.
func (s *AIServiceImpl) ChatWithAI(ctx context.Context, req *ai.ChatWithAIReq) (resp *ai.ChatWithAIResp, err error) {
	// TODO: Your code here...
	return
}
