package front

import (
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RtrwControllerType struct {
	Gets                     func(*gin.Context)
	GetByUUID                func(*gin.Context)
	GetByPosition            func(*gin.Context)
	GetRegenciesByProvinceId func(*gin.Context)
}

func RtrwController() RtrwControllerType {
	RtrwService := service.RtrwService{
		DB: gorm_support.DB,
	}
	RtrwGeojsonService := service.RtrwGeojsonService{
		DB: gorm_support.DB,
	}
	getRtrws := func(ctx *gin.Context) {
		reg_province_id := ctx.DefaultQuery("reg_province_id", "51")
		latitude := ctx.DefaultQuery("lat", "")
		longitude := ctx.DefaultQuery("lng", "")
		rtrw_datas := []model.RtrwType{}
		rtrwDAtaDb := RtrwService.Gets()
		if reg_province_id == "" {
			_reg_province_id, err := GetNearProvinceByPosition(latitude, longitude)
			if err != nil {
				ctx.JSON(400, gin.H{
					"status":      "error",
					"status_code": 400,
					"return":      err.Error(),
				})
			}
			reg_province_id = strconv.Itoa(int(*_reg_province_id))
		}
		fmt.Println("reg_province_id:: ", reg_province_id)
		err := rtrwDAtaDb.
			Preload("Reg_province").
			Preload("Reg_regency").
			Preload("Reg_district").
			Preload("Rtrw_mbtiles").Where("reg_province_id = ?", reg_province_id).Where("status = ?", "active").Find(&rtrw_datas).Error
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
			"return":      rtrw_datas,
			"status":      "success",
			"status_code": 200,
		})
	}

	getByPosition := func(ctx *gin.Context) {
		latlng := ctx.Param("latlng")
		latlngArr := strings.Split(latlng, ",")
		lat := latlngArr[0]
		lng := latlngArr[1]

		rtrwGeojson := []model.RtrwGeojsonView{}
		rtrwGeoDb := RtrwGeojsonService.Gets()
		err := rtrwGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Select("rtrw_geojson.*, ST_AsGeoJSON(geojson) as geojson").Find(&rtrwGeojson).Error
		if err != nil {
			if err != nil {
				log.Println(err)
				ctx.JSON(200, gin.H{
					"type":        "rtrw",
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
			"type":        "rtrw",
			"return":      rtrwGeojson,
			"status":      "success",
			"status_code": 200,
			"lat":         lat,
			"lng":         lng,
		})
	}

	getRtrwByUUID := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getRtrwByUUID endpoint"})
	}

	getRegenciesByPronvinceId := func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "getRtrwByUUID endpoint"})
	}

	return RtrwControllerType{
		Gets:                     getRtrws,
		GetByUUID:                getRtrwByUUID,
		GetByPosition:            getByPosition,
		GetRegenciesByProvinceId: getRegenciesByPronvinceId,
	}
}
