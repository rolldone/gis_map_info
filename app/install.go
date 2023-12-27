package app

import (
	Model "gis_map_info/app/model"
	Seeder "gis_map_info/db/seeders"
)

func Install() {

	Model.ConnectDatabase()

	var InstallLocation = Seeder.InstallLocation{}
	InstallLocation.Install()
}
