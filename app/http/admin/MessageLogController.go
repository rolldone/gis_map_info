package admin

import (
	"log"
	"net/http"
	"time"

	Helper "gis_map_info/app/helper"
	Model "gis_map_info/app/model"
	Service "gis_map_info/app/service"
	"gis_map_info/support/gorm_support"

	"github.com/gin-gonic/gin"
)

func MessageLogControllerConstruct() MessageLogController {
	gg := MessageLogController{}
	return gg
}

type MessageLogController struct{}

type Filter_MessageLogController struct {
	Cache_filter *string `json:"cache_filter,omitempty"`
	Uuid         string  `json:"uuid,omitempty"`
	Take         int     `json:"take,omitempty"`
	Skip         int     `json:"skip,omitempty"`
	From_at      *string `json:"from_at,omitempty"`
}

func (c *MessageLogController) Gets(ctx *gin.Context) {
	filter := Filter_MessageLogController{
		Take: 20,
	}
	err := ctx.ShouldBindJSON(&filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		log.Println(err)
		return
	}
	jobLogDatas := []Model.MessageLog{}
	messageLogService := Service.AsynqJobService{
		DB: gorm_support.DB,
	}
	jobLogDB := messageLogService.Gets()

	jobLogDB.Limit(filter.Take)
	jobLogDB.Offset(filter.Skip * filter.Take)

	if filter.From_at != nil {
		parsedTime, err := time.Parse(time.RFC3339, *filter.From_at)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":      "error",
				"status_code": http.StatusBadRequest,
				"return":      err.Error(),
			})
			log.Println(err)
			return
		}
		jobLogDB.Where("created_at > ?", parsedTime)
	}
	jobLogDB.Where("uuid = ?", filter.Uuid)
	if filter.From_at != nil {
		jobLogDB.Order("created_at ASC")
	} else {
		jobLogDB.Order("created_at DESC")
	}
	err = jobLogDB.Find(&jobLogDatas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		log.Println(err)
		return
	}

	if filter.From_at == nil {
		jobLogDatas = Helper.ReverseArray(jobLogDatas)
	}

	// cache_filter := helper.SaveParameter(ctx, filter, time.Duration(time.Minute*10))
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      jobLogDatas,
		// "cache_filter": cache_filter,
	})
}
