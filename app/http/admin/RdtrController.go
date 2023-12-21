package admin

import (
	"github.com/gin-gonic/gin"
)

func GetRdtrs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func GetRdtrByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func AddRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func UpdateRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}

func DeleteRdtr(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
}
