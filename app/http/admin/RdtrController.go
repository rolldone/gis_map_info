package admin

import (
	"encoding/json"
	"errors"
	"fmt"
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

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"github.com/nxadm/tail"
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
	if err := rdtrDB.Preload("Rdtr_groups").
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
	if err := rdtrDB.Preload("Rdtr_groups", func(db *gorm.DB) *gorm.DB {
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
				ll, err3 := rdtrFileService.Update(rdtrFileProps)
				if err3 != nil {
					err = err3
				}
				fmt.Println("rdtrFileItem ", ll)

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
		"return":      props,
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

type TheClient struct {
	Id           string
	Conn         *websocket.Conn
	loopingCheck func()
}

var conWs []TheClient

func (a *RdtrController) HandleWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ttail := []*tail.Tail{}
	rdtr_group_uuids := []string{}
	closeAllTail := func() {
		for _, v := range ttail {
			fmt.Println("File are closed")
			v.Stop()
		}
	}

	defer func() {
		for i, v := range conWs {
			if conn == v.Conn {
				fmt.Printf("websocket %v from client is closed \n", i)
				closeAllTail()
				v.Conn.Close()
				conWs = Helper.RemoveIndex(conWs, i)
				break
			}
		}
	}()

	clientThe := TheClient{
		Id:   uuid.New().String(),
		Conn: conn,
	}
	conWs = append(conWs, clientThe)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("You are join"))
		if err != nil {
			log.Println(err)
			return
		}

		if clientThe.loopingCheck == nil {
			ids := []int64{}
			err = json.Unmarshal([]byte(message), &ids)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				conn.Close()
				return
			}
			clientThe.loopingCheck = func() {
				go func(conn *websocket.Conn, ids []int64) {
					for {
						isExist := false
						for _, v := range conWs {
							if v.Conn == conn {
								isExist = true
							}
						}
						if !isExist {
							break
						}
						time.Sleep(3 * time.Second)
						rdtr_group_datas, err := checkStatusRdtrGroups(ids)
						if err != nil {
							conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
							conn.Close()
							return
						}
						rdtr_group_datas_str, err := json.Marshal(rdtr_group_datas)
						if err != nil {
							conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
							conn.Close()
							return
						}
						// fmt.Println("Jalan terus")
						// fmt.Println("rdtr_group_datas :: ", string(rdtr_group_datas_str))
						conn.WriteMessage(websocket.TextMessage, []byte(rdtr_group_datas_str))
					}
					fmt.Println("Exit goroutine :: Websocket")
				}(conn, ids)

				for _, v := range ids {
					// Get the rdtr_group
					rdtrService := service.RdtrService{
						DB: gorm_support.DB,
					}
					rdtr_groupItem := model.RdtrGroup{}
					rdtr_groupDB := rdtrService.GetRdtrGroups()
					err = rdtr_groupDB.Where("id = ?", v).First(&rdtr_groupItem).Error
					if err != nil {
						log.Fatalf(err.Error())
						return
					}
					// Get the asynq_job
					asynqJobService := service.AsynqJobService{
						DB: gorm_support.DB,
					}
					asynqJobData := model.AsynqJob{}
					asynqJob_DB := asynqJobService.Gets()
					err = asynqJob_DB.Where("app_uuid = ?", rdtr_groupItem.Uuid).Order("created_at DESC").First(&asynqJobData).Error
					if err != nil {
						log.Fatalf(err.Error())
						return
					}
					// Create a tail
					t, err := tail.TailFile(
						fmt.Sprint("./storage/log/job/", asynqJobData.Asynq_uuid, ".log"), tail.Config{Follow: true, ReOpen: true})
					if err != nil {
						log.Fatalf(err.Error())
						return
					}
					ttail = append(ttail, t)
					rdtr_group_uuids = append(rdtr_group_uuids, asynqJobData.App_uuid)
				}
				for i, v := range ttail {
					go func(conn *websocket.Conn, v *tail.Tail, s string) {
						// Print the text of each received line
						for line := range v.Lines {
							textT := map[string]interface{}{}
							textT["from"] = s
							textT["message"] = line.Text
							textTSTrng, _ := json.Marshal(textT)
							conn.WriteMessage(websocket.TextMessage, textTSTrng)
						}
						fmt.Println("Read file Done :: ", s)
					}(conn, v, rdtr_group_uuids[i])
				}
			}
			clientThe.loopingCheck()
		} else {
			// fmt.Println("Udah di register")
			// This is allready register
		}
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
