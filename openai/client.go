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
	ApiKey                 = "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"
)

// HttpJsonSend 发送json数据的http请求
func HttpJsonSend(method, reqUrl string, param interface{}) ([]byte, error) {
	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	log.Printf("request json string :\nurl:%v\n%v", reqUrl, string(b))

	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ApiKey))

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

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response error status:%v", response.Status))
	}
	log.Printf("response status:%v", response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// HttpMultiPartSend 发送 文件+文本 的http请求
func HttpMultiPartSend(method, url string, file *os.File, param interface{}) ([]byte, error) {
	file, err := os.Open("filePath")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	filePart, err := writer.CreateFormFile("image", "filePath")
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
	req, err := http.NewRequest(method, url, body)
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

// AudioTranslations 将音频翻译成英语文本。
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
