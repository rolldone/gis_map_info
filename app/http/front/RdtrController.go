package front

import (
	"encoding/json"
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RdtrControllerType struct {
	Gets                     func(*gin.Context)
	GetByUUID                func(*gin.Context)
	GetByPosition            func(*gin.Context)
	GetRegenciesByProvinceId func(*gin.Context)
	GetRdtrRtrwPos           func(*gin.Context)
}

func GetNearProvinceByPosition(latitude string, longitude string) (*int64, error) {
	regLocationService := service.RegLocationService{
		DB: gorm_support.DB,
	}
	latitude_64, err := strconv.ParseFloat(latitude, 64)
	if err != nil {

		return nil, err
	}
	longitude_64, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, err
	}
	regProvince_datas, err := regLocationService.GetNearProvinceByLocation(latitude_64, longitude_64, 1, 10000)
	if err != nil {
		return nil, err
	}
	if len(regProvince_datas) == 0 {
		return nil, fmt.Errorf("there is no province near your position")
	}
	fmt.Println("regProvince_datas :: ", regProvince_datas)
	return &regProvince_datas[0].Id, nil
}

func RdtrController() RdtrControllerType {
	RdtrService := service.RdtrService{
		DB: gorm_support.DB,
	}
	RdtrGeojsonService := service.RdtrGeojsonService{
		DB: gorm_support.DB,
	}
	getRdtrs := func(ctx *gin.Context) {
		reg_province_id := ctx.DefaultQuery("reg_province_id", "51")
		latitude := ctx.DefaultQuery("lat", "")
		longitude := ctx.DefaultQuery("lng", "")
		rdtr_datas := []model.RdtrType{}
		rdtrDAtaDb := RdtrService.Gets()
		if reg_province_id == "" {
			_reg_province_id, err := GetNearProvinceByPosition(latitude, longitude)
			if err != nil {
				ctx.JSON(400, gin.H{
					"status":      "error",
					"status_code": 400,
					"return":      err.Error(),
				})
				return
			}
			reg_province_id = strconv.Itoa(int(*_reg_province_id))
		}
		fmt.Println("reg_province_id:: ", reg_province_id)
		err := rdtrDAtaDb.
			Preload("Reg_province").
			Preload("Reg_regency").
			Preload("Reg_district").
			Preload("Rdtr_mbtiles").Where("reg_province_id = ?", reg_province_id).Where("status = ?", "active").
			Order("name ASC").
			Find(&rdtr_datas).Error
		if err != nil {
			ctx.JSON(400, gin.H{
				"status":      "error",
				"status_code": 400,
				"return":      err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"return":      rdtr_datas,
			"status":      "success",
			"status_code": 200,
		})
	}

	getByPosition := func(ctx *gin.Context) {
		latlng := ctx.Param("latlng")
		ids_query := ctx.Query("ids")

		// Declare a variable to hold the slice of strings (or int)
		var ids []int64
		log.Println("ids_query :: ", ids_query)
		if ids_query != "" {
			// Unmarshal the JSON string into a slice
			err := json.Unmarshal([]byte(ids_query), &ids)
			if err != nil {
				log.Printf("Error parsing ids: %v", err)
				ctx.JSON(400, gin.H{
					"error": "invalid JSON format",
				})
				return
			}
		}

		latlngArr := strings.Split(latlng, ",")
		lat := latlngArr[0]
		lng := latlngArr[1]

		rdtrGeojson := []model.RdtrGeojsonView{}
		rdtrGeoDb := RdtrGeojsonService.Gets()
		err := rdtrGeoDb.Where("ST_Within(ST_SetSRID(ST_MakePoint(?, ?), 4326), geojson)", lng, lat).Where("rdtr_id IN ?", ids).Select("rdtr_geojson.*, ST_AsGeoJSON(geojson) as geojson").Find(&rdtrGeojson).Error
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

	getRdtrRtrwPos := func(ctx *gin.Context) {

	}

	return RdtrControllerType{
		Gets:                     getRdtrs,
		GetByUUID:                getRdtrByUUID,
		GetByPosition:            getByPosition,
		GetRegenciesByProvinceId: getRegenciesByPronvinceId,
		GetRdtrRtrwPos:           getRdtrRtrwPos,
	}
}
