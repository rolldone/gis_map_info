package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"gis_map_info/app/model"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserServiceConstruct(DB *gorm.DB) UserService {
	gg := UserService{
		db: DB,
	}
	return gg
}

type userStatus struct {
	ACTIVE   string
	INACTIVE string
}

type AddPayload_UserService struct {
	Name                  string
	Username              string
	Email                 string
	Password              *string
	Password_confirmation *string
	Status                string
}

type UserService struct {
	db *gorm.DB
}

func (c *UserService) Add(props AddPayload_UserService) (*model.User, error) {
	userModel := model.User{}
	userModel.Uuid = uuid.New().String()
	userModel.Email = props.Email
	userModel.Username = props.Username
	userModel.Name = props.Name
	userModel.Salt = uuid.NewString()
	userModel.Passkey = c.GeneratePassword(*props.Password, userModel.Salt)
	userModel.Status = props.Status
	err := c.db.Model(&model.User{}).Create(&userModel).Error
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

type UpdatePayload_UserService struct {
	AddPayload_UserService
	Uuid string
}

func (c *UserService) Update(props UpdatePayload_UserService) (*model.User, error) {
	userModel := model.User{}
	err := c.db.Model(&model.User{}).Where("uuid = ?", props.Uuid).First(&userModel).Error
	if err != nil {
		return nil, err
	}
	userModel.Username = props.Username
	userModel.Email = props.Email
	userModel.Name = props.Name
	val := reflect.ValueOf(props.Password)
	if !val.IsNil() {
		if *props.Password != *props.Password_confirmation {
			err := errors.New("password and password confirmation not same")
			return &userModel, err
		}
		userModel.Passkey = c.GeneratePassword(*props.Password, userModel.Salt)
	}
	userModel.Status = props.Status
	err = c.db.Model(&model.User{}).Save(&userModel).Error
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

func (c *UserService) Delete(uuids []string) error {
	err := c.db.Model(&model.User{}).Where("uuid IN ?", uuids).Delete(model.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *UserService) Gets() *gorm.DB {
	return c.db.Model(&model.User{})
}

func (c *UserService) GetStatus() userStatus {
	return userStatus{
		ACTIVE:   "active",
		INACTIVE: "inactive",
	}
}

func (c *UserService) GetByUUID(uuid string) (*model.UserView, error) {
	gg := &model.UserView{}
	err := c.db.Model(&model.User{}).Where("uuid = ?", uuid).First(gg).Error
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func (c *UserService) GeneratePassword(password string, salt string) string {
	var saltedText = fmt.Sprintf("text: '%s', salt: %s", password, salt)
	sha := sha1.New()
	sha.Write([]byte(saltedText))
	var encrypted = hex.EncodeToString(sha.Sum(nil))
	return string(encrypted)
}
