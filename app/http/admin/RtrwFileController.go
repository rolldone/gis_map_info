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

type RtrwFileController struct{}

func (c *RtrwFileController) GetByUUID(ctx *gin.Context) {

	tx := gorm_support.DB
	rtrwFileService := Service.RtrwFileService{
		DB: tx,
	}
	rtrwData, err := rtrwFileService.GetByUUID(ctx.Param("uuid"))
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
		"return":      rtrwData,
	})
}
func (c *RtrwFileController) Add(ctx *gin.Context) {
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
	uploadDir := "./storage/rtrw_files"
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
	rtrw_file := filepath.Join(uploadDir, fileName)
	if err := ctx.SaveUploadedFile(formData.File, rtrw_file); err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	tx := gorm_support.DB.Begin()
	rtrwFIleService := Service.RtrwFileService{
		DB: tx,
	}
	rdTrFileAdd := rtrwFIleService.RtrwFileAdd
	rdTrFileAdd.Uuid = formData.Uuid
	rtrwFile, err := rtrwFIleService.Add(rdTrFileAdd)
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
	rrr, err := rtrwFIleService.GetById(int(rtrwFile.Id))
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
