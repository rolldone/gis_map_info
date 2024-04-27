package service

import (
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type ZlpFileService struct {
	DB *gorm.DB
	// Custom own payload
	ZlpFileAdd    random_92821m23unjnvjfv
	ZlpFileUpdate random_o28mamivm286
}

type random_92821m23unjnvjfv struct {
	Uuid string
}

func (c *ZlpFileService) Add(props random_92821m23unjnvjfv) (Model.ZlpFile, error) {
	zlpFile := Model.ZlpFile{}
	zlpFile.UUID = props.Uuid
	if err := c.DB.Create(&zlpFile).Error; err != nil {
		return Model.ZlpFile{}, err
	}
	return zlpFile, nil
}

type random_o28mamivm286 struct {
	Id           int64
	Zlp_group_id int64
	Zlp_id       int64
}

func (c *ZlpFileService) Update(props random_o28mamivm286) (Model.ZlpFile, error) {
	zlpFile := Model.ZlpFile{}
	c.DB.Where("id = ?", props.Id).First(&zlpFile)
	zlpFile.Zlp_group_id = Model.NullInt64{
		Valid: true,
		Int64: props.Zlp_group_id,
	}
	zlpFile.Zlp_id = Model.NullInt64{
		Valid: true,
		Int64: props.Zlp_id,
	}
	err := c.DB.Save(&zlpFile).Error
	if err != nil {
		return Model.ZlpFile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpFileService) DeleteById(id int) bool {
	return true
}

func (c *ZlpFileService) Gets(props interface{}) *gorm.DB {
	zlpFiles := Model.ZlpFile{}
	return c.DB.Model(&zlpFiles)
}

func (c *ZlpFileService) GetByUUID(uuid string) (Model.ZlpFile, error) {
	zlpFile := Model.ZlpFile{}
	err := c.DB.Where("uuid = ?", uuid).First(&zlpFile).Error
	if err != nil {
		return Model.ZlpFile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpFileService) GetById(id int) (Model.ZlpFile, error) {
	zlpFile := Model.ZlpFile{}
	err := c.DB.Where("id = ?", id).First(&zlpFile).Error
	if err != nil {
		return Model.ZlpFile{}, err
	}
	return zlpFile, nil
}

func (c *ZlpFileService) Unvalidated(ids []int) error {
	err := c.DB.Model(model.ZlpFile{}).Where("zlp_group_id IN ?", ids).Update("validated_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}
