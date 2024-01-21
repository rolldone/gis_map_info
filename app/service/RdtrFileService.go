package service

import (
	"gis_map_info/app/model"
	Model "gis_map_info/app/model"

	"gorm.io/gorm"
)

type RdtrFileService struct {
	DB *gorm.DB
	// Custom own payload
	RdtrFileAdd    random_92923m23unjnvjfv
	RdtrFileUpdate random_o29mamivm289
}

type random_92923m23unjnvjfv struct {
	Uuid string
}

func (c *RdtrFileService) Add(props random_92923m23unjnvjfv) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	rdtrFile.UUID = props.Uuid
	if err := c.DB.Create(&rdtrFile).Error; err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

type random_o29mamivm289 struct {
	Id            int64
	Rdtr_group_id int64
	Rdtr_id       int64
}

func (c *RdtrFileService) Update(props random_o29mamivm289) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	c.DB.Where("id = ?", props.Id).First(&rdtrFile)
	rdtrFile.Rdtr_group_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_group_id,
	}
	rdtrFile.Rdtr_id = Model.NullInt64{
		Valid: true,
		Int64: props.Rdtr_id,
	}
	err := c.DB.Save(&rdtrFile).Error
	if err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrFileService) DeleteById(id int) bool {
	return true
}

func (c *RdtrFileService) Gets(props interface{}) *gorm.DB {
	rdtrFiles := Model.RdtrFile{}
	return c.DB.Model(&rdtrFiles)
}

func (c *RdtrFileService) GetByUUID(uuid string) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	err := c.DB.Where("uuid = ?", uuid).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrFileService) GetById(id int) (Model.RdtrFile, error) {
	rdtrFile := Model.RdtrFile{}
	err := c.DB.Where("id = ?", id).First(&rdtrFile).Error
	if err != nil {
		return Model.RdtrFile{}, err
	}
	return rdtrFile, nil
}

func (c *RdtrFileService) Unvalidated(ids []int) error {
	err := c.DB.Model(model.RdtrFile{}).Where("rdtr_group_id IN ?", ids).Update("validated_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}
