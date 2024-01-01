package admin

import (
	"errors"
	"fmt"
	Helper "gis_map_info/app/helper"
	Model "gis_map_info/app/model"
	Service "gis_map_info/app/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	if err := rdtrDB.Preload("Rdtr_groups").Limit(limit).Offset(offset).Order("updated_at DESC").Find(&rdtrDatas).Error; err != nil {
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
	if err := rdtrDB.Preload("Rdtr_groups").First(&rdtrData).Error; err != nil {
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
			errMsg := fmt.Errorf("Validation Error: Field %s has invalid value: %v", err.Field(), err.Value())
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
	tx := Model.DB.Begin()
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
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		tx.Rollback()
		return
	}
	rdtrGroups := props.Rdtr_groups
	fmt.Println(rdtrGroups)
	for i := 0; i < len(rdtrGroups); i++ {
		rdtrGroupItem, _ := rdtrGroups[i].(map[string]interface{})
		rdtrGroupItem["rdtr_id"] = rdtrData.Id
		rdtrGroupData := RdtrService.RdtrGroupAddType
		rdtrGroupData.Name = rdtrGroupItem["name"].(string)
		rdtrGroupData.Rdtr_id = int64(rdtrGroupItem["rdtr_id"].(float64))
		rdtrGroupData.Asset_key = rdtrGroupItem["asset_key"].(string)
		_rdtrGroupItem_properties, _ := Helper.GetValue(rdtrGroupItem["properties"], Helper.MapEmpty).(map[string]interface{})
		rdtrGroupData.Properties = _rdtrGroupItem_properties
		err = validate.Struct(rdtrGroupData)
		if err != nil {
			var errArr []error
			var errValidators = err.(validator.ValidationErrors)
			for i := 0; i < len(errValidators); i++ {
				errValidItem := errValidators[i]
				errMess := fmt.Errorf("Validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
				errArr = append(errArr, errMess)
			}
			err = errors.Join(errArr...)
			break
		}
		_, err2 := RdtrService.AddGroup(rdtrGroupData)
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
	rr := Model.RdtrType{}
	tx.Preload("Rdtr_groups").Where("id = ?", rdtrData.Id).First(&rr)
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
			errMess := fmt.Errorf("Validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
			errArr = append(errArr, errMess)
		}
		err = errors.Join(errArr...)
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err,
		})
	}

	tx := Model.DB.Begin()
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
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
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
				errMess := fmt.Errorf("Validation Error: Field %s has invalid value: %v", errValidItem.Field(), errValidItem.Value())
				errArr = append(errArr, errMess)
			}
			err = errors.Join(errArr...)
			break
		}
		_, err2 := RdtrService.AddGroup(rdtrGroupData)
		if err2 != nil {
			err = err2
			break
		}
	}
	if err != nil {
		tx.Rollback()
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	rr := Model.RdtrType{}
	tx.Preload("Rdtr_groups").Where("id = ?", rdtrData.Id).First(&rr)
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rr,
		"status":      "success",
		"status_code": 200,
	})
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
