package seeders

import (
	"encoding/json"
	"fmt"
	Model "gis_map_info/app/model"
	"gis_map_info/support/gorm_support"
	"log"
	"os"
	"strconv"
	"strings"
)

type RegProvinceJson struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Alt_name  string  `json:"alt_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RegRegencyJson struct {
	Id          string  `json:"id"`
	Province_id string  `json:"province_id"`
	Name        string  `json:"name"`
	Alt_name    string  `json:"alt_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type RegDistrctJson struct {
	Id         string  `json:"id"`
	Regency_id string  `json:"regency_id"`
	Name       string  `json:"name"`
	Alt_name   string  `json:"alt_name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type RegVillageJson struct {
	Id          string  `json:"id"`
	District_id string  `json:"district_id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type InstallLocation struct{}

func (a *InstallLocation) Install() {

	regProvinceModel := Model.RegProvince{}
	regProvinceDatas := gorm_support.DB.Find(&regProvinceModel)
	fmt.Println("Install province")
	if regProvinceDatas.RowsAffected == 0 {
		data, err := os.ReadFile("db/seeders/json/reg_provinces_latlng.json")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var province []RegProvinceJson
		// Unmarshal the JSON data into the slice of Person

		if err := json.Unmarshal(data, &province); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Loop through the slice and print each person's name and age
		for _, p := range province {
			// fmt.Printf("Name: %s, AltName: %s\n", p.Name, p.Alt_name)
			pId, err := strconv.Atoi(p.Id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regProvinceModel.AltName = p.Alt_name
			regProvinceModel.Id = int64(pId)
			regProvinceModel.Name = p.Name
			regProvinceModel.Latitude = p.Latitude
			regProvinceModel.Longitude = p.Longitude
			gorm_support.DB.Create(&regProvinceModel)
		}
		var regProvinceDatas []Model.RegProvince
		gorm_support.DB.Find(&regProvinceDatas)
		// Loop through the retrieved records
		// for _, regProvince := range regProvinceDatas {
		// 	// Print each record (or perform other operations as needed)
		// 	fmt.Printf("ID: %d, Name: %s\n", regProvince.Id, regProvince.Name)
		// }
	}

	fmt.Println("Install regency")
	var count int64
	regRegencyModel := Model.RegRegency{}
	gorm_support.DB.Model(&regRegencyModel).Count(&count)
	if count == 0 {
		data, err := os.ReadFile("db/seeders/json/reg_regencies_latlng.json")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var regency []RegRegencyJson
		// Unmarshal the JSON data into the slice of Person

		if err := json.Unmarshal(data, &regency); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Loop through the slice and print each person's name and age
		for _, p := range regency {
			// fmt.Printf("Name: %s, AltName: %s\n", p.Name, p.Alt_name)
			pId, err := strconv.Atoi(p.Id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regRegencyModel := Model.RegRegency{}
			regRegencyModel.Id = int64(pId)
			regRegencyModel.FullName = p.Name
			regRegencyModel.Latitude = p.Latitude
			regRegencyModel.Longitude = p.Longitude
			regRegencyModel.AltName = p.Alt_name
			// Check if the name contains "KABUPATEN"
			if strings.Contains(strings.ToUpper(p.Name), "KABUPATEN") {
				regRegencyModel.Name = strings.ReplaceAll(strings.ToUpper(p.Name), "KABUPATEN ", "")
				regRegencyModel.Type = "Kabupaten"
			} else if strings.Contains(strings.ToUpper(p.Name), "KOTA") {
				regRegencyModel.Name = strings.ReplaceAll(strings.ToUpper(p.Name), "KOTA ", "")
				regRegencyModel.Type = "Kota"
			}
			pId, err = strconv.Atoi(p.Province_id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regRegencyModel.RegProvinceID = int64(pId)
			gorm_support.DB.Create(&regRegencyModel)
		}
		var regRegencyDatas []Model.RegRegency
		gorm_support.DB.Find(&regRegencyDatas)
		// Loop through the retrieved records
		// for _, regRegency := range regRegencyDatas {
		// 	// Print each record (or perform other operations as needed)
		// 	fmt.Printf("ID: %d, Name: %s\n", regRegency.Id, regRegency.Name)
		// }
	}

	fmt.Println("Install district")
	regDistrictModel := Model.RegDistrict{}
	gorm_support.DB.Model(&regDistrictModel).Count(&count)
	if count == 0 {
		data, err := os.ReadFile("db/seeders/json/reg_districts_latlng.json")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var district []RegDistrctJson
		// Unmarshal the JSON data into the slice of Person

		if err := json.Unmarshal(data, &district); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Loop through the slice and print each person's name and age
		for _, p := range district {
			// fmt.Printf("Name: %s, AltName: %s\n", p.Name, p.Alt_name)

			pId, err := strconv.Atoi(p.Id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regDistrictModel.Id = int64(pId)
			regDistrictModel.Latitude = p.Latitude
			regDistrictModel.Longitude = p.Longitude
			regDistrictModel.Name = p.Name
			regDistrictModel.AltName = p.Alt_name
			pId, err = strconv.Atoi(p.Regency_id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regDistrictModel.RegRegencyID = int64(pId)
			gorm_support.DB.Create(&regDistrictModel)
		}
		var regDistrictDatas []Model.RegRegency
		gorm_support.DB.Find(&regDistrictDatas)
		// Loop through the retrieved records
		// for _, regDistrict := range regDistrictDatas {
		// 	// Print each record (or perform other operations as needed)
		// 	fmt.Printf("ID: %d, Name: %s\n", regDistrict.Id, regDistrict.Name)
		// }
	}

	fmt.Println("Install village")
	regVillageModel := Model.RegVillage{}
	gorm_support.DB.Model(&regVillageModel).Count(&count)
	if count == 0 {
		data, err := os.ReadFile("db/seeders/json/reg_villages_latlng.json")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var village []RegVillageJson
		// Unmarshal the JSON data into the slice of Person

		if err := json.Unmarshal(data, &village); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Loop through the slice and print each person's name and age
		for _, p := range village {
			// fmt.Printf("Name: %s, AltName: %s\n", p.Name, p.Alt_name)

			pId, err := strconv.Atoi(p.Id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regVillageModel.Id = int64(pId)
			regVillageModel.Name = p.Name
			regVillageModel.Latitude = p.Latitude
			regVillageModel.Longitude = p.Longitude
			pId, err = strconv.Atoi(p.District_id)
			if err != nil {
				// Handle error if conversion fails
				fmt.Println("Failed to convert string to int:", err)
				os.Exit(1)
				return
			}
			regVillageModel.RegDistrictID = int64(pId)
			gorm_support.DB.Create(&regVillageModel)
		}
		var regVillageDatas []Model.RegRegency
		gorm_support.DB.Find(&regVillageDatas)
		// Loop through the retrieved records
		// for _, regVillage := range regVillageDatas {
		// 	// Print each record (or perform other operations as needed)
		// 	fmt.Printf("ID: %d, Name: %s\n", regVillage.Id, regVillage.Name)
		// }
	}
}
