package middleware

import (
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var contextTopicAdmin = "admin_data"

func AdminTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/dashboard/auth/login" {
			c.Next()
			return
		}

		_, exists := c.Request.Header["Authorization"]
		if !exists && len(c.Request.Header["Authorization"]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"status_code": http.StatusUnauthorized,
				"return":      "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := c.Request.Header["Authorization"][0]
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		userService := service.UserServiceConstruct(gorm_support.DB)
		jwtToken, err := userService.CheckJWTTOken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"status_code": http.StatusUnauthorized,
				"return":      "Unauthorized",
			})
			c.Abort()
			return
		}
		fmt.Println("jwtTokenClaimns", jwtToken)
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"status_code": http.StatusUnauthorized,
				"return":      "Unauthorized",
			})
			c.Abort()
			return
		}

		userData := model.UserView{}
		userDB := userService.Gets()
		err = userDB.Where("uuid = ?", claims["uuid"]).First(&userData).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"status_code": http.StatusUnauthorized,
				"return":      "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Set(contextTopicAdmin, userData)
		c.Next()
		// For testing enable this code
		// c.JSON(http.StatusAccepted, gin.H{
		// 	"status":      "error",
		// 	"status_code": http.StatusAccepted,
		// 	"return":      jwtToken,
		// })
		// c.Done()
	}
}

func GetAdminAccess(ctx *gin.Context) *model.UserView {
	// Get personal data from token middleware
	project_data, ok := ctx.Get(contextTopicAdmin)
	if !ok {
		return nil
	}
	personalData := project_data.(model.UserView)
	return &personalData
}
