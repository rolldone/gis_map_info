package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"gis_map_info/app/helper"
	Helper "gis_map_info/app/helper"
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"
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

type ZlpController struct{}

func (a *ZlpController) GetZlps(ctx *gin.Context) {
	var ZlpService = Service.ZlpService{}
	zlpDB := ZlpService.Gets()
	var zlpDatas = []Model.ZlpType{}
	// Fetch data from the database
	if err := zlpDB.Find(&zlpDatas).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      zlpDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *ZlpController) GetZlpsPaginate(ctx *gin.Context) {
	var ZlpService = Service.ZlpService{}
	zlpDB := ZlpService.Gets()
	var zlpDatas = []Model.ZlpType{}
	var page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	var limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	var offset = (page - 1) * limit

	// Fetch data from the database
	if err := zlpDB.
		Preload("Reg_province").
		Preload("Reg_regency").
		Preload("Reg_district").
		Preload("Zlp_groups").
		Preload("Zlp_groups.Datas").
		Limit(limit).Offset(offset).Order("updated_at DESC").Find(&zlpDatas).Error; err != nil {
		fmt.Println("Error - 29mvamfivm2 ", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      zlpDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *ZlpController) GetZlpById(ctx *gin.Context) {
	var ZlpService = Service.ZlpService{}
	id, err := strconv.Atoi(ctx.Param("id"))
	// Check for any conversion error
	if err != nil {
		fmt.Println("GetZlpById Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Something wrong with parameter",
		})
		return
	}
	zlpDB := ZlpService.GetById(id)
	var zlpData = Model.ZlpType{}
	// Fetch data from the database
	if err := zlpDB.
		Preload("Reg_province").
		Preload("Reg_district").
		Preload("Zlp_groups", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Datas", func(db *gorm.DB) *gorm.DB {
				return db.Where("zlp_id IS NOT NULL")
			}).Select("zlp_group.*, COALESCE((SELECT COUNT(*) FROM zlp_file  WHERE zlp_file.zlp_group_id = zlp_group.id AND zlp_file.validated_at IS NULL AND zlp_file.zlp_id IS NOT NULL),0) AS unvalidated, " +
				"COALESCE((SELECT COUNT(*) FROM zlp_file WHERE zlp_file.zlp_group_id = zlp_group.id AND zlp_file.validated_at IS NOT NULL AND zlp_file.zlp_id IS NOT NULL),0) as validated")
		}).
		Preload("Zlp_mbtiles").
		Preload("Zlp_groups.Datas").First(&zlpData).Error; err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(404, gin.H{
			"status":      "error",
			"status_code": 404,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      zlpData,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *ZlpController) AddZlp(ctx *gin.Context) {
	var props struct {
		Name            string        `json:"name" validate:"required"`
		Place_string    string        `json:"place_string" validate:"required"`
		Reg_Province_id int64         `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64         `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64         `json:"reg_district_id"`
		Reg_Village_id  int64         `json:"reg_village_id"`
		Status          string        `json:"status"`
		Zlp_groups      []interface{} `json:"zlp_groups"`
		Zlp_mbtiles     []interface{} `json:"zlp_mbtiles"`
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
	var ZlpService = Service.ZlpService{
		DB: tx,
	}

	zlpServiceAddData := ZlpService.ZlpServiceAddType
	zlpServiceAddData.Name = props.Name
	zlpServiceAddData.Place_string = props.Place_string
	zlpServiceAddData.RegProvince_id = props.Reg_Province_id
	zlpServiceAddData.RegRegency_id = props.Reg_Regency_id
	zlpServiceAddData.RegDistrict_id = props.Reg_District_id
	zlpServiceAddData.RegVillage_id = props.Reg_Village_id
	zlpServiceAddData.Status = props.Status
	zlpData, err := ZlpService.Add(zlpServiceAddData)
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
	// zlpGroups := props.Zlp_groups
	// fmt.Println(zlpGroups)
	// for i := 0; i < len(zlpGroups); i++ {
	// 	zlpGroupItem, _ := zlpGroups[i].(map[string]interface{})
	// 	zlpGroupItem["zlp_id"] = zlpData.Id
	// 	zlpGroupData := ZlpService.ZlpGroupAddType
	// 	zlpGroupData.Name = zlpGroupItem["name"].(string)
	// 	zlpGroupData.Zlp_id = int64(zlpGroupItem["zlp_id"].(float64))
	// 	zlpGroupData.Asset_key = zlpGroupItem["asset_key"].(string)
	// 	zlpGroupData.Uuid = uuid.NewString()
	// 	_zlpGroupItem_properties, _ := Helper.GetValue(zlpGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
	// 	zlpGroupData.Properties = _zlpGroupItem_properties
	// 	err = validate.Struct(zlpGroupData)
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
	// 	_, err2 := ZlpService.AddGroup(zlpGroupData)
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

	// zlpMbtiles := props.Zlp_mbtiles
	// // Delete zlp mbtile by zlp_id first
	// ZlpService.DeleteMbtileByZlpId(int(zlpData.Id))
	// for i := 0; i < len(zlpMbtiles); i++ {
	// 	zlpMbtileItem, _ := zlpMbtiles[i].(map[string]interface{})
	// 	zlpMbtileItem["zlp_id"] = zlpData.Id
	// 	int64MbtileId := int64(zlpMbtileItem["id"].(float64))
	// 	ZlpService.ZlpMbtilePayload = struct {
	// 		Id        *int64
	// 		File_name string
	// 		Uuid      string
	// 		Zlp_id   int64
	// 	}{
	// 		// For id Because pointer need vessel first
	// 		Id:        &int64MbtileId,
	// 		File_name: zlpMbtileItem["file_name"].(string),
	// 		Uuid:      zlpMbtileItem["uuid"].(string),
	// 		Zlp_id:   zlpMbtileItem["zlp_id"].(int64),
	// 	}

	// 	_, err2 := ZlpService.AddMbtile()
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
	rr := Model.ZlpType{}
	tx.Preload("Zlp_groups").
		Preload("Zlp_groups.Datas").Where("id = ?", zlpData.Id).First(&rr)
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rr,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *ZlpController) UpdateZlp(ctx *gin.Context) {
	var props struct {
		Id              int64                    `json:"id" validate:"required"`
		Name            string                   `json:"name" validate:"required"`
		Place_string    string                   `json:"place_string" validate:"required"`
		Reg_Province_id int64                    `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64                    `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64                    `json:"reg_district_id"`
		Reg_Village_id  int64                    `json:"reg_village_id"`
		Status          string                   `json:"status" validate:"required"`
		Zlp_groups      []map[string]interface{} `json:"zlp_groups"`
		Zlp_mbtiles     []interface{}            `json:"zlp_mbtiles"`
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
	var ZlpService = Service.ZlpService{
		DB: tx,
	}

	zlpServiceUpdateData := ZlpService.ZlpServiceUpdateType
	zlpServiceUpdateData.Id = props.Id
	zlpServiceUpdateData.Name = props.Name
	zlpServiceUpdateData.Place_string = props.Place_string
	zlpServiceUpdateData.RegProvince_id = props.Reg_Province_id
	zlpServiceUpdateData.RegRegency_id = props.Reg_Regency_id
	zlpServiceUpdateData.RegDistrict_id = props.Reg_District_id
	zlpServiceUpdateData.RegVillage_id = props.Reg_Village_id
	zlpServiceUpdateData.Status = props.Status
	zlpData, err := ZlpService.Update(zlpServiceUpdateData)
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
	zlpFileService := Service.ZlpFileService{
		DB: tx,
	}
	zlpGroups := props.Zlp_groups
	ZlpService.DeleteGroupByZlpId(int(zlpData.Id))
	for i := 0; i < len(zlpGroups); i++ {
		zlpGroupItem := zlpGroups[i]
		zlpGroupItem["zlp_id"] = zlpData.Id
		zlpGroupData := ZlpService.ZlpGroupAddType
		if zlpGroupItem["id"] != nil {
			zlpGroupData.Id = int64(zlpGroupItem["id"].(float64))
		}
		if zlpGroupItem["uuid"] != nil {
			zlpGroupData.Uuid = reflect.ValueOf(zlpGroupItem["uuid"]).String()
		}
		zlpGroupData.Zlp_id = zlpGroupItem["zlp_id"].(int64)
		zlpGroupData.Asset_key = zlpGroupItem["asset_key"].(string)
		_zlpGroupData_name := strings.ReplaceAll(zlpGroupData.Asset_key, "_", " ")
		zlpGroupData.Name = _zlpGroupData_name
		zlpGroupData.Properties = Helper.GetValue(zlpGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
		err = validate.Struct(zlpGroupData)
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
		zlpGroupuResult, err2 := ZlpService.AddGroup(zlpGroupData)
		if err2 != nil {
			err = err2
			break
		}
		if zlpGroupItem["datas"] != nil {
			zlpFileDatas := zlpGroupItem["datas"].([]interface{})
			for j := 0; j < len(zlpFileDatas); j++ {
				zlpFileItem := zlpFileDatas[j].(map[string]interface{})
				zlpFileProps := zlpFileService.ZlpFileUpdate
				zlpFileProps.Id = int64(Helper.GetValue(zlpFileItem["id"], 0).(float64))
				zlpFileProps.Zlp_group_id = zlpGroupuResult.Id
				zlpFileProps.Zlp_id = zlpData.Id
				_, err3 := zlpFileService.Update(zlpFileProps)
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

	zlpMbtiles := props.Zlp_mbtiles
	var zlp_mbtile_ids []int
	for i := 0; i < len(zlpMbtiles); i++ {
		zlpMbtileItem, _ := zlpMbtiles[i].(map[string]interface{})
		intMbtileId := int(zlpMbtileItem["id"].(float64))
		zlp_mbtile_ids = append(zlp_mbtile_ids, intMbtileId)
	}
	// Delete zlp mbtile by zlp_id first
	err = ZlpService.DeleteMbtileExceptZlpMbtileIds_withZlp_id(zlp_mbtile_ids, int(zlpData.Id))
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
	for i := 0; i < len(zlpMbtiles); i++ {
		zlpMbtileItem, _ := zlpMbtiles[i].(map[string]interface{})
		zlpMbtileItem["zlp_id"] = zlpData.Id
		int64MbtileId := int64(zlpMbtileItem["id"].(float64))
		zlp_group_id_float64, ok := zlpMbtileItem["zlp_group_id"].(float64)
		if !ok {
			zlp_group_id_float64 = 0
		}
		zlp_group_id := int64(zlp_group_id_float64)
		ZlpService.ZlpMbtilePayload = Service.ZlpMbtilePayload{
			Id:              &int64MbtileId,
			File_name:       zlpMbtileItem["file_name"].(string),
			Uuid:            zlpMbtileItem["uuid"].(string),
			Zlp_id:          int64(zlpMbtileItem["zlp_id"].(int64)),
			Asset_key:       zlpMbtileItem["asset_key"].(string),
			Reg_province_id: props.Reg_Province_id,
			Zlp_group_id:    zlp_group_id,
		}
		_, err2 := ZlpService.AddMbtile()
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
	ctx.AddParam("id", strconv.Itoa(int(zlpData.Id)))
	a.GetZlpById(ctx)
	// rr := Model.ZlpType{}
	// tx.Preload("Zlp_mbtiles").
	// 	Preload("Zlp_groups").
	// 	Preload("Zlp_groups.Datas").Where("id = ?", zlpData.Id).First(&rr)

	// ctx.JSON(200, gin.H{
	// 	"return":      rr,
	// 	"status":      "success",
	// 	"status_code": 200,
	// })
}

func (a *ZlpController) DeleteZlp(ctx *gin.Context) {
	var ZlpService = Service.ZlpService{}
	var props = map[string]interface{}{} // Bind the request body to the newUser struct
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var propsSt struct {
		Ids []int
	}
	Helper.ToStructFromMap(props, &propsSt)
	err := ZlpService.DeleteByIds(propsSt.Ids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Deleted Successfuly"})
}

func (a *ZlpController) ValidateMbtile(ctx *gin.Context) {
	var props struct {
		Zlp_mbtile_ids         []int `json:"zlp_mbtile_ids"`
		Zlp_mbtile_ids_uncheck []int `json:"zlp_mbtile_ids_uncheck"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zlpMbtileService := Service.ZlpMbtileService{
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

	martin_mbtiles := martinMap["mbtiles"].(map[string]interface{})
	martin_mbtile_sources := martin_mbtiles["sources"].(map[string]interface{})

	var mbtile_datas = []model.ZlpMbtile{}

	// First get uncheck data first
	zlpMbtileDB := zlpMbtileService.Gets()
	err = zlpMbtileDB.Where("id IN ?", props.Zlp_mbtile_ids_uncheck).Find(&mbtile_datas).Error
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
	zlpMbtileDB = zlpMbtileService.Gets()
	zlpMbtileDB.Where("id IN ?", props.Zlp_mbtile_ids_uncheck).Update("checked_at", nil)

	// Next get all check datas
	zlpMbtileDB = zlpMbtileService.Gets()
	err = zlpMbtileDB.Where("id IN ?", props.Zlp_mbtile_ids).Find(&mbtile_datas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Iterate over the map
	for _, value2 := range mbtile_datas {
		martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/zlp/", value2.UUID, ".mbtiles")
	}

	// Redefine again to martin config mbtiles
	martin_mbtiles["sources"] = martin_mbtile_sources
	martinMap["mbtiles"] = martin_mbtiles

	// Then set checked_at
	zlpMbtileDB = zlpMbtileService.Gets()
	currentTime := time.Now()
	zlpMbtileDB.Where("id IN ?", props.Zlp_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

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
	zlpMbtileDB = zlpMbtileService.Gets()
	zlpMbtileDB.Where("id IN ?", props.Zlp_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

	ctx.JSON(200, gin.H{
		"return":      martinConfig,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *ZlpController) ValidateKml(ctx *gin.Context) {
	var props struct {
		Zlp_group_ids []int `json:"zlp_group_ids"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zlpService := Service.ZlpService{
		DB: gorm_support.DB,
	}
	asyncJObService := Service.AsynqJobService{
		DB: gorm_support.DB,
	}

	// Reset zlp file remove validated_at
	zlpFileService := Service.ZlpFileService{
		DB: gorm_support.DB,
	}
	zlpFileService.Unvalidated(props.Zlp_group_ids)

	var queue_ids []map[string]interface{}

	for i := 0; i < len(props.Zlp_group_ids); i++ {
		fmt.Println("props.Zlp_group_ids[i]", props.Zlp_group_ids[i])
		taskk, err := asynq_support.NewValidateKmlTask(props.Zlp_group_ids[i], "zlp")
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
			"id":       props.Zlp_group_ids[i],
			"asynq_id": gg.ID,
		})

		zlpGroupItemData := Model.ZlpGroup{}
		err = zlpService.GetZlpGroups().Where("id = ?", props.Zlp_group_ids[i]).First(&zlpGroupItemData).Error
		if err != nil {
			log.Fatalf("could not scheudle task: %v", err)
			break
		}
		resAsynqData, err := asyncJObService.Add(Service.AsynqJobAddPayload{
			App_uuid:     zlpGroupItemData.Uuid,
			Asynq_uuid:   gg.ID,
			Payload:      string(gg.Payload),
			Status:       asyncJObService.GetStatus().STATUS_PENDING,
			Table_name:   zlpGroupItemData.TableName(),
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

func (a *ZlpController) HandleWS(ctx *gin.Context) {
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
	var checkStatusZlpGroupClose func() bool

	defer func() {
		fmt.Println("Websocket from client is closed")
		if checkStatusZlpGroupClose != nil {
			checkStatusZlpGroupClose()
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
			if checkStatusZlpGroupClose == nil {
				jj, ok := dataParse["group_ids"]
				if !ok {
					log.Println("Problem check interface")
					return
				}

				// Create a new slice of int64
				int64Slice := helper.SliceOfInt64(jj)
				checkStatusZlpGroupClose = checkStatusZlpGroupClousure(bus, eventBusKey, int64Slice)
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

func checkStatusZlpGroupClousure(bus EventBus.Bus, eventBusKey string, ids []int64) func() bool {
	is_loop := true
	go func(ids []int64, s *bool) {
		for *s {
			var group_ids []int64 = ids
			zlp_group_datas, err := checkStatusZlpGroups(group_ids)
			if err != nil {
				bus.Publish(eventBusKey, []byte(err.Error()))
				return
			}
			time.Sleep(5 * time.Second)
			fmt.Println("checkStatusZlpGroupClousure - running")
			textT := map[string]interface{}{}
			textT["from"] = "check_group"
			textT["message"] = zlp_group_datas
			textTSTrng, _ := json.Marshal(textT)
			bus.Publish(eventBusKey, textTSTrng)
		}
		fmt.Println("checkStatusZlpGroupClousure - stop")
	}(ids, &is_loop)
	return func() bool {
		is_loop = false
		return is_loop
	}
}

func checkStatusZlpGroups(ids []int64) ([]Model.ZlpGroupView, error) {
	zlp_service := Service.ZlpService{
		DB: gorm_support.DB,
	}
	zlp_groupModel := zlp_service.GetZlpGroups()
	zlp_group_datas := []Model.ZlpGroupView{}
	err := zlp_groupModel.Preload("Datas").
		Select("zlp_group.*, " +
			"COALESCE((SELECT COUNT(*) FROM zlp_file  WHERE zlp_file.zlp_group_id = zlp_group.id AND zlp_file.validated_at IS NULL AND zlp_file.zlp_id != 0),0) AS unvalidated, " +
			"COALESCE((SELECT COUNT(*) FROM zlp_file WHERE zlp_file.zlp_group_id = zlp_group.id AND zlp_file.validated_at IS NOT NULL AND zlp_file.zlp_id != 0),0) as validated").
		Where([]int64(ids)).Find(&zlp_group_datas).Error
	if err != nil {
		return []Model.ZlpGroupView{}, err
	}

	return zlp_group_datas, err
}
