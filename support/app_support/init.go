package app_support

import (
	"os"
)

var App *AppSupport

type AppSupport struct {
	App_host         string
	Time_zone        string
	File_zlp_path    string
	File_rtrw_path   string
	File_rdtr_path   string
	Mbtile_zlp_path  string
	Mbtile_rtrw_path string
	Mbtile_rdtr_path string
}

func Init() {

	App_host := os.Getenv("APP_HOST")
	Time_zone := os.Getenv("TIME_ZONE")

	App = &AppSupport{
		App_host:         App_host,
		Time_zone:        Time_zone,
		File_zlp_path:    "./storage/zlp_files",
		File_rtrw_path:   "./storage/rtrw_files",
		File_rdtr_path:   "./storage/rdtr_files",
		Mbtile_zlp_path:  "./sub_app/martin/mbtiles/zlp",
		Mbtile_rtrw_path: "./sub_app/martin/mbtiles/rtrw",
		Mbtile_rdtr_path: "./sub_app/martin/mbtiles/rdtr",
	}
}
