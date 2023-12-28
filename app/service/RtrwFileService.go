package service

import (
	Model "gis_map_info/app/model"
)

type RtrwFileService struct {
}

func (c *RtrwFileService) add(props interface{}) Model.RtrwFile {
	return Model.RtrwFile{}
}

func (c *RtrwFileService) deleteById(id int) bool {
	return true
}

func (c *RtrwFileService) gets(props interface{}) []Model.RtrwFile {
	return []Model.RtrwFile{}
}

func (c *RtrwFileService) getById(id int) Model.RtrwFile {
	return Model.RtrwFile{}
}
