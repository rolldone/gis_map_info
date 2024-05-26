package admin

import (
	"fmt"
	Service "gis_map_info/app/service"
	"gis_map_info/support/app_support"
	"gis_map_info/support/gorm_support"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type ZlpMbtileController struct {
	DB *gorm.DB
}

func (c *ZlpMbtileController) GetbyUUID(ctx *gin.Context) {
	tx := gorm_support.DB
	zlpMbtileService := Service.ZlpMbtileService{
		DB: tx,
	}
	zlpData, err := zlpMbtileService.GetByUUID(ctx.Param("uuid"))
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
		"return":      zlpData,
	})
}

func (c *ZlpMbtileController) Add(ctx *gin.Context) {
	var formData struct {
		Uuid            string
		File            *multipart.FileHeader `form:"file" validate:"required"`
		Asset_key       string                `form:"asset_key"`
		Reg_province_id int64                 `form:"reg_province_id"`
		Zlp_group_id    int64                 `form:"zlp_group_id"`
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
	uploadDir := app_support.App.Mbtile_zlp_path
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
	zlp_mbtile_file := filepath.Join(uploadDir, fileName)
	if err := ctx.SaveUploadedFile(formData.File, zlp_mbtile_file); err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	tx := gorm_support.DB.Begin()
	zlpFIleService := Service.ZlpMbtileService{
		DB: tx,
	}
	zlpFIleService.AddPayload = &struct {
		Uuid            string
		File_name       string
		Asset_key       string
		Reg_province_id int64
		Zlp_group_id    int64
	}{
		Uuid:            formData.Uuid,
		File_name:       formData.File.Filename,
		Asset_key:       formData.Asset_key,
		Reg_province_id: formData.Reg_province_id,
		Zlp_group_id:    formData.Zlp_group_id,
	}
	zlpMbtile, err := zlpFIleService.Add()
	if err != nil {
		fmt.Println("Error - KI2M9VC5MVKFMV :: ", err)
		tx.Rollback()
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	rrr, err := zlpFIleService.GetById(int(zlpMbtile.Id))
	if err != nil {
		fmt.Println("Error - IM25MC9SDJM2 :: ", err)
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

func (c *ZlpMbtileController) GetMartinConfig(ctx *gin.Context) {

	martinConfig, err := os.ReadFile("./sub_app/martin/config.yaml")
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	var martinJson map[string]interface{}
	err = yaml.Unmarshal(martinConfig, &martinJson)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":      "error",
			"status_code": 400,
			"return":      err.Error(),
		})
		return
	}
	fmt.Println("martin_config", martinJson)
	ctx.JSON(200, gin.H{
		"status":      "success",
		"status_code": 200,
		"return":      martinJson,
	})
}
