package service

import (
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RtrwFileService struct {
	DB *gorm.DB
	// Custom own payload
	RtrwFileAdd    random_92823m23unjnvjfv
	RtrwFileUpdate random_o28mamivm289
}

type random_92823m23unjnvjfv struct {
	Uuid string
}

func (c *RtrwFileService) Add(props random_92823m23unjnvjfv) (Model.RtrwFile, error) {
	rtrwFile := Model.RtrwFile{}
	rtrwFile.UUID = props.Uuid
	if err := c.DB.Create(&rtrwFile).Error; err != nil {
		return Model.RtrwFile{}, err
	}
	return rtrwFile, nil
}

type random_o28mamivm289 struct {
	Id            int64
	Rtrw_group_id int64
	Rtrw_id       int64
}

func (c *RtrwFileService) Update(props random_o28mamivm289) (Model.RtrwFile, error) {
	rtrwFile := Model.RtrwFile{}
	c.DB.Where("id = ?", props.Id).First(&rtrwFile)
	rtrwFile.Rtrw_group_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_group_id,
	}
	rtrwFile.Rtrw_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rtrw_id,
	}
	err := c.DB.Save(&rtrwFile).Error
	if err != nil {
		return Model.RtrwFile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwFileService) DeleteById(id int) bool {
	return true
}

func (c *RtrwFileService) Gets(props interface{}) *gorm.DB {
	rtrwFiles := Model.RtrwFile{}
	return c.DB.Model(&rtrwFiles)
}

func (c *RtrwFileService) GetByUUID(uuid string) (Model.RtrwFile, error) {
	rtrwFile := Model.RtrwFile{}
	err := c.DB.Where("uuid = ?", uuid).First(&rtrwFile).Error
	if err != nil {
		return Model.RtrwFile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwFileService) GetById(id int) (Model.RtrwFile, error) {
	rtrwFile := Model.RtrwFile{}
	err := c.DB.Where("id = ?", id).First(&rtrwFile).Error
	if err != nil {
		return Model.RtrwFile{}, err
	}
	return rtrwFile, nil
}

func (c *RtrwFileService) Unvalidated(ids []int) error {
	err := c.DB.Model(model.RtrwFile{}).Where("rtrw_group_id IN ?", ids).Update("validated_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}
