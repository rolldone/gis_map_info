package service

import (
	Model "gis_map_info/app/model"
)

type RdtrFileService struct {
}

func (c *RdtrFileService) add(props interface{}) Model.RdtrFile {
	return Model.RdtrFile{}
}

func (c *RdtrFileService) deleteById(id int) bool {
	return true
}

func (c *RdtrFileService) gets(props interface{}) []Model.RdtrFile {
	return []Model.RdtrFile{}
}

func (c *RdtrFileService) getById(id int) Model.RdtrFile {
	return Model.RdtrFile{}
}
