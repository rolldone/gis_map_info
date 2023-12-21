package admin

import "github.com/gin-gonic/gin"

func GetZLPs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "getZLPs endpoint"})
}

func GetZLPByUUId(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func AddZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func UpdateZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}

func DeleteZLP(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "GetZLPs endpoint"})
}
