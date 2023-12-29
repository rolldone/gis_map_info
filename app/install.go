package app

import (
	Seeder "gis_map_info/db/seeders"
)

func Install() {

	var InstallLocation = Seeder.InstallLocation{}
	InstallLocation.Install()
}
