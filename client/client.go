package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	BASEURL                = "https://api.openai.com/v1"
	ModelName              = "text-davinci-003"
	CompletionsApi         = "/completions"
	ChatApi                = "/chat/completions"
	EditApi                = "/edits"
	ImagesGenerationsApi   = "/images/generations"   // 给定信息返回图像
	ImagesEditsApi         = "/images/edits"         // 在给定原始图像和提示的情况下创建编辑或扩展的图像
	ImagesVariationsApi    = "/images/variations"    // 在给定原始图像和提示的情况下创建编辑或扩展的图像
	AudioTranscriptionsApi = "/audio/transcriptions" // 将音频转录为输入语言
	AudioTranslationsApi   = "/audio/translations"   // 将音频翻译成英语
)

type BaseReq struct {
	Model            string  `json:"model"`
	MaxTokens        int     `json:"max_tokens" default:"2048"`
	N                int     `json:"n" default:"1"`
	Temperature      float32 `json:"temperature" default:"0.7"` // 较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定。
	TopP             int     `json:"top_p" default:"1"`
	FrequencyPenalty int     `json:"frequency_penalty" default:"0"`
	PresencePenalty  int     `json:"presence_penalty" default:"0"`
	User             string  `json:"user" default:""`
}

type CompletionsReq struct {
	BaseReq
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type CompletionsRsp struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Usage   map[string]interface{}   `json:"usage"`
	Choices []map[string]interface{} `json:"choices"`
}

type ChatReq struct {
	Model   string        `json:"model"`
	Message []ChatMessage `json:"message"`
	BaseReq
}

type ChatMessage struct {
	Role    string `json:"role"` // system user assistant
	Content string `json:"content"`
}

type ChatRsp struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type EditsReq struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Instruction string  `json:"instruction"`
	N           int     `json:"n" default:"1"`
	Temperature float32 `json:"temperature" default:"0.7"` // 较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定。
	TopP        int     `json:"top_p" default:"1"`
}

type EditsRsp struct {
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ImagesGenReq struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n" default:"1"`
	Size   string `json:"size" default:"1024x1024"`
}

type ImagesGenRsp struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	}
}

type ImagesVarReq struct {
	Image string `json:"image"`
	N     int    `json:"n" default:"1"`
	Size  string `json:"size" default:"1024x1024"`
}

// HttpMethodSend 发送http请求
func HttpMethodSend(method, reqUrl string, param interface{}) ([]byte, error) {
	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	apiKey := ""
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Completions 为提供的提示和参数创建补全。
func Completions(model, msg string) (string, error) {
	req := CompletionsReq{
		Model:  model,
		Prompt: msg,
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	body, err := HttpMethodSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, CompletionsApi), reqData)
	if err != nil {
		return "", err
	}

	var completionsData CompletionsRsp
	if err := json.Unmarshal(body, &completionsData); err != nil {
		return "", err
	}

	var reply string
	if len(completionsData.Choices) > 0 {
		for _, v := range completionsData.Choices {
			reply = v["text"].(string)
			break
		}
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// Chat 给定描述对话的消息列表，模型将返回响应。
func Chat(model string, msg []string) (string, error) {
	var msgList []ChatMessage
	for _, v := range msg {
		tmp := ChatMessage{
			Role:    "",
			Content: v,
		}
		msgList = append(msgList, tmp)
	}

	req := ChatReq{
		Model:   model,
		Message: msgList,
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	body, err := HttpMethodSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, ChatApi), reqData)
	if err != nil {
		return "", err
	}

	var chatData ChatRsp
	if err := json.Unmarshal(body, &chatData); err != nil {
		return "", err
	}

	var reply string
	if len(chatData.Choices) > 0 {
		for _, v := range chatData.Choices {
			reply = v["content"].(string)
			break
		}
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// Edits 给定提示和指令，模型将返回提示的编辑版本。
func Edits(model, input, instruction string) (string, error) {
	req := EditsReq{
		Model:       model,
		Input:       input,
		Instruction: instruction,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	body, err := HttpMethodSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, EditApi), reqData)
	if err != nil {
		return "", err
	}

	var editsData EditsRsp
	if err := json.Unmarshal(body, &editsData); err != nil {
		return "", err
	}

	var reply string
	if len(editsData.Choices) > 0 {
		for _, v := range editsData.Choices {
			reply = v["text"].(string)
			break
		}
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// Images 创建给定提示的图像。 或者 在给定原始图像和提示的情况下创建编辑或扩展的图像。
func Images(prompt string, n int, size, apiType string) (interface{}, error) {
	req := ImagesGenReq{
		Prompt: prompt,
		N:      n,
		Size:   size,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body, err := HttpMethodSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, apiType), reqData)
	if err != nil {
		return nil, err
	}
	var data ImagesGenRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data.Data, nil
}

// ImagesVariations 创建给定图像的变体。
func ImagesVariations(image string, n int, size, apiType string) (interface{}, error) {
	req := ImagesVarReq{
		Image: image,
		N:     n,
		Size:  size,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body, err := HttpMethodSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, apiType), reqData)
	if err != nil {
		return nil, err
	}
	var data ImagesGenRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data.Data, nil
}
