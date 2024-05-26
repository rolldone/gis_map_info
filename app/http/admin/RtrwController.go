package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"gis_map_info/app/helper"
	Helper "gis_map_info/app/helper"
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"
	"gis_map_info/app/service"
	Service "gis_map_info/app/service"
	"gis_map_info/support/asynq_support"
	"gis_map_info/support/gorm_support"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/benpate/convert"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type RtrwController struct{}

func (a *RtrwController) GetRtrws(ctx *gin.Context) {
	var RtrwService = Service.RtrwService{}
	rtrwDB := RtrwService.Gets()
	var rtrwDatas = []Model.RtrwType{}
	// Fetch data from the database
	if err := rtrwDB.Find(&rtrwDatas).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rtrwDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RtrwController) GetRtrwsPaginate(ctx *gin.Context) {
	var RtrwService = Service.RtrwService{}
	rtrwDB := RtrwService.Gets()
	var rtrwDatas = []Model.RtrwType{}
	var page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	var limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	var offset = (page - 1) * limit

	// Fetch data from the database
	if err := rtrwDB.
		Preload("Reg_province").
		Preload("Reg_regency").
		Preload("Reg_district").
		Preload("Rtrw_groups").
		Preload("Rtrw_groups.Datas").
		Limit(limit).Offset(offset).Order("updated_at DESC").Find(&rtrwDatas).Error; err != nil {
		fmt.Println("Error - 29mvamfivm2 ", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rtrwDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RtrwController) GetRtrwById(ctx *gin.Context) {
	var RtrwService = Service.RtrwService{}
	id, err := strconv.Atoi(ctx.Param("id"))
	// Check for any conversion error
	if err != nil {
		fmt.Println("GetRtrwById Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Something wrong with parameter",
		})
		return
	}
	rtrwDB := RtrwService.GetById(id)
	var rtrwData = Model.RtrwType{}
	// Fetch data from the database
	if err := rtrwDB.
		Preload("Reg_province").
		Preload("Reg_district").
		Preload("Rtrw_groups", func(db *gorm.DB) *gorm.DB {
			return db.Select("rtrw_group.*, COALESCE((SELECT COUNT(*) FROM rtrw_file  WHERE rtrw_file.rtrw_group_id = rtrw_group.id AND rtrw_file.validated_at IS NULL),0) AS unvalidated, " +
				"COALESCE((SELECT COUNT(*) FROM rtrw_file WHERE rtrw_file.rtrw_group_id = rtrw_group.id AND rtrw_file.validated_at IS NOT NULL),0) as validated")
		}).
		Preload("Rtrw_mbtiles").
		Preload("Rtrw_groups.Datas").First(&rtrwData).Error; err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(404, gin.H{
			"status":      "error",
			"status_code": 404,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rtrwData,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RtrwController) AddRtrw(ctx *gin.Context) {
	var props struct {
		Name            string        `json:"name" validate:"required"`
		Place_string    string        `json:"place_string" validate:"required"`
		Reg_Province_id int64         `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64         `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64         `json:"reg_district_id"`
		Reg_Village_id  int64         `json:"reg_village_id"`
		Status          string        `json:"status"`
		Rtrw_groups     []interface{} `json:"rtrw_groups"`
		Rtrw_mbtiles    []interface{} `json:"rtrw_mbtiles"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	err := validate.Struct(props)
	if err != nil {
		var errArr []error
		for _, err := range err.(validator.ValidationErrors) {
			errMsg := fmt.Errorf("validation Error: Field %s has invalid value: %v", err.Field(), err.Value())
			errArr = append(errArr, errMsg)
		}
		err = errors.Join(errArr...)
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err,
		})
		return
	}
	tx := gorm_support.DB.Begin()
	var RtrwService = Service.RtrwService{
		DB: tx,
	}

	rtrwServiceAddData := RtrwService.RtrwServiceAddType
	rtrwServiceAddData.Name = props.Name
	rtrwServiceAddData.Place_string = props.Place_string
	rtrwServiceAddData.RegProvince_id = props.Reg_Province_id
	rtrwServiceAddData.RegRegency_id = props.Reg_Regency_id
	rtrwServiceAddData.RegDistrict_id = props.Reg_District_id
	rtrwServiceAddData.RegVillage_id = props.Reg_Village_id
	rtrwServiceAddData.Status = props.Status
	rtrwData, err := RtrwService.Add(rtrwServiceAddData)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error - 2mvcaisdfmv29 :", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	// rtrwGroups := props.Rtrw_groups
	// fmt.Println(rtrwGroups)
	// for i := 0; i < len(rtrwGroups); i++ {
	// 	rtrwGroupItem, _ := rtrwGroups[i].(map[string]interface{})
	// 	rtrwGroupItem["rtrw_id"] = rtrwData.Id
	// 	rtrwGroupData := RtrwService.RtrwGroupAddType
	// 	rtrwGroupData.Name = rtrwGroupItem["name"].(string)
	// 	rtrwGroupData.Rtrw_id = int64(rtrwGroupItem["rtrw_id"].(float64))
	// 	rtrwGroupData.Asset_key = rtrwGroupItem["asset_key"].(string)
	// 	rtrwGroupData.Uuid = uuid.NewString()
	// 	_rtrwGroupItem_properties, _ := Helper.GetValue(rtrwGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
	// 	rtrwGroupData.Properties = _rtrwGroupItem_properties
	// 	err = validate.Struct(rtrwGroupData)
	// 	if err != nil {
	// 		var errArr []error
	// 		var errValidators = err.(validator.ValidationErrors)
	// 		for i := 0; i < len(errValidators); i++ {
	// 			errValidItem := errValidators[i]
	// 			errMess := fmt.Errorf("validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
	// 			errArr = append(errArr, errMess)
	// 		}
	// 		err = errors.Join(errArr...)
	// 		break
	// 	}
	// 	_, err2 := RtrwService.AddGroup(rtrwGroupData)
	// 	if err2 != nil {
	// 		err = err2
	// 		break
	// 	}
	// }
	// if err != nil {
	// 	tx.Rollback()
	// 	fmt.Println("Error:", err)
	// 	ctx.JSON(400, gin.H{
	// 		"status":      "error",
	// 		"status_code": 400,
	// 		"return":      err.Error(),
	// 	})
	// 	return
	// }

	// rtrwMbtiles := props.Rtrw_mbtiles
	// // Delete rtrw mbtile by rtrw_id first
	// RtrwService.DeleteMbtileByRtrwId(int(rtrwData.Id))
	// for i := 0; i < len(rtrwMbtiles); i++ {
	// 	rtrwMbtileItem, _ := rtrwMbtiles[i].(map[string]interface{})
	// 	rtrwMbtileItem["rtrw_id"] = rtrwData.Id
	// 	int64MbtileId := int64(rtrwMbtileItem["id"].(float64))
	// 	RtrwService.RtrwMbtilePayload = struct {
	// 		Id        *int64
	// 		File_name string
	// 		Uuid      string
	// 		Rtrw_id   int64
	// 	}{
	// 		// For id Because pointer need vessel first
	// 		Id:        &int64MbtileId,
	// 		File_name: rtrwMbtileItem["file_name"].(string),
	// 		Uuid:      rtrwMbtileItem["uuid"].(string),
	// 		Rtrw_id:   rtrwMbtileItem["rtrw_id"].(int64),
	// 	}

	// 	_, err2 := RtrwService.AddMbtile()
	// 	if err2 != nil {
	// 		err = err2
	// 		break
	// 	}
	// }
	// if err != nil {
	// 	tx.Rollback()
	// 	fmt.Println("Error:", err)
	// 	ctx.JSON(400, gin.H{
	// 		"status":      "error",
	// 		"status_code": 400,
	// 		"return":      err.Error(),
	// 	})
	// 	return
	// }
	rr := Model.RtrwType{}
	tx.Preload("Rtrw_groups").
		Preload("Rtrw_groups.Datas").Where("id = ?", rtrwData.Id).First(&rr)
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rr,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RtrwController) UpdateRtrw(ctx *gin.Context) {
	var props struct {
		Id              int64                    `json:"id" validate:"required"`
		Name            string                   `json:"name" validate:"required"`
		Place_string    string                   `json:"place_string" validate:"required"`
		Reg_Province_id int64                    `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64                    `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64                    `json:"reg_district_id"`
		Reg_Village_id  int64                    `json:"reg_village_id"`
		Status          string                   `json:"status" validate:"required"`
		Rtrw_groups     []map[string]interface{} `json:"rtrw_groups"`
		Rtrw_mbtiles    []interface{}            `json:"rtrw_mbtiles"`
	}

	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err := validate.Struct(props)
	if err != nil {
		var errArr []error
		var errValidators = err.(validator.ValidationErrors)
		for i := 0; i < len(errValidators); i++ {
			errValidItem := errValidators[i]
			errMess := fmt.Errorf("validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
			errArr = append(errArr, errMess)
		}
		err = errors.Join(errArr...)
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}

	tx := gorm_support.DB.Begin()
	var RtrwService = Service.RtrwService{
		DB: tx,
	}

	rtrwServiceUpdateData := RtrwService.RtrwServiceUpdateType
	rtrwServiceUpdateData.Id = props.Id
	rtrwServiceUpdateData.Name = props.Name
	rtrwServiceUpdateData.Place_string = props.Place_string
	rtrwServiceUpdateData.RegProvince_id = props.Reg_Province_id
	rtrwServiceUpdateData.RegRegency_id = props.Reg_Regency_id
	rtrwServiceUpdateData.RegDistrict_id = props.Reg_District_id
	rtrwServiceUpdateData.RegVillage_id = props.Reg_Village_id
	rtrwServiceUpdateData.Status = props.Status
	rtrwData, err := RtrwService.Update(rtrwServiceUpdateData)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error - E2MVAIDFMVIRTIEM", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}

	rtrwFileService := Service.RtrwFileService{
		DB: tx,
	}
	err = rtrwFileService.Gets(map[string]interface{}{}).Where("rtrw_id = ?", rtrwData.Id).Update("rtrw_id", nil).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Error - 29915223489:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}

	rtrwGroups := props.Rtrw_groups
	RtrwService.DeleteGroupByRtrwId(int(rtrwData.Id))
	for i := 0; i < len(rtrwGroups); i++ {
		rtrwGroupItem := rtrwGroups[i]
		rtrwGroupItem["rtrw_id"] = rtrwData.Id
		rtrwGroupData := RtrwService.RtrwGroupAddType
		if rtrwGroupItem["id"] != nil {
			rtrwGroupData.Id = int64(rtrwGroupItem["id"].(float64))
		}
		if rtrwGroupItem["uuid"] != nil {
			rtrwGroupData.Uuid = reflect.ValueOf(rtrwGroupItem["uuid"]).String()
		}
		rtrwGroupData.Rtrw_id = rtrwGroupItem["rtrw_id"].(int64)
		rtrwGroupData.Asset_key = rtrwGroupItem["asset_key"].(string)
		_rtrwGroupData_name := strings.ReplaceAll(rtrwGroupData.Asset_key, "_", " ")
		rtrwGroupData.Name = _rtrwGroupData_name
		rtrwGroupData.Properties = Helper.GetValue(rtrwGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
		err = validate.Struct(rtrwGroupData)
		if err != nil {
			var errArr []error
			var errValidators = err.(validator.ValidationErrors)
			for i := 0; i < len(errValidators); i++ {
				errValidItem := errValidators[i]
				errMess := fmt.Errorf("validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
				errArr = append(errArr, errMess)
			}
			err = errors.Join(errArr...)
			break
		}
		rtrwGroupuResult, err2 := RtrwService.AddGroup(rtrwGroupData)
		if err2 != nil {
			err = err2
			break
		}
		if rtrwGroupItem["datas"] != nil {
			rtrwFileDatas := rtrwGroupItem["datas"].([]interface{})
			for j := 0; j < len(rtrwFileDatas); j++ {
				rtrwFileItem := rtrwFileDatas[j].(map[string]interface{})
				rtrwFileProps := rtrwFileService.RtrwFileUpdate
				rtrwFileProps.Id = int64(Helper.GetValue(rtrwFileItem["id"], 0).(float64))
				rtrwFileProps.Rtrw_group_id = rtrwGroupuResult.Id
				rtrwFileProps.Rtrw_id = rtrwData.Id
				_, err3 := rtrwFileService.Update(rtrwFileProps)
				if err3 != nil {
					err = err3
				}
			}
		}
	}
	if err != nil {
		tx.Rollback()
		fmt.Println("Error - 2998923489:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}

	rtrwMbtiles := props.Rtrw_mbtiles
	var rtrw_mbtile_ids []int
	for i := 0; i < len(rtrwMbtiles); i++ {
		rtrwMbtileItem, _ := rtrwMbtiles[i].(map[string]interface{})
		intMbtileId := int(rtrwMbtileItem["id"].(float64))
		rtrw_mbtile_ids = append(rtrw_mbtile_ids, intMbtileId)
	}

	// Load RtrwMbtileService
	rtrwMbtileService := service.RtrwMbtileService{
		DB: tx,
	}
	// Delete rtrw mbtile by rtrw_id first
	err = rtrwMbtileService.Gets().Where("rtrw_id = ?", int(rtrwData.Id)).Update("rtrw_id", nil).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Error:", err)
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}

	for i := 0; i < len(rtrwMbtiles); i++ {
		rtrwMbtileItem, _ := rtrwMbtiles[i].(map[string]interface{})
		rtrwMbtileItem["rtrw_id"] = rtrwData.Id
		int64MbtileId := int64(rtrwMbtileItem["id"].(float64))
		RtrwService.RtrwMbtilePayload = struct {
			Id         *int64
			File_name  string
			Uuid       string
			Rtrw_id    int64
			Created_at string
			Updated_at string
			Checked_at string
		}{
			// For id Because pointer need vessel first
			Id:        &int64MbtileId,
			File_name: rtrwMbtileItem["file_name"].(string),
			Uuid:      rtrwMbtileItem["uuid"].(string),
			Rtrw_id:   int64(rtrwMbtileItem["rtrw_id"].(int64)),
		}

		_, err2 := RtrwService.AddMbtile()
		if err2 != nil {
			err = err2
			break
		}
	}
	if err != nil {
		tx.Rollback()
		fmt.Println("Error:", err)
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	tx.Commit()
	ctx.AddParam("id", strconv.Itoa(int(rtrwData.Id)))
	a.GetRtrwById(ctx)
	// rr := Model.RtrwType{}
	// tx.Preload("Rtrw_mbtiles").
	// 	Preload("Rtrw_groups").
	// 	Preload("Rtrw_groups.Datas").Where("id = ?", rtrwData.Id).First(&rr)

	// ctx.JSON(200, gin.H{
	// 	"return":      rr,
	// 	"status":      "success",
	// 	"status_code": 200,
	// })
}

func (a *RtrwController) DeleteRtrw(ctx *gin.Context) {
	var RtrwService = Service.RtrwService{}
	var props = map[string]interface{}{} // Bind the request body to the newUser struct
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var propsSt struct {
		Ids []int
	}
	Helper.ToStructFromMap(props, &propsSt)
	err := RtrwService.DeleteByIds(propsSt.Ids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Deleted Successfuly"})
}

func (a *RtrwController) ValidateMbtile(ctx *gin.Context) {
	var props struct {
		Rtrw_mbtile_ids         []int `json:"rtrw_mbtile_ids"`
		Rtrw_mbtile_ids_uncheck []int `json:"rtrw_mbtile_ids_uncheck"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rtrwMbtileService := Service.RtrwMbtileService{
		DB: gorm_support.DB,
	}

	martinConfig, err := os.ReadFile("./sub_app/martin/config.yaml")
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	var martinMap map[string]interface{}
	err = yaml.Unmarshal(martinConfig, &martinMap)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	martin_mbtile_sources := map[string]interface{}{}
	martin_mbtiles := martinMap["mbtiles"].(map[string]interface{})
	martin_mbtile_sources_parse, ok := martin_mbtiles["sources"].(map[string]interface{})
	if ok {
		martin_mbtile_sources = martin_mbtile_sources_parse
	}

	var mbtile_datas = []model.RtrwMbtile{}

	// First get uncheck data first
	rtrwMbtileDB := rtrwMbtileService.Gets()
	err = rtrwMbtileDB.Where("id IN ?", props.Rtrw_mbtile_ids_uncheck).Find(&mbtile_datas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Iterate over the map
	for key, value := range martin_mbtile_sources {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
		for _, value2 := range mbtile_datas {
			if key == value2.UUID {
				delete(martin_mbtile_sources, value2.UUID)
				break
			}
		}
	}

	// Then set checked_at to be null
	rtrwMbtileDB = rtrwMbtileService.Gets()
	rtrwMbtileDB.Where("id IN ?", props.Rtrw_mbtile_ids_uncheck).Update("checked_at", nil)

	// Next get all check datas
	rtrwMbtileDB = rtrwMbtileService.Gets()
	err = rtrwMbtileDB.Where("id IN ?", props.Rtrw_mbtile_ids).Find(&mbtile_datas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Iterate over the map
	for _, value2 := range mbtile_datas {
		martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/rtrw/", value2.UUID, ".mbtiles")
	}

	// Redefine again to martin config mbtiles
	martin_mbtiles["sources"] = martin_mbtile_sources
	martinMap["mbtiles"] = martin_mbtiles

	// Then set checked_at
	rtrwMbtileDB = rtrwMbtileService.Gets()
	currentTime := time.Now()
	rtrwMbtileDB.Where("id IN ?", props.Rtrw_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

	// The last save config martin
	martinConfig, err = yaml.Marshal(martinMap)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = os.WriteFile("./sub_app/martin/config.yaml", []byte(martinConfig), 0755)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Then set checked_at
	rtrwMbtileDB = rtrwMbtileService.Gets()
	rtrwMbtileDB.Where("id IN ?", props.Rtrw_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

	ctx.JSON(200, gin.H{
		"return":      martinConfig,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RtrwController) ValidateKml(ctx *gin.Context) {
	var props struct {
		Rtrw_group_ids []int `json:"rtrw_group_ids"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rtrwService := Service.RtrwService{
		DB: gorm_support.DB,
	}
	asyncJObService := Service.AsynqJobService{
		DB: gorm_support.DB,
	}

	// Reset rtrw file remove validated_at
	rtrwFileService := Service.RtrwFileService{
		DB: gorm_support.DB,
	}
	rtrwFileService.Unvalidated(props.Rtrw_group_ids)

	var queue_ids []map[string]interface{}

	for i := 0; i < len(props.Rtrw_group_ids); i++ {
		fmt.Println("props.Rtrw_group_ids[i]", props.Rtrw_group_ids[i])
		taskk, err := asynq_support.NewValidateKmlTask(props.Rtrw_group_ids[i], "rtrw")
		if err != nil {
			log.Fatalf("Could not schedule task : %v", err)
		}
		// RUNNING QUEUE WITH DATA
		gg, err := asynq_support.Client.Enqueue(taskk, asynq.Queue("default"), asynq.ProcessIn(3*time.Second))
		if err != nil {
			log.Fatalf("could not scheudle task: %v", err)
		}

		fmt.Println("Task info :: ", gg.ID, " :: ", gg.Queue, " :: ", gg.Type)

		queue_ids = append(queue_ids, map[string]interface{}{
			"id":       props.Rtrw_group_ids[i],
			"asynq_id": gg.ID,
		})

		rtrwGroupItemData := Model.RtrwGroup{}
		err = rtrwService.GetRtrwGroups().Where("id = ?", props.Rtrw_group_ids[i]).First(&rtrwGroupItemData).Error
		if err != nil {
			log.Fatalf("could not scheudle task: %v", err)
			break
		}
		resAsynqData, err := asyncJObService.Add(Service.AsynqJobAddPayload{
			App_uuid:     rtrwGroupItemData.Uuid,
			Asynq_uuid:   gg.ID,
			Payload:      string(gg.Payload),
			Status:       asyncJObService.GetStatus().STATUS_PENDING,
			Table_name:   rtrwGroupItemData.TableName(),
			Message_text: "Job is created",
		})
		if err != nil {
			log.Fatalf("could not scheudle task: %v", err)
			break
		}
		fmt.Println("resAsynqData :: ", resAsynqData)
		// log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
	ctx.JSON(200, gin.H{
		"return":      queue_ids,
		"status":      "success",
		"status_code": 200,
	})
}

// Handle Websocket

func (a *RtrwController) HandleWS(ctx *gin.Context) {
	var asynq_ids = map[string]*aty2type{}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	bus := EventBus.New()

	eventBusKey := helper.RandStringBytes(10)
	bus.Subscribe(eventBusKey, func(message []byte) {
		conn.WriteMessage(websocket.TextMessage, message)
	})

	var checkStatusAsynqClose func() bool
	var checkStatusRtrwGroupClose func() bool

	defer func() {
		fmt.Println("Websocket from client is closed")
		if checkStatusRtrwGroupClose != nil {
			checkStatusRtrwGroupClose()
		}
		if checkStatusAsynqClose != nil {
			checkStatusAsynqClose()
		}
		for _, v := range asynq_ids {
			fmt.Println("File are closed")
			v.Is_run = false
		}
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		dataParse := map[string]interface{}{}
		err = json.Unmarshal(message, &dataParse)
		if err != nil {
			log.Println(err)
			return
		}

		var action string = dataParse["action"].(string)
		switch action {
		case "CHECK_ASYNQ_STATUS":
			if checkStatusAsynqClose == nil {
				jj, ok := dataParse["uuids"]
				if !ok {
					log.Println("Problem check interface")
					return
				}
				// Create a new slice of string
				stringSlice := convert.SliceOfString(jj)
				checkStatusAsynqClose = checkAsynqStatusClosure(bus, eventBusKey, stringSlice)
			}
		case "CHECK_VALIDATED":
			if checkStatusRtrwGroupClose == nil {
				jj, ok := dataParse["group_ids"]
				if !ok {
					log.Println("Problem check interface")
					return
				}

				// Create a new slice of int64
				int64Slice := helper.SliceOfInt64(jj)
				checkStatusRtrwGroupClose = checkStatusRtrwGroupClousure(bus, eventBusKey, int64Slice)
			}

		case "TAIL":
			var asynq_id string = dataParse["asynq_id"].(string)
			if (asynq_ids)[asynq_id] == nil {
				(asynq_ids)[asynq_id] = &aty2type{
					Is_run:      true,
					Line_number: 0,
				}

			} else {
				asynq_ids[asynq_id].Is_run = !asynq_ids[asynq_id].Is_run
			}
			if asynq_ids[asynq_id].Is_run {
				go func(bus EventBus.Bus, eventBusKey string, ass *aty2type) {
					for ass.Is_run {
						time.Sleep(2 * time.Second)
						tailLog(bus, eventBusKey, asynq_id, ass)
					}
				}(bus, eventBusKey, asynq_ids[asynq_id])
			}
		}
	}
}

func checkStatusRtrwGroupClousure(bus EventBus.Bus, eventBusKey string, ids []int64) func() bool {
	is_loop := true
	go func(ids []int64, s *bool) {
		for *s {
			var group_ids []int64 = ids
			rtrw_group_datas, err := checkStatusRtrwGroups(group_ids)
			if err != nil {
				bus.Publish(eventBusKey, []byte(err.Error()))
				return
			}
			time.Sleep(5 * time.Second)
			fmt.Println("checkStatusRtrwGroupClousure - running")
			textT := map[string]interface{}{}
			textT["from"] = "check_group"
			textT["message"] = rtrw_group_datas
			textTSTrng, _ := json.Marshal(textT)
			bus.Publish(eventBusKey, textTSTrng)
		}
		fmt.Println("checkStatusRtrwGroupClousure - stop")
	}(ids, &is_loop)
	return func() bool {
		is_loop = false
		return is_loop
	}
}

func checkStatusRtrwGroups(ids []int64) ([]Model.RtrwGroupView, error) {
	rtrw_service := Service.RtrwService{
		DB: gorm_support.DB,
	}
	rtrw_groupModel := rtrw_service.GetRtrwGroups()
	rtrw_group_datas := []Model.RtrwGroupView{}
	err := rtrw_groupModel.Preload("Datas").
		Select("rtrw_group.*, " +
			"COALESCE((SELECT COUNT(*) FROM rtrw_file  WHERE rtrw_file.rtrw_group_id = rtrw_group.id AND rtrw_file.validated_at IS NULL),0) AS unvalidated, " +
			"COALESCE((SELECT COUNT(*) FROM rtrw_file WHERE rtrw_file.rtrw_group_id = rtrw_group.id AND rtrw_file.validated_at IS NOT NULL),0) as validated").
		Where([]int64(ids)).Find(&rtrw_group_datas).Error
	if err != nil {
		return []Model.RtrwGroupView{}, err
	}

	return rtrw_group_datas, err
}
