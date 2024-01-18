package admin

import (
	"encoding/json"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AsynqJobController struct{}

func (c *AsynqJobController) GetsAsynqJob(ctx *gin.Context) {
	asynqJobService := service.AsynqJobService{
		DB: gorm_support.DB,
	}
	asyncJobDatas := []*model.AsyncJobView{}
	asyncJobDB := asynqJobService.Gets()
	err := asyncJobDB.Find(&asyncJobDatas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      asyncJobDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *AsynqJobController) GetAsynqJobByAppUuid(ctx *gin.Context) {
	asynqJobService := service.AsynqJobService{
		DB: gorm_support.DB,
	}
	asyncJobData := model.AsyncJobView{}
	asyncJobDB := asynqJobService.Gets()
	err := asyncJobDB.Where("uuid = ?", ctx.DefaultQuery("uuid", "")).First(&asyncJobData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      asyncJobData,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *AsynqJobController) GetAsynqJobByUUIDS(ctx *gin.Context) {
	uuids := []string{}
	err := json.Unmarshal([]byte(ctx.DefaultQuery("uuids", "[]")), &uuids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	asynqJobService := service.AsynqJobService{
		DB: gorm_support.DB,
	}
	asyncJobData := model.AsyncJobView{}
	asyncJobDB := asynqJobService.Gets()
	err = asyncJobDB.Where("uuid IN ?", []string(uuids)).First(&asyncJobData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"return":      asyncJobData,
		"status":      "success",
		"status_code": 200,
	})
}

func (c *AsynqJobController) DeleteAsynqJobByUUIDS(ctx *gin.Context) {
	uuids := []string{}
	err := json.Unmarshal([]byte(ctx.DefaultQuery("uuids", "[]")), &uuids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	asynqjobService := service.AsynqJobService{}
	asynqjobService.DeleteByUuids(uuids)
	ctx.JSON(200, gin.H{"message": "DeleteAsynqJobByUUIDS deleted"})
}
