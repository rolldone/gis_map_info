package admin

import "github.com/gin-gonic/gin"

type ZLPController struct{}

func (a *ZLPController) GetZLPs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "getZLPs endpoint"})
}

func (a *ZLPController) GetZLPByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func (a *ZLPController) AddZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func (a *ZLPController) UpdateZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func (a *ZLPController) DeleteZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}
