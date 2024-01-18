package front

import (
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
}

func (c *LocationController) GetsProvinceDistincExist(ctx *gin.Context) {
	locationService := service.LocationService{
		DB: gorm_support.DB,
	}

	regProvinceDatas := []model.RegProvinceView{}
	locationServiceDB := locationService.GetsProvince()
	err := locationServiceDB.
		Where("EXISTS(SELECT * FROM rdtr where rdtr.reg_province_id = reg_province.id)").
		Or("EXISTS(SELECT * FROM rtrw where rtrw.reg_province_id = reg_province.id)").
		Find(&regProvinceDatas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      regProvinceDatas,
		"status":      "success",
		"status_code": 200,
	})
}
