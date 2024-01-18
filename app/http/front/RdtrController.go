package front

import (
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"strings"

	"github.com/gin-gonic/gin"
)

type RdtrControllerType struct {
	Gets          func(*gin.Context)
	GetByUUID     func(*gin.Context)
	GetByPosition func(*gin.Context)
}

func RdtrController() RdtrControllerType {
	RdtrService := service.RdtrService{
		DB: gorm_support.DB,
	}
	RdtrGeojsonService := service.RdtrGeojsonService{
		DB: gorm_support.DB,
	}
	getRdtrs := func(ctx *gin.Context) {
		rdtr_datas := []model.RdtrType{}
		rdtrDAtaDb := RdtrService.Gets()
		err := rdtrDAtaDb.Preload("Rdtr_mbtiles").Where("status = ?", "active").Find(&rdtr_datas).Error
		if err != nil {
			if err != nil {
				ctx.JSON(400, gin.H{
					"status":      "error",
					"status_code": 400,
					"return":      err.Error(),
				})
				return
			}
		}
		ctx.JSON(200, gin.H{
			"return":      rdtr_datas,
			"status":      "success",
			"status_code": 200,
		})
	}

	getByPosition := func(ctx *gin.Context) {
		latlng := ctx.Param("latlng")
		latlngArr := strings.Split(latlng, ",")
		lat := latlngArr[0]
		lng := latlngArr[1]

		rdtrGeojson := model.RdtrGeojsonView{}
		rdtrGeoDb := RdtrGeojsonService.Gets()
		err := rdtrGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Select("rdtr_geojson.*, ST_AsGeoJSON(geojson) as geojson").First(&rdtrGeojson).Error
		if err != nil {
			if err != nil {
				ctx.JSON(400, gin.H{
					"status":      "error",
					"status_code": 400,
					"return":      err.Error(),
				})
				return
			}
		}
		ctx.JSON(200, gin.H{
			"return":      rdtrGeojson,
			"status":      "success",
			"status_code": 200,
		})
	}

	getRdtrByUUID := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getRdtrByUUID endpoint"})
	}

	return RdtrControllerType{
		Gets:          getRdtrs,
		GetByUUID:     getRdtrByUUID,
		GetByPosition: getByPosition,
	}
}
