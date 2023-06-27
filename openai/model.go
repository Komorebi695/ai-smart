package openai

import "os"

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
