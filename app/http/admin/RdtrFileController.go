package admin

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	Service "gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
)

type RdtrFileController struct{}

func (c *RdtrFileController) GetByUUID(ctx *gin.Context) {

	tx := gorm_support.DB
	rdtrFileService := Service.RdtrFileService{
		DB: tx,
	}
	rdtrData, err := rdtrFileService.GetByUUID(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
			"FF":          ctx.DefaultQuery("uuid", ""),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status":      "success",
		"status_code": 200,
		"return":      rdtrData,
	})
}
func (c *RdtrFileController) Add(ctx *gin.Context) {
	var formData struct {
		Uuid string
		File *multipart.FileHeader `form:"file" validate:"required"`
	}
	if err := ctx.ShouldBind(&formData); err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      "Bad request",
		})
		return
	}

	if err := ctx.ShouldBindUri(&formData); err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      "Bad request",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(formData)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
	}

	// Create directory if it doesn't exist
	uploadDir := "./storage/rdtr_files"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			// ctx.String(http.StatusInternalServerError, "Failed to create directory")
			ctx.JSON(400, gin.H{
				"status":      "error",
				"status_code": 400,
				"return":      "Bad request",
			})
			return
		}
	}

	// Generate a new UUID
	uuidValue, _ := uuid.NewRandom()
	formData.Uuid = uuidValue.String()

	// Get file extension
	fileExtension := filepath.Ext(formData.File.Filename)
	fileName := formData.Uuid + fileExtension
	rdtr_file := filepath.Join(uploadDir, fileName)
	if err := ctx.SaveUploadedFile(formData.File, rdtr_file); err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	tx := gorm_support.DB.Begin()
	rdtrFIleService := Service.RdtrFileService{
		DB: tx,
	}
	rdTrFileAdd := rdtrFIleService.RdtrFileAdd
	rdTrFileAdd.Uuid = formData.Uuid
	rdtrFile, err := rdtrFIleService.Add(rdTrFileAdd)
	if err != nil {
		fmt.Println("Error - KI2M9VCAMVKFMV :: ", err)
		tx.Rollback()
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	rrr, err := rdtrFIleService.GetById(int(rdtrFile.Id))
	if err != nil {
		fmt.Println("Error - IM2IMC9SDJM2 :: ", err)
		tx.Rollback()
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	tx.Commit()
	ctx.JSON(200, gin.H{
		"status":      "success",
		"status_code": 200,
		"return":      rrr,
	})
}
