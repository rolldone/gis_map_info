package admin

import (
	"fmt"
	"gis_map_info/app/helper"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"gis_map_info/support/gorm_support"
	"gis_map_info/support/redis_support"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthController struct {
}

var APPLICATION_NAME = "My Simple JWT App"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("the secret of kalimdor")

type MyClaims struct {
	jwt.RegisteredClaims
	Uuid string `json:"uuid"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	props := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.BindJSON(&props)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	userData := model.User{}
	userService := service.UserServiceConstruct(gorm_support.DB)
	userdataDB := userService.Gets()
	err = userdataDB.Where("username = ?", props.Username).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	passwordEncrypte := userService.GeneratePassword(props.Password, userData.Salt)

	if passwordEncrypte == userData.Passkey {
		signedToken, err := GenerateToken(userData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":      "error",
				"status_code": http.StatusBadRequest,
				"return":      err.Error(),
			})
			ctx.Done()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":      "Success",
			"status_code": http.StatusOK,
			"return":      signedToken,
		})
		ctx.Done()
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":      "error",
		"status_code": http.StatusBadRequest,
		"return":      "Username or password are incorect",
	})
	ctx.Done()
}

func (c *AuthController) Register(ctx *gin.Context) {
	props := struct {
		Name                  string `json:"name"`
		Username              string `json:"username"`
		Email                 string `json:"email"`
		Password              string `json:"password"`
		Password_confirmation string `json:"password_confirmation"`
	}{}
	err := ctx.BindJSON(&props)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)

	userData, err := userService.Add(service.AddPayload_UserService{
		Name:     props.Name,
		Username: props.Username,
		Email:    props.Email,
		Password: &props.Password,
		Status:   userService.GetStatus().ACTIVE,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	userDataView, err := userService.GetByUUID(userData.Uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      userDataView,
	})
	ctx.Done()
}

func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	props := struct {
		Uuid                  string  `json:"uuid"`
		Name                  string  `json:"name"`
		Username              string  `json:"username"`
		Email                 string  `json:"email"`
		Status                string  `json:"status"`
		Password              *string `json:"password,omitempty"`
		Password_confirmation *string `json:"password_confirmation,omitempty"`
		Current_password      *string `json:"current_password,omitempty"`
	}{}

	err := ctx.BindJSON(&props)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}

	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)

	userData := model.User{}
	userdataDB := userService.Gets()
	err = userdataDB.Where("uuid = ?", props.Uuid).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}

	if props.Current_password != nil || props.Password != nil || props.Password_confirmation != nil {
		passwordEncrypte := userService.GeneratePassword(*props.Current_password, userData.Salt)
		if passwordEncrypte != userData.Passkey {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":      "error",
				"status_code": http.StatusBadRequest,
				"return":      "Username or password are incorect",
			})
			ctx.Done()
			return
		}
	}

	userNewData, err := userService.Update(service.UpdatePayload_UserService{
		Uuid: props.Uuid,
		AddPayload_UserService: service.AddPayload_UserService{
			Name:                  props.Name,
			Username:              props.Username,
			Email:                 props.Email,
			Password:              props.Password,
			Password_confirmation: props.Password_confirmation,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}

	userDataView, err := userService.GetByUUID(userNewData.Uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      userDataView,
	})
	ctx.Done()
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	props := struct {
		Email string `json:"email"`
	}{}
	err := ctx.BindJSON(&props)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	userData := model.User{}
	userService := service.UserServiceConstruct(gorm_support.DB)
	userDB := userService.Gets()
	err = userDB.Where("email = ?", props.Email).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	key := helper.RandStringBytes(6)
	redis_support.RedisClient.Set(ctx, key, userData.Uuid, time.Duration(time.Minute*30))
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      key,
		"uuid":        userData.Uuid,
	})
	ctx.Done()
}

func (c *AuthController) RecoveryPassword(ctx *gin.Context) {
	props := struct {
		Password              string `json:"password"`
		Password_confirmation string `json:"password_confirmation"`
		Key                   string `json:"key"`
	}{}
	err := ctx.BindJSON(&props)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}

	val, err := redis_support.RedisClient.Get(ctx, props.Key).Result()
	fmt.Println("val", val)
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      "They key is expired",
		})
		ctx.Done()
		return
	}
	userData := model.User{}
	DB := gorm_support.DB.Begin()
	userService := service.UserServiceConstruct(DB)
	userDB := userService.Gets()
	err = userDB.Where("uuid = ?", val).First(&userData).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		return
	}
	userService.Update(service.UpdatePayload_UserService{
		AddPayload_UserService: service.AddPayload_UserService{
			Name:     userData.Name,
			Username: userData.Username,
			Email:    userData.Email,
			Password: &props.Password,
			Status:   userData.Status,
		},
		Uuid: userData.Uuid,
	})
	userDataView, err := userService.GetByUUID(userData.Uuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      err.Error(),
		})
		ctx.Done()
		DB.Rollback()
		return
	}
	DB.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      userDataView,
	})
	ctx.Done()
}

func (c *AuthController) Auth(ctx *gin.Context) {
	admin_data, ok := ctx.Get("admin_data")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"status_code": http.StatusBadRequest,
			"return":      "User data is empty",
		})
		ctx.Done()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "Success",
		"status_code": http.StatusOK,
		"return":      admin_data.(model.UserView),
	})
	ctx.Done()
}

func (c *AuthController) Logout(ctx *gin.Context) {

}

func GenerateToken(userData model.User) (*string, error) {
	now := time.Now()
	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: &jwt.NumericDate{now.Add(time.Hour * 24)},
		},
		Uuid: userData.Uuid,
	}
	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)
	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	return &signedToken, err
}

func CheckJWTTOken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("signing method invalid")
		}

		return JWT_SIGNATURE_KEY, nil
	})
	return token, err
}
