package controller

import (
	model2 "ai-smart/model"
	"github.com/gin-gonic/gin"
)

type OpenAIController struct {
}

func NewOpenAIController() *OpenAIController {
	return &OpenAIController{}
}

func (h *OpenAIController) Completions(c *gin.Context, req *model2.BaseHeaderReq) (rsp *model2.BaseResponse) {

	return nil
}
