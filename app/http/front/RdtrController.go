package front

import (
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type RdtrControllerType struct {
	Gets                     func(*gin.Context)
	GetByUUID                func(*gin.Context)
	GetByPosition            func(*gin.Context)
	GetRegenciesByProvinceId func(*gin.Context)
}

func RdtrController() RdtrControllerType {
	RdtrService := service.RdtrService{
		DB: gorm_support.DB,
	}
	RdtrGeojsonService := service.RdtrGeojsonService{
		DB: gorm_support.DB,
	}
	getRdtrs := func(ctx *gin.Context) {
		reg_province_id := ctx.DefaultQuery("reg_province_id", "")
		rdtr_datas := []model.RdtrType{}
		rdtrDAtaDb := RdtrService.Gets()
		fmt.Println("reg_province_id:: ", reg_province_id)
		err := rdtrDAtaDb.
			Preload("Reg_province").
			Preload("Reg_regency").
			Preload("Reg_district").
			Preload("Rdtr_mbtiles").Where("reg_province_id = ?", reg_province_id).Where("status = ?", "active").Find(&rdtr_datas).Error
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

		rdtrGeojson := []model.RdtrGeojsonView{}
		rdtrGeoDb := RdtrGeojsonService.Gets()
		err := rdtrGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Select("rdtr_geojson.*, ST_AsGeoJSON(geojson) as geojson").Find(&rdtrGeojson).Error
		if err != nil {
			if err != nil {
				log.Println(err)
				ctx.JSON(200, gin.H{
					"type":        "rdtr",
					"status":      "success",
					"status_code": 200,
					"return":      make([]map[string]interface{}, 0),
					"lat":         lat,
					"lng":         lng,
				})
				return
			}
		}
		ctx.JSON(200, gin.H{
			"type":        "rdtr",
			"return":      rdtrGeojson,
			"status":      "success",
			"status_code": 200,
			"lat":         lat,
			"lng":         lng,
		})
	}

	getRdtrByUUID := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getRdtrByUUID endpoint"})
	}

	getRegenciesByPronvinceId := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getRdtrByUUID endpoint"})
	}

	return RdtrControllerType{
		Gets:                     getRdtrs,
		GetByUUID:                getRdtrByUUID,
		GetByPosition:            getByPosition,
		GetRegenciesByProvinceId: getRegenciesByPronvinceId,
	}
}
