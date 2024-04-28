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

type ZlpControllerType struct {
	Gets                     func(*gin.Context)
	GetByUUID                func(*gin.Context)
	GetByPosition            func(*gin.Context)
	GetRegenciesByProvinceId func(*gin.Context)
	GetsByZlpGroup           func(*gin.Context)
	GetPositionByZlpGroup    func(*gin.Context)
}

func ZlpController() ZlpControllerType {
	ZlpService := service.ZlpService{
		DB: gorm_support.DB,
	}
	ZlpGeojsonService := service.ZlpGeojsonService{
		DB: gorm_support.DB,
	}
	getZlps := func(ctx *gin.Context) {
		reg_province_id := ctx.DefaultQuery("reg_province_id", "51")
		zlp_datas := []model.ZlpType{}
		zlpDAtaDb := ZlpService.Gets()
		fmt.Println("reg_province_id:: ", reg_province_id)
		err := zlpDAtaDb.
			Preload("Reg_province").
			Preload("Reg_regency").
			Preload("Reg_district").
			Preload("Zlp_mbtiles").Where("reg_province_id = ?", reg_province_id).Where("status = ?", "active").Find(&zlp_datas).Error
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
			"return":      zlp_datas,
			"status":      "success",
			"status_code": 200,
		})
	}

	getByPosition := func(ctx *gin.Context) {
		latlng := ctx.Param("latlng")
		latlngArr := strings.Split(latlng, ",")
		lat := latlngArr[0]
		lng := latlngArr[1]

		zlpGeojson := []model.ZlpGeojsonView{}
		zlpGeoDb := ZlpGeojsonService.Gets()
		err := zlpGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Select("zlp_geojson.*, ST_AsGeoJSON(geojson) as geojson").Find(&zlpGeojson).Error
		if err != nil {
			if err != nil {
				log.Println(err)
				ctx.JSON(200, gin.H{
					"type":        "zlp",
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
			"type":        "zlp",
			"return":      zlpGeojson,
			"status":      "success",
			"status_code": 200,
			"lat":         lat,
			"lng":         lng,
		})
	}

	getZlpByUUID := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getZlpByUUID endpoint"})
	}

	getRegenciesByPronvinceId := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getZlpByUUID endpoint"})
	}

	getsByZlpGroup := func(ctx *gin.Context) {
		reg_province_id := ctx.DefaultQuery("reg_province_id", "51")
		zlp_group_datas := []model.ZlpGroupDistinctAssetView{}
		zlpService := service.ZlpService{
			DB: gorm_support.DB,
		}
		zlpGroupDB := zlpService.GetZlpGroups()
		fmt.Println("reg_province_id:: ", reg_province_id)
		zlpGroupDB.Preload("Mbtiles", "zlp_id IS NOT NULL").Joins("left join zlp_mbtile i on zlp_group.asset_key = i.asset_key").
			Where("i.reg_province_id = ?", reg_province_id)
		err := zlpGroupDB.Distinct("zlp_group.asset_key", "zlp_group.name").Find(&zlp_group_datas).Error
		if err != nil {
			ctx.JSON(400, gin.H{
				"status":      "error",
				"status_code": 400,
				"return":      err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"return":      zlp_group_datas,
			"status":      "success",
			"status_code": 200,
		})
	}

	GetPositionByZlpGroup := func(ctx *gin.Context) {
		latlng := ctx.Param("latlng")
		latlngArr := strings.Split(latlng, ",")
		lat := latlngArr[0]
		lng := latlngArr[1]

		zlpGeojson := []model.ZlpGeojsonView{}
		zlpGeoDb := ZlpGeojsonService.Gets()
		err := zlpGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Select("zlp_geojson.*, ST_AsGeoJSON(geojson) as geojson").Find(&zlpGeojson).Error
		if err != nil {
			if err != nil {
				log.Println(err)
				ctx.JSON(200, gin.H{
					"type":        "zlp",
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
			"type":        "zlp",
			"return":      zlpGeojson,
			"status":      "success",
			"status_code": 200,
			"lat":         lat,
			"lng":         lng,
		})
	}

	return ZlpControllerType{
		Gets:                     getZlps,
		GetByUUID:                getZlpByUUID,
		GetByPosition:            getByPosition,
		GetRegenciesByProvinceId: getRegenciesByPronvinceId,
		GetsByZlpGroup:           getsByZlpGroup,
		GetPositionByZlpGroup:    GetPositionByZlpGroup,
	}
}
