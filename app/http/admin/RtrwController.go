package admin

import "github.com/gin-gonic/gin"

type RtrwController struct{}

func (a *RtrwController) GetRtrws(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func (a *RtrwController) GetRtrwByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func (a *RtrwController) AddRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func (a *RtrwController) UpdateRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func (a *RtrwController) DeleteRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}
