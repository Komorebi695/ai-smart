package controller

import (
	model "ai-smart/model"
	"github.com/gin-gonic/gin"
)

type OpenAIController struct{}

func NewOpenAIController() *OpenAIController {
	return &OpenAIController{}
}

// Completions
// @Summary  Completions
// @Description Completions
// @Router /v1/chat/completions [POST]
// @Tags FleetDriver
// @Param data body model.BaseHeaderReq true "参数data"
// @Success 200 {object} model.BaseResponse
func (h *OpenAIController) Completions(c *gin.Context, req *model.BaseHeaderReq) (rsp *model.BaseResponse) {

	return nil
}
