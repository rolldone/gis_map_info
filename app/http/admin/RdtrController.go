package admin

import (
	"github.com/gin-gonic/gin"
)

type RdtrController struct{}

func (a *RdtrController) GetRdtrs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func (a *RdtrController) GetRdtrByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func (a *RdtrController) AddRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func (a *RdtrController) UpdateRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func (a *RdtrController) DeleteRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}
