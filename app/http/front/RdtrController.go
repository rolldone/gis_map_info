package front

import "github.com/gin-gonic/gin"

type RdtrControllerType struct {
	Gets      func(*gin.Context)
	GetByUUID func(*gin.Context)
}

func RdtrController() RdtrControllerType {

	getRdtrs := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
	}

	getRdtrByUUID := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "GetRdtrs endpoint"})
	}

	return RdtrControllerType{
		Gets:      getRdtrs,
		GetByUUID: getRdtrByUUID,
	}
}
