package admin

import (
	"fmt"
	Helper "gis_map_info/app/helper"
	Model "gis_map_info/app/model"
	Service "gis_map_info/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	if err := rdtrDB.First(&rdtrData).Error; err != nil {
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
	tx := Model.DB.Begin()
	var RdtrService = Service.RdtrService{}
	RdtrService.DB = tx
	var props = map[string]interface{}{} // Bind the request body to the newUser struct
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rdtrData, err := RdtrService.Add(props)
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
	rdtrGroups, _ := props["rdtr_groups"].([]interface{})
	if rdtrGroups == nil {
		rdtrGroups = []interface{}{}
	}
	fmt.Println(rdtrGroups)
	for i := 0; i < len(rdtrGroups); i++ {
		rdtrGroupItem, _ := rdtrGroups[i].(map[string]interface{})
		rdtrGroupItem["rdtr_id"] = rdtrData.Id
		_, err2 := RdtrService.AddGroup(rdtrGroups[i])
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
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rdtrData,
		"status":      "success",
		"status_code": 200,
	})
}

func (a *RdtrController) UpdateRdtr(ctx *gin.Context) {
	tx := Model.DB.Begin()
	var RdtrService = Service.RdtrService{}
	RdtrService.DB = tx
	var props = map[string]interface{}{} // Bind the request body to the newUser struct
	if err := ctx.ShouldBindJSON(&props); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rdtrData, err := RdtrService.Update(props)
	if err != nil {
		fmt.Println("Error:", err)
		ctx.JSON(500, gin.H{
			"status":      "error",
			"status_code": 500,
			"return":      err.Error(),
		})
		return
	}
	rdtrGroups, _ := props["rdtr_groups"].([]interface{})
	if rdtrGroups == nil {
		rdtrGroups = []interface{}{}
	}
	fmt.Println(rdtrGroups)
	RdtrService.DeleteGroupByRdtrId(int(rdtrData.Id))
	for i := 0; i < len(rdtrGroups); i++ {
		rdtrGroupItem, _ := rdtrGroups[i].(map[string]interface{})
		rdtrGroupItem["rdtr_id"] = rdtrData.Id
		_, err2 := RdtrService.AddGroup(rdtrGroups[i])
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
	tx.Commit()
	ctx.JSON(200, gin.H{
		"return":      rdtrData,
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
