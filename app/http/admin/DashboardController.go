package admin

import "github.com/gin-gonic/gin"

type GetRdtrPeriodeType func(ctx *gin.Context)
type GetRtrwPeriodeType func(ctx *gin.Context)

type DashboardType struct {
	GetRdtrPeriode GetRdtrPeriodeType
	GetRtrwPeriode GetRtrwPeriodeType
}

func Dashboard() DashboardType {
	dashboard := DashboardType{
		GetRdtrPeriode: getRdtrPeriode(),
		GetRtrwPeriode: getRtrwPeriode(),
	}
	return dashboard
}

func getRdtrPeriode() GetRdtrPeriodeType {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "GetRdtrPeriode endpoint"})
	}
}

func getRtrwPeriode() GetRtrwPeriodeType {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "GetRtrwPeriode endpoint"})
	}
}
