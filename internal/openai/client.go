package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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
	ContentTypeJson        = "application/json"
	ContentTypeMultipart   = "multipart/form-data"
	Temperature            = 0.8
)

type CompletionsReq struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"` // 较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定。
	//N           int     `json:"n" default:"1"`
	//TopP        int     `json:"top_p"`
}

type CompletionsRsp struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int                 `json:"created"`
	Model   string              `json:"model"`
	Choices []CompletionsChoice `json:"choices"`
	Usage   Usage               `json:"usage"`
	Error   Error               `json:"error"`
}

type CompletionsChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type ChatReq struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float32       `json:"temperature"`
}

type ChatMessage struct {
	Role    string `json:"role"` // system user assistant
	Content string `json:"content"`
}

type ChatRsp struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Choices []ChatChoice `json:"choices"`
	Usage   Usage        `json:"usage"`
	Error   Error        `json:"error"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type EditsReq struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Instruction string  `json:"instruction"`
	Temperature float32 `json:"temperature" default:"0.7"` // 较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定。
	//N           int     `json:"n" default:"1"`
	//TopP        int     `json:"top_p" default:"1"`
}

type EditsRsp struct {
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   Usage                    `json:"usage"`
	Error   Error                    `json:"error"`
}

type ImagesGenReq struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n" default:"1"` // 响应次数
	Size   string `json:"size" default:"1024x1024"`
}

type ImagesGenRsp struct {
	Created int   `json:"created"`
	Data    []Url `json:"data"`
	Error   Error `json:"error"`
}

type Url struct {
	Url string `json:"url"`
}

type ImagesEditsReq struct {
	Image  *os.File `json:"image"` // 必须是有效的 PNG 文件，小于 4MB，并且是正方形。
	Prompt string   `json:"prompt"`
	N      int      `json:"n" default:"1"` // 响应次数
	Size   string   `json:"size" default:"1024x1024"`
}

type ImagesVarReq struct {
	Image *os.File `json:"image"`
	N     int      `json:"n" default:"1"`
	Size  string   `json:"size" default:"1024x1024"`
}

type AudioReq struct {
	File  *os.File `json:"file"`
	Model string   `json:"model"`
}

type AudioRsp struct {
	Text  string `json:"text"`
	Error Error  `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// HttpJsonSend 发送http请求
func HttpJsonSend(method, reqUrl string, param interface{}) ([]byte, error) {
	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	log.Printf("request json string :\nurl:%v\n%v", reqUrl, string(b))

	apiKey := ""
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	proxyURL, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		return nil, err
	}
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	log.Printf("response status:%v", response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpMultiPartSend(method, reqUrl string, param interface{}) ([]byte, error) {
	file, err := os.Open("filePath")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	filePart, err := writer.CreateFormFile("file", "filePath")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		return nil, err
	}

	// 添加文本字段
	textPart, err := writer.CreateFormField("text")
	if err != nil {
		return nil, err
	}
	_, err = textPart.Write([]byte(""))
	if err != nil {
		return nil, err
	}

	// 写入结束边界
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", "url", body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送HTTP请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("上传请求失败，状态码：%d", resp.StatusCode)
	}

	return nil, nil
}

// Completions 为提供的提示和参数创建补全。
func Completions(model, msg string, temperature float32) (string, error) {
	req := CompletionsReq{
		Model:       model,
		Prompt:      msg,
		MaxTokens:   2048,
		Temperature: temperature,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, CompletionsApi), req)
	if err != nil {
		return "", err
	}
	var data CompletionsRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data.Error.Type) > 0 {
		return "", errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}
	var reply string
	if len(data.Choices) > 0 {
		reply = data.Choices[0].Text
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// Chat 给定描述对话的消息列表，模型将返回响应。
func Chat(model string, msg []string, temperature float32) (string, error) {
	var msgList []ChatMessage
	for _, v := range msg {
		tmp := ChatMessage{
			Role:    "user",
			Content: v,
		}
		msgList = append(msgList, tmp)
	}
	req := ChatReq{
		Model:       model,
		Messages:    msgList,
		Temperature: temperature,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, ChatApi), req)
	if err != nil {
		return "", err
	}
	var data ChatRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data.Error.Type) > 0 {
		return "", errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}
	var reply string
	if len(data.Choices) > 0 {
		reply = strings.ReplaceAll(data.Choices[0].Message.Content, "\n", "")
	}

	return reply, nil
}

// Edits 给定提示和指令，模型将返回提示的编辑版本。
func Edits(model, input, instruction string, temperature float32) (string, error) {
	req := EditsReq{
		Model:       model,
		Input:       input,
		Instruction: instruction,
		Temperature: temperature,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, EditApi), req)
	if err != nil {
		return "", err
	}

	var data EditsRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data.Error.Type) > 0 {
		return "", errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}
	var reply string
	if len(data.Choices) > 0 {
		for _, v := range data.Choices {
			reply = v["text"].(string)
			break
		}
	}

	return strings.ReplaceAll(reply, "\n", ""), nil
}

// ImagesGenerations 创建给定提示的图像。
func ImagesGenerations(prompt, size string, n int) ([]Url, error) {
	req := ImagesGenReq{
		Prompt: prompt,
		N:      n,
		Size:   size,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, ImagesGenerationsApi), req)
	if err != nil {
		return nil, err
	}
	var data ImagesGenRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if len(data.Error.Type) > 0 {
		return data.Data, errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}

	return data.Data, nil
}

// ImagesEdits 在给定原始图像和提示的情况下创建编辑或扩展的图像。 image 必须是有效的 PNG 文件，小于 4MB，并且是正方形。
func ImagesEdits(image *os.File, prompt, size string, n int) ([]Url, error) {
	req := ImagesEditsReq{
		Image:  image,
		Prompt: prompt,
		N:      n,
		Size:   size,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, ImagesEditsApi), req)
	if err != nil {
		return nil, err
	}
	var data ImagesGenRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if len(data.Error.Type) > 0 {
		return data.Data, errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}

	return data.Data, nil
}

// ImagesVariations 创建给定图像的变体。
func ImagesVariations(image *os.File, n int, size string) ([]Url, error) {
	req := ImagesVarReq{
		Image: image,
		N:     n,
		Size:  size,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, ImagesVariationsApi), req)
	if err != nil {
		return nil, err
	}
	var data ImagesGenRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if len(data.Error.Type) > 0 {
		return data.Data, errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}

	return data.Data, nil
}

// AudioTranscriptions 将音频转换为文本。
func AudioTranscriptions(model string, file *os.File) (string, error) {
	return audio(model, file, AudioTranscriptionsApi)
}

// AudioTranslations 将音频翻译成英语。
func AudioTranslations(model string, file *os.File) (string, error) {
	return audio(model, file, AudioTranslationsApi)
}

// Audio 将音频转换为文本。 格式之一：mp3、mp4、mpeg、mpga、m4a、wav 或 webm
func audio(model string, file *os.File, audioType string) (string, error) {
	audios := AudioReq{
		Model: model,
		File:  file,
	}
	body, err := HttpJsonSend(http.MethodPost, fmt.Sprintf("%s%s", BASEURL, audioType), audios)
	if err != nil {
		return "", err
	}
	var data AudioRsp
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data.Error.Type) > 0 {
		return "", errors.New(fmt.Sprintf("send request error! type:%v. message:%v", data.Error.Type, data.Error.Message))
	}

	return data.Text, nil
}
