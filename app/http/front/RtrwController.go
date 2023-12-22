package front

import "github.com/gin-gonic/gin"

type RtrwControllerType struct {
	Gets      func(*gin.Context)
	GetByUUID func(*gin.Context)
}

func RtrwController() RtrwControllerType {

	return RtrwControllerType{
		Gets: func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
		},
		GetByUUID: func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "GetRtrws endpoint"})
		},
	}
}
