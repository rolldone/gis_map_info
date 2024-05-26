package admin

import (
	"gis_map_info/app/helper"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func UserControllerConstruct() UserController {
	gg := UserController{}
	return gg
}

func (c *UserController) Add(ctx *gin.Context) {
	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)
	gg := struct {
		Name                  string `json:"name"`
		Username              string `json:"username"`
		Email                 string `json:"email"`
		Password              string `json:"password"`
		Password_confirm      string `json:"password_confirm"`
		Password_confirmation string `json:"password_confirmation,omitempty"`
	}{}
	err := ctx.BindJSON(&gg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	addPayload := service.AddPayload_UserService{
		Name:     gg.Name,
		Username: gg.Username,
		Email:    gg.Email,
		Password: &gg.Password,
		Status:   userService.GetStatus().ACTIVE,
	}
	jobData, err := userService.Add(addPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      jobData,
	})
	ctx.Done()
}

func (c *UserController) Update(ctx *gin.Context) {
	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)
	gg := struct {
		Name                  string `json:"name"`
		Username              string `json:"username"`
		Email                 string `json:"email"`
		Password              string `json:"password,omitempty"`
		Password_confirmation string `json:"password_confirmation,omitempty"`
		Status                string `json:"status"`
		Uuid                  string `json:"uuid"`
	}{}
	err := ctx.BindJSON(&gg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	updatePayload := service.UpdatePayload_UserService{
		Uuid: gg.Uuid,
		AddPayload_UserService: service.AddPayload_UserService{
			Name:                  gg.Name,
			Username:              gg.Username,
			Email:                 gg.Email,
			Password:              &gg.Password,
			Password_confirmation: &gg.Password_confirmation,
			Status:                gg.Status,
		},
	}
	jobData, err := userService.Update(updatePayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      jobData,
	})
	ctx.Done()
}

func (c *UserController) Delete(ctx *gin.Context) {
	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)
	gg := struct {
		Uuids []string `json:"uuids"`
	}{}
	err := ctx.BindJSON(&gg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	err = userService.Delete(gg.Uuids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      "deleted",
	})
	ctx.Done()
}

type Filter_user struct {
	Take         int      `json:"take,omitempty"`
	Skip         int      `json:"skip,omitempty"`
	Search       string   `json:"search,omitempty"`
	Uuids        []string `json:"uuids,omitempty"`
	Cache_filter *string  `json:"cache_filter,omitempty"`
}

func (c *UserController) Gets(ctx *gin.Context) {
	DB := gorm_support.DB
	userService := service.UserServiceConstruct(DB)
	filter := Filter_user{
		Take: 20,
		Skip: 0,
	}
	err := ctx.BindJSON(&filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}

	if filter.Cache_filter != nil {
		err := helper.GetParameter(ctx, *filter.Cache_filter, &filter)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":      "error",
				"status_code": http.StatusBadRequest,
				"return":      err.Error(),
			})
			log.Println(err)
		}
	}

	cache_filter := helper.SaveParameter(ctx, filter, time.Duration(time.Minute*15))
	filter.Cache_filter = &cache_filter

	userDatas := []model.UserView{}
	userDataDB := userService.Gets()
	userDataDB.Limit(filter.Take).Offset((filter.Skip - 1) * filter.Take)
	userDataDB.Order("updated_at DESC")
	err = userDataDB.Find(&userDatas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      userDatas,
		"filter":      filter,
	})
	ctx.Done()
}

func (c *UserController) GetByUUID(ctx *gin.Context) {
	DB := gorm_support.DB
	userService := service.UserServiceConstruct(DB)
	uuid := ctx.Param("uuid")
	userData := model.User{}
	userDataDB := userService.Gets()
	err := userDataDB.Where("uuid = ?", uuid).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      userData,
	})
	ctx.Done()
}
