package parser

import (
	"encoding/xml"
	"fmt"
	model "gis_map_info/app/helper/parser/parsermodel"
	"io/ioutil"
	"os"
)

func UnMarshallKml(filePath string) model.Folder {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var root model.Root

	err = xml.Unmarshal(byteValue, &root)
	if err != nil {
		fmt.Println(err, "--", filePath)
		panic(err)
	}

	return root.Folder
}

// Func helper
func GetFolderRecursive(folder *model.Folder) *model.Folder {
	if folder != nil {
		return GetFolderRecursive(folder.Folder)
	}
	return nil
}
