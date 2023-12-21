package model

type RdtrType struct {
	id          int    `json:id`
	uuid        string `json:"uuid"`
	province_id int    `json:province_id`
	regency_id  int    `json:regency_id`
	district_id int    `json:district_id`
	village_id  int    `json:village_id`
}
