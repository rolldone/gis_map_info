package task

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gis_map_info/app/helper/parser"
	"gis_map_info/app/helper/parser/parsermodel"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/hibiken/asynq"
)

const (
	TypeValidateKml = "validate:kml"
)

type ValidateKmlPayload struct {
	Rdtr_group_id int
	Segment_group string
}

func NewValidateKmlTask(rdtr_group_id int, segment_group string) (*asynq.Task, error) {
	payload, err := json.Marshal(ValidateKmlPayload{
		Rdtr_group_id: rdtr_group_id,
		Segment_group: segment_group,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeValidateKml, payload), nil
}

func HandleValidateKmlTask(ctx context.Context, t *asynq.Task) error {
	var p ValidateKmlPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	fmt.Printf("Validate is running :: %v", p.Rdtr_group_id)
	switch p.Segment_group {
	case "rdtr":
		tx := model.DB
		fileService := service.RdtrFileService{
			DB: tx,
		}
		rdtrFileDB := fileService.Gets(nil)
		rdtrFIleDatas := []model.RdtrFile{}
		err := rdtrFileDB.Where("rdtr_group_id = ? ", p.Rdtr_group_id).Find(&rdtrFIleDatas).Error
		if err != nil {
			fmt.Println("Validate rdtr :: ", err.Error())
		}
		ff, err := json.Marshal(rdtrFIleDatas)
		if err != nil {
			fmt.Println("Validate rdtr - marshal :: ", err.Error())
		}
		fmt.Println("\nrdtrFileDatas :: ", string(ff))
		// for i := 0; i < len(rdtrFIleDatas); i++ {
		// 	rdtrFileItem := rdtrFIleDatas[i]
		// 	generateGeoJSON("./storage/rdtr_files/" + rdtrFileItem.UUID + ".kml")
		// 	// readKmltoString("./storage/rdtr_files/" + rdtrFileItem.UUID + ".kml")
		// 	break
		// }
	case "rtrw":
	case "zlp":

	}

	// k, err :=
	// if err != nil {
	// 	log.Fatalf("Failed to parse KML: %v", err)
	// }

	// // Convert KML to GeoJSON
	// geoJSON, err := k.MarshalJSON()
	// if err != nil {
	// 	log.Fatalf("Failed to convert KML to GeoJSON: %v", err)
	// }

	// // Output GeoJSON
	// fmt.Println(string(geoJSON))
	// for i := 0; i < 100; i++ {

	// 	pubsub.NATS.Publish("foo", []byte("This is from task"))
	// }
	// : user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	return nil
}

func generateGeoJSON(path string) {
	// fileName := strings.ReplaceAll(path, ".kml", ".geojson")
	var kml = parser.UnMarshallKml(path)
	placeMarks := kml.Document.Placemark
	if parser.GetFolderRecursive(kml.Folder) != nil {
		placeMarks = parser.GetFolderRecursive(kml.Folder).Document.Placemark
	}
	for i := 0; i < len(placeMarks); i++ {
		placeMarkItem := placeMarks[i]
		description := map[string]interface{}{}
		if !reflect.ValueOf(placeMarkItem.ExtendedData).IsZero() && !reflect.ValueOf(placeMarkItem.ExtendedData.SchemaData).IsZero() {
			fmt.Println("placeMarkIteme Name", placeMarkItem.ExtendedData.SchemaData.SimpleData[0].Value)
			for i_index := 0; i_index < len(placeMarkItem.ExtendedData.SchemaData.SimpleData); i_index++ {
				simpleDataItem := placeMarkItem.ExtendedData.SchemaData.SimpleData[i_index]
				description[simpleDataItem.Name] = simpleDataItem.Value
			}
		}
		_descriptionString, err := json.Marshal(description)
		if err != nil {
			fmt.Println("json.Marshal([]byte(description)) :: ", err)
		}
		placeMarkItem.Description.Value = string(_descriptionString)
		fmt.Println("PlaceMarkerItem", placeMarkItem.Description)
		placeMarks[i] = placeMarkItem
	}
	kml.Document.Placemark = placeMarks
	newKML := parsermodel.Root{
		Folder: kml,
	}

	fileKml, err := xml.Marshal(newKML)
	if err != nil {
		fmt.Println("xml.Marshal(placeMarks) :: ", err)
	}
	_ = ioutil.WriteFile(path, []byte(fileKml), 0644)
	// var geoJson = parser.ToGeoJson(kml)
	// file, _ := json.MarshalIndent(geoJson, "", " ")

	// _ = ioutil.WriteFile(fileName, file, 0644)
}

func readKmltoString(path string) {
	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	p, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		// handle error
		fmt.Println(err.Error())
		return
	}
	fmt.Println("aaaaaaaaaaaaaaaaa :: ", p)
	// asString := string(p)
	// fmt.Println("aaaaaaaaaa", asString)
	// // var buf bytes.Buffer
	// // io.Copy(&buf, xmlFile)
	// // asString := buf.String()
	// xml := strings.NewReader(asString)
	// json, err := xj.Convert(xml)
	// if err != nil {
	// 	panic("That's embarrassing...")
	// }
	// fileName := strings.ReplaceAll(path, ".kml", ".geojson")

	// // file, _ := json.MarshalIndent(json.String(), "", " ")

	// _ = ioutil.WriteFile(fileName, []byte(json.String()), 0644)
	// fmt.Println(json.String())
}
