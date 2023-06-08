package controller

import (
	"ai-smart/internal/model"
	"github.com/gin-gonic/gin"
)

type OpenAIController struct {
}

func NewOpenAIController() *OpenAIController {
	return &OpenAIController{}
}

func (h *OpenAIController) Completions(c *gin.Context, req *model.BaseHeaderReq) (rsp *model.BaseResponse) {

	return nil
}
