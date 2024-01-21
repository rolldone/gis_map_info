package admin

import (
	"bufio"
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

	"github.com/benpate/convert"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type RdtrController struct{}

func (a *RdtrController) GetRdtrs(ctx *gin.Context) {
	var RdtrService = Service.RdtrService{}
	rdtrDB := RdtrService.Gets()
	var rdtrDatas = []Model.RdtrType{}
	// Fetch data from the database
	if err := rdtrDB.Find(&rdtrDatas).Error; err != nil {
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rdtrDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) GetRdtrsPaginate(ctx *gin.Context) {
	var RdtrService = Service.RdtrService{}
	rdtrDB := RdtrService.Gets()
	var rdtrDatas = []Model.RdtrType{}
	var page, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	var limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	var offset = (page - 1) * limit

	// Fetch data from the database
	if err := rdtrDB.
		Preload("Reg_province").
		Preload("Reg_regency").
		Preload("Reg_district").
		Preload("Rdtr_groups").
		Preload("Rdtr_groups.Datas").
		Limit(limit).Offset(offset).Order("updated_at DESC").Find(&rdtrDatas).Error; err != nil {
		fmt.Println("Error - 29mvamfivm2 ", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rdtrDatas,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) GetRdtrById(ctx *gin.Context) {
	var RdtrService = Service.RdtrService{}
	id, err := strconv.Atoi(ctx.Param("id"))
	// Check for any conversion error
	if err != nil {
		fmt.Println("GetRdtrById Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      "Something wrong with parameter",
		})
		return
	}
	rdtrDB := RdtrService.GetById(id)
	var rdtrData = Model.RdtrType{}
	// Fetch data from the database
	if err := rdtrDB.
		Preload("Reg_province").
		Preload("Reg_district").
		Preload("Rdtr_groups", func(db *gorm.DB) *gorm.DB {
			return db.Select("rdtr_group.*, COALESCE((SELECT COUNT(*) FROM rdtr_file  WHERE rdtr_file.rdtr_group_id = rdtr_group.id AND rdtr_file.validated_at IS NULL),0) AS unvalidated, " +
				"COALESCE((SELECT COUNT(*) FROM rdtr_file WHERE rdtr_file.rdtr_group_id = rdtr_group.id AND rdtr_file.validated_at IS NOT NULL),0) as validated")
		}).
		Preload("Rdtr_mbtiles").
		Preload("Rdtr_groups.Datas").First(&rdtrData).Error; err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(404, gin.H{
			"status":      "error",
			"status_code": 404,
			"return":      "Failed to fetch data",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"return":      rdtrData,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) AddRdtr(ctx *gin.Context) {
	var props struct {
		Name            string        `json:"name" validate:"required"`
		Place_string    string        `json:"place_string" validate:"required"`
		Reg_Province_id int64         `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64         `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64         `json:"reg_district_id"`
		Reg_Village_id  int64         `json:"reg_village_id"`
		Status          string        `json:"status"`
		Rdtr_groups     []interface{} `json:"rdtr_groups"`
		Rdtr_mbtiles    []interface{} `json:"rdtr_mbtiles"`
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
	var RdtrService = Service.RdtrService{
		DB: tx,
	}

	rdtrServiceAddData := RdtrService.RdtrServiceAddType
	rdtrServiceAddData.Name = props.Name
	rdtrServiceAddData.Place_string = props.Place_string
	rdtrServiceAddData.RegProvince_id = props.Reg_Province_id
	rdtrServiceAddData.RegRegency_id = props.Reg_Regency_id
	rdtrServiceAddData.RegDistrict_id = props.Reg_District_id
	rdtrServiceAddData.RegVillage_id = props.Reg_Village_id
	rdtrServiceAddData.Status = props.Status
	rdtrData, err := RdtrService.Add(rdtrServiceAddData)
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
	// rdtrGroups := props.Rdtr_groups
	// fmt.Println(rdtrGroups)
	// for i := 0; i < len(rdtrGroups); i++ {
	// 	rdtrGroupItem, _ := rdtrGroups[i].(map[string]interface{})
	// 	rdtrGroupItem["rdtr_id"] = rdtrData.Id
	// 	rdtrGroupData := RdtrService.RdtrGroupAddType
	// 	rdtrGroupData.Name = rdtrGroupItem["name"].(string)
	// 	rdtrGroupData.Rdtr_id = int64(rdtrGroupItem["rdtr_id"].(float64))
	// 	rdtrGroupData.Asset_key = rdtrGroupItem["asset_key"].(string)
	// 	rdtrGroupData.Uuid = uuid.NewString()
	// 	_rdtrGroupItem_properties, _ := Helper.GetValue(rdtrGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
	// 	rdtrGroupData.Properties = _rdtrGroupItem_properties
	// 	err = validate.Struct(rdtrGroupData)
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
	// 	_, err2 := RdtrService.AddGroup(rdtrGroupData)
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

	// rdtrMbtiles := props.Rdtr_mbtiles
	// // Delete rdtr mbtile by rdtr_id first
	// RdtrService.DeleteMbtileByRdtrId(int(rdtrData.Id))
	// for i := 0; i < len(rdtrMbtiles); i++ {
	// 	rdtrMbtileItem, _ := rdtrMbtiles[i].(map[string]interface{})
	// 	rdtrMbtileItem["rdtr_id"] = rdtrData.Id
	// 	int64MbtileId := int64(rdtrMbtileItem["id"].(float64))
	// 	RdtrService.RdtrMbtilePayload = struct {
	// 		Id        *int64
	// 		File_name string
	// 		Uuid      string
	// 		Rdtr_id   int64
	// 	}{
	// 		// For id Because pointer need vessel first
	// 		Id:        &int64MbtileId,
	// 		File_name: rdtrMbtileItem["file_name"].(string),
	// 		Uuid:      rdtrMbtileItem["uuid"].(string),
	// 		Rdtr_id:   rdtrMbtileItem["rdtr_id"].(int64),
	// 	}

	// 	_, err2 := RdtrService.AddMbtile()
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
	rr := Model.RdtrType{}
	tx.Preload("Rdtr_groups").
		Preload("Rdtr_groups.Datas").Where("id = ?", rdtrData.Id).First(&rr)
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rr,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) UpdateRdtr(ctx *gin.Context) {
	var props struct {
		Id              int64                    `json:"id" validate:"required"`
		Name            string                   `json:"name" validate:"required"`
		Place_string    string                   `json:"place_string" validate:"required"`
		Reg_Province_id int64                    `json:"reg_province_id" validate:"required"`
		Reg_Regency_id  int64                    `json:"reg_regency_id" validate:"required"`
		Reg_District_id int64                    `json:"reg_district_id"`
		Reg_Village_id  int64                    `json:"reg_village_id"`
		Status          string                   `json:"status" validate:"required"`
		Rdtr_groups     []map[string]interface{} `json:"rdtr_groups"`
		Rdtr_mbtiles    []interface{}            `json:"rdtr_mbtiles"`
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
	var RdtrService = Service.RdtrService{
		DB: tx,
	}

	rdtrServiceUpdateData := RdtrService.RdtrServiceUpdateType
	rdtrServiceUpdateData.Id = props.Id
	rdtrServiceUpdateData.Name = props.Name
	rdtrServiceUpdateData.Place_string = props.Place_string
	rdtrServiceUpdateData.RegProvince_id = props.Reg_Province_id
	rdtrServiceUpdateData.RegRegency_id = props.Reg_Regency_id
	rdtrServiceUpdateData.RegDistrict_id = props.Reg_District_id
	rdtrServiceUpdateData.RegVillage_id = props.Reg_Village_id
	rdtrServiceUpdateData.Status = props.Status
	rdtrData, err := RdtrService.Update(rdtrServiceUpdateData)
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
	rdtrFileService := Service.RdtrFileService{
		DB: tx,
	}
	rdtrGroups := props.Rdtr_groups
	RdtrService.DeleteGroupByRdtrId(int(rdtrData.Id))
	for i := 0; i < len(rdtrGroups); i++ {
		rdtrGroupItem := rdtrGroups[i]
		rdtrGroupItem["rdtr_id"] = rdtrData.Id
		rdtrGroupData := RdtrService.RdtrGroupAddType
		if rdtrGroupItem["id"] != nil {
			rdtrGroupData.Id = int64(rdtrGroupItem["id"].(float64))
		}
		if rdtrGroupItem["uuid"] != nil {
			rdtrGroupData.Uuid = reflect.ValueOf(rdtrGroupItem["uuid"]).String()
		}
		rdtrGroupData.Rdtr_id = rdtrGroupItem["rdtr_id"].(int64)
		rdtrGroupData.Asset_key = rdtrGroupItem["asset_key"].(string)
		_rdtrGroupData_name := strings.ReplaceAll(rdtrGroupData.Asset_key, "_", " ")
		rdtrGroupData.Name = _rdtrGroupData_name
		rdtrGroupData.Properties = Helper.GetValue(rdtrGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
		err = validate.Struct(rdtrGroupData)
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
		rdtrGroupuResult, err2 := RdtrService.AddGroup(rdtrGroupData)
		if err2 != nil {
			err = err2
			break
		}
		if rdtrGroupItem["datas"] != nil {
			rdtrFileDatas := rdtrGroupItem["datas"].([]interface{})
			for j := 0; j < len(rdtrFileDatas); j++ {
				rdtrFileItem := rdtrFileDatas[j].(map[string]interface{})
				rdtrFileProps := rdtrFileService.RdtrFileUpdate
				rdtrFileProps.Id = int64(Helper.GetValue(rdtrFileItem["id"], 0).(float64))
				rdtrFileProps.Rdtr_group_id = rdtrGroupuResult.Id
				rdtrFileProps.Rdtr_id = rdtrData.Id
				_, err3 := rdtrFileService.Update(rdtrFileProps)
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

	rdtrMbtiles := props.Rdtr_mbtiles
	var rdtr_mbtile_ids []int
	for i := 0; i < len(rdtrMbtiles); i++ {
		rdtrMbtileItem, _ := rdtrMbtiles[i].(map[string]interface{})
		intMbtileId := int(rdtrMbtileItem["id"].(float64))
		rdtr_mbtile_ids = append(rdtr_mbtile_ids, intMbtileId)
	}
	// Delete rdtr mbtile by rdtr_id first
	err = RdtrService.DeleteMbtileExceptRdtrMbtileIds_withRdtr_id(rdtr_mbtile_ids, int(rdtrData.Id))
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
	for i := 0; i < len(rdtrMbtiles); i++ {
		rdtrMbtileItem, _ := rdtrMbtiles[i].(map[string]interface{})
		rdtrMbtileItem["rdtr_id"] = rdtrData.Id
		int64MbtileId := int64(rdtrMbtileItem["id"].(float64))
		RdtrService.RdtrMbtilePayload = struct {
			Id         *int64
			File_name  string
			Uuid       string
			Rdtr_id    int64
			Created_at string
			Updated_at string
			Checked_at string
		}{
			// For id Because pointer need vessel first
			Id:        &int64MbtileId,
			File_name: rdtrMbtileItem["file_name"].(string),
			Uuid:      rdtrMbtileItem["uuid"].(string),
			Rdtr_id:   int64(rdtrMbtileItem["rdtr_id"].(int64)),
		}

		_, err2 := RdtrService.AddMbtile()
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
	ctx.AddParam("id", strconv.Itoa(int(rdtrData.Id)))
	a.GetRdtrById(ctx)
	// rr := Model.RdtrType{}
	// tx.Preload("Rdtr_mbtiles").
	// 	Preload("Rdtr_groups").
	// 	Preload("Rdtr_groups.Datas").Where("id = ?", rdtrData.Id).First(&rr)

	// ctx.JSON(200, gin.H{
	// 	"return":      rr,
	// 	"status":      "success",
	// 	"status_code": 200,
	// })
}

func (a *RdtrController) DeleteRdtr(ctx *gin.Context) {
	var RdtrService = Service.RdtrService{}
	var props = map[string]interface{}{} // Bind the request body to the newUser struct
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var propsSt struct {
		Ids []int
	}
	Helper.ToStructFromMap(props, &propsSt)
	err := RdtrService.DeleteByIds(propsSt.Ids)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Deleted Successfuly"})
}

func (a *RdtrController) ValidateMbtile(ctx *gin.Context) {
	var props struct {
		Rdtr_mbtile_ids         []int `json:"rdtr_mbtile_ids"`
		Rdtr_mbtile_ids_uncheck []int `json:"rdtr_mbtile_ids_uncheck"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rdtrMbtileService := Service.RdtrMbtileService{
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

	var mbtile_datas = []model.RdtrMbtile{}

	// First get uncheck data first
	rdtrMbtileDB := rdtrMbtileService.Gets()
	err = rdtrMbtileDB.Where("id IN ?", props.Rdtr_mbtile_ids_uncheck).Find(&mbtile_datas).Error
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
	rdtrMbtileDB = rdtrMbtileService.Gets()
	rdtrMbtileDB.Where("id IN ?", props.Rdtr_mbtile_ids_uncheck).Update("checked_at", nil)

	// Next get all check datas
	rdtrMbtileDB = rdtrMbtileService.Gets()
	err = rdtrMbtileDB.Where("id IN ?", props.Rdtr_mbtile_ids).Find(&mbtile_datas).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Iterate over the map
	for _, value2 := range mbtile_datas {
		martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/rdtr/", value2.UUID, ".mbtiles")
	}

	// Redefine again to martin config mbtiles
	martin_mbtiles["sources"] = martin_mbtile_sources
	martinMap["mbtiles"] = martin_mbtiles

	// Then set checked_at
	rdtrMbtileDB = rdtrMbtileService.Gets()
	currentTime := time.Now()
	rdtrMbtileDB.Where("id IN ?", props.Rdtr_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

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
	rdtrMbtileDB = rdtrMbtileService.Gets()
	rdtrMbtileDB.Where("id IN ?", props.Rdtr_mbtile_ids).Update("checked_at", currentTime.Format("2006-01-02 15:04:05"))

	ctx.JSON(200, gin.H{
		"return":      martinConfig,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) ValidateKml(ctx *gin.Context) {
	var props struct {
		Rdtr_group_ids []int `json:"rdtr_group_ids"`
	}
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rdtrService := Service.RdtrService{
		DB: gorm_support.DB,
	}
	asyncJObService := Service.AsynqJobService{
		DB: gorm_support.DB,
	}

	// Reset rdtr file remove validated_at
	rdtrFileService := Service.RdtrFileService{
		DB: gorm_support.DB,
	}
	rdtrFileService.Unvalidated(props.Rdtr_group_ids)

	var queue_ids []map[string]interface{}

	for i := 0; i < len(props.Rdtr_group_ids); i++ {
		fmt.Println("props.Rdtr_group_ids[i]", props.Rdtr_group_ids[i])
		taskk, err := asynq_support.NewValidateKmlTask(props.Rdtr_group_ids[i], "rdtr")
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
			"id":       props.Rdtr_group_ids[i],
			"asynq_id": gg.ID,
		})

		rdtrGroupItemData := Model.RdtrGroup{}
		err = rdtrService.GetRdtrGroups().Where("id = ?", props.Rdtr_group_ids[i]).First(&rdtrGroupItemData).Error
		if err != nil {
			log.Fatalf("could not scheudle task: %v", err)
			break
		}
		resAsynqData, err := asyncJObService.Add(Service.AsynqJobAddPayload{
			App_uuid:     rdtrGroupItemData.Uuid,
			Asynq_uuid:   gg.ID,
			Payload:      string(gg.Payload),
			Status:       asyncJObService.GetStatus().STATUS_PENDING,
			Table_name:   rdtrGroupItemData.TableName(),
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type aty2type struct {
	Line_number int
	Is_run      bool
}

func (a *RdtrController) HandleWS(ctx *gin.Context) {
	var asynq_ids = map[string]*aty2type{}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var checkStatusAsynqClose func() bool
	var checkStatusRdtrGroupClose func() bool

	defer func() {
		fmt.Println("Websocket from client is closed")
		if checkStatusRdtrGroupClose != nil {
			checkStatusRdtrGroupClose()
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
				checkStatusAsynqClose = checkAsynqStatusClosure(conn, stringSlice)
			}
		case "CHECK_VALIDATED":
			if checkStatusRdtrGroupClose == nil {
				jj, ok := dataParse["group_ids"]
				if !ok {
					log.Println("Problem check interface")
					return
				}

				// Create a new slice of int64
				int64Slice := helper.SliceOfInt64(jj)
				checkStatusRdtrGroupClose = checkStatusRdtrGroupClousure(conn, int64Slice)
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
				go func(conn *websocket.Conn, ass *aty2type) {
					for ass.Is_run {
						tailLog(conn, asynq_id, ass)
						time.Sleep(2 * time.Second)
					}
				}(conn, asynq_ids[asynq_id])
			}
		}
	}
}

func checkAsynqStatusClosure(conn *websocket.Conn, uuids []string) func() bool {
	is_loop := true
	go func(uuids []string, s *bool) {
		for *s {
			fmt.Println("checkAsynqStatusClosure - running")
			asynq_datas, err := checkStatusAsynqStatus(uuids)
			if err != nil {
				log.Println("checkAsynqStatusClosure - err", err)
				return
			}
			textT := map[string]interface{}{}
			textT["from"] = "check_asynq_status"
			textT["message"] = asynq_datas
			textTSTrng, _ := json.Marshal(textT)
			conn.WriteMessage(websocket.TextMessage, textTSTrng)
			time.Sleep(5 * time.Second)
		}
		fmt.Println("checkAsynqStatusClosure - stop")
	}(uuids, &is_loop)
	return func() bool {
		is_loop = false
		return is_loop
	}
}

func checkStatusRdtrGroupClousure(conn *websocket.Conn, ids []int64) func() bool {
	is_loop := true
	go func(ids []int64, s *bool) {
		for *s {
			var group_ids []int64 = ids
			rdtr_group_datas, err := checkStatusRdtrGroups(group_ids)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return
			}
			fmt.Println("checkStatusRdtrGroupClousure - running")
			textT := map[string]interface{}{}
			textT["from"] = "check_group"
			textT["message"] = rdtr_group_datas
			textTSTrng, _ := json.Marshal(textT)
			conn.WriteMessage(websocket.TextMessage, textTSTrng)
			time.Sleep(5 * time.Second)
		}
		fmt.Println("checkStatusRdtrGroupClousure - stop")
	}(ids, &is_loop)
	return func() bool {
		is_loop = false
		return is_loop
	}
}

func tailLog(conn *websocket.Conn, asynq_id string, asyncItem *aty2type) {
	ff, err := os.Open(fmt.Sprint("./storage/log/job/", asynq_id, ".log"))

	if err != nil {
		log.Println("tailLog :: ", err)
	}

	fmt.Println("Tail ", asynq_id, " still running")

	ggNewReader := bufio.NewReader(ff)

	// Skip lines until the desired starting line
	for i := 1; i < int((*asyncItem).Line_number); i++ {
		_, err := ggNewReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error skipping lines:", err)
			return
		}
	}

	// Create a scanner to read the file line by line from the desired starting line
	kkNewScanner := bufio.NewScanner(ggNewReader)

	// Process or print lines as needed
	for kkNewScanner.Scan() {
		(*asyncItem).Line_number = (*asyncItem).Line_number + 1
		textT := map[string]interface{}{}
		textT["from"] = asynq_id
		textT["message"] = kkNewScanner.Text()
		textTSTrng, _ := json.Marshal(textT)
		conn.WriteMessage(websocket.TextMessage, textTSTrng)
	}

	if kkNewScanner.Err() != nil {
		fmt.Println("Error reading file:", kkNewScanner.Err())
		// (*asyncItem).Is_run = false
		return
	}
}

func checkStatusRdtrGroups(ids []int64) ([]Model.RdtrGroupView, error) {
	rdtr_service := Service.RdtrService{
		DB: gorm_support.DB,
	}
	rdtr_groupModel := rdtr_service.GetRdtrGroups()
	rdtr_group_datas := []Model.RdtrGroupView{}
	err := rdtr_groupModel.Preload("Datas").
		Select("rdtr_group.*, " +
			"COALESCE((SELECT COUNT(*) FROM rdtr_file  WHERE rdtr_file.rdtr_group_id = rdtr_group.id AND rdtr_file.validated_at IS NULL),0) AS unvalidated, " +
			"COALESCE((SELECT COUNT(*) FROM rdtr_file WHERE rdtr_file.rdtr_group_id = rdtr_group.id AND rdtr_file.validated_at IS NOT NULL),0) as validated").
		Where([]int64(ids)).Find(&rdtr_group_datas).Error
	if err != nil {
		return []Model.RdtrGroupView{}, err
	}

	return rdtr_group_datas, err
}

func checkStatusAsynqStatus(uuids []string) ([]Model.AsyncJobView, error) {

	asynq_job_service := Service.AsynqJobService{
		DB: gorm_support.DB,
	}
	asynq_datas := []Model.AsyncJobView{}
	async_job_db := asynq_job_service.Gets()
	err := async_job_db.Where("asynq_uuid IN ?", uuids).Find(&asynq_datas).Error
	if err != nil {
		log.Println("checkStatusAsynqStatus - err :: ", err)
		return []Model.AsyncJobView{}, nil
	}
	return asynq_datas, nil
}
