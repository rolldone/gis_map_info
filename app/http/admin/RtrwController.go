package admin

import "github.com/gin-gonic/gin"

func GetRtrws(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func GetRtrwByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func AddRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func UpdateRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}

func DeleteRtrw(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
}
