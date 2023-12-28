package service

import Rtrw "gis_map_info/app/model"

type RtrwService struct {
}

func (c *RtrwService) Gets(props interface{}) []Rtrw.RtrwType {
	return []Rtrw.RtrwType{}
}

func (c *RtrwService) GetByUUId(uuid string) Rtrw.RtrwType {
	return Rtrw.RtrwType{}
}

func (c *RtrwService) Add(props interface{}) Rtrw.RtrwType {
	return Rtrw.RtrwType{}
}

func (c *RtrwService) Update(props interface{}) Rtrw.RtrwType {
	return Rtrw.RtrwType{}
}

func (c *RtrwService) Delete(arr []int) bool {
	return true
}

func (c *RtrwService) GetGroupsByRtrwId(rtrw_id int) []Rtrw.RtrwGroup {
	return []Rtrw.RtrwGroup{}
}

func (c *RtrwService) AddGroup(props interface{}) Rtrw.RtrwGroup {
	return Rtrw.RtrwGroup{}
}

func (c *RtrwService) UpdateGroup(props interface{}) Rtrw.RtrwGroup {
	return Rtrw.RtrwGroup{}
}

func (c *RtrwService) DeleteGroupByRtrwId(rtrw_id int) bool {
	return true
}
