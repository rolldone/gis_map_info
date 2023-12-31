package admin

import (
	"fmt"
	Service "gis_map_info/app/service"
	"strconv"

	"github.com/gin-gonic/gin"

	Model "gis_map_info/app/model"
)

type RegLocationController struct {
}

func (c *RegLocationController) GetProvinces(ctx *gin.Context) {
	tx := Model.DB
	regLocationService := Service.RegLocationService{}
	regLocationService.DB = tx
	regProvinceDatas, err := regLocationService.GetProvinces()
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      regProvinceDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *RegLocationController) GetRegencies(ctx *gin.Context) {
	tx := Model.DB
	regLocationService := Service.RegLocationService{
		DB: tx,
	}
	province_id, err := strconv.Atoi(ctx.Query("reg_province_id"))
	if err != nil {
		fmt.Println("GetRegencies Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	regRegencyDatas, err := regLocationService.GetRegenciesByProvinceId(province_id)
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      regRegencyDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *RegLocationController) GetDistricts(ctx *gin.Context) {
	tx := Model.DB
	regLocationService := Service.RegLocationService{
		DB: tx,
	}
	regency_id, err := strconv.Atoi(ctx.Query("reg_regency_id"))
	if err != nil {
		fmt.Println("GetDistricts Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	regRegencyDatas, err := regLocationService.GetDistrictsByRegencyId(regency_id)
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      regRegencyDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *RegLocationController) GetVillages(ctx *gin.Context) {
	tx := Model.DB
	regLocationService := Service.RegLocationService{
		DB: tx,
	}
	district_id, err := strconv.Atoi(ctx.Query("reg_district_id"))
	if err != nil {
		fmt.Println("GetVillages Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	regVillageDatas, err := regLocationService.GetVillagesByDistrictId(district_id)
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      regVillageDatas,
		"status":      "success",
		"status_code": 200,
	})
}
