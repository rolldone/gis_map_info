package task

import (
	"context"
	"encoding/json"
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/pubsub"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	service "gis_map_info/app/service"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
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
	// fmt.Printf("Validate is running :: %v", p.Rdtr_group_id)
	switch p.Segment_group {
	case "rdtr":
		tx := model.DB
		fileService := service.RdtrFileService{
			DB: tx,
		}
		geoJsonService := service.RdtrGeojsonService{
			DB: tx,
		}
		err := geoJsonService.DeleteByRdtrGroupId(int64(p.Rdtr_group_id))
		if err != nil {
			fmt.Println("DeleteByRdtrGroupId(int64(p.Rdtr_group_id)) :: ", err)
		}
		rdtrFileDB := fileService.Gets(nil)
		rdtrFIleDatas := []model.RdtrFile{}
		err = rdtrFileDB.Preload("Rdtr_group").Preload("Rdtr").Where("rdtr_group_id = ? ", p.Rdtr_group_id).Find(&rdtrFIleDatas).Error
		if err != nil {
			fmt.Println("Validate rdtr :: ", err.Error())
		}
		ff, err := json.Marshal(rdtrFIleDatas)
		if err != nil {
			fmt.Println("Validate rdtr - marshal :: ", err.Error())
		}
		fmt.Println("\nrdtrFileDatas :: ", string(ff))
		for i := 0; i < len(rdtrFIleDatas); i++ {
			rdtrFileItem := rdtrFIleDatas[i]
			payload := map[string]interface{}{}
			payload["uuid"] = rdtrFileItem.UUID
			payload["download"] = os.Getenv("APP_HOST") + "/api/admin/rdtr_file/assets/" + string(payload["uuid"].(string)+".kml")
			payloadString, _ := json.Marshal(payload)
			pubsub.NATS.Publish("convert_kml_geojson", payloadString)
			uuidItem := payload["uuid"].(string)
			ubsubcribeProcess, err := pubsub.NATS.Subscribe(uuidItem+"_process", func(msg *nats.Msg) {
				fmt.Printf("\nReceived a message process: %s\n", string(msg.Data))
			})
			if err != nil {
				fmt.Println("Validate rdtr - pubsub.NATS.Subscribe(uuidItem_process :: ", err.Error())
				continue
			}
			unSubscribeFinish, err := pubsub.NATS.SubscribeSync(uuidItem + "_done")
			if err != nil {
				log.Fatal(err)
				continue
			}

			// Counter to track processed messages
			messageCount := 0
			maxMessages := 60 // Maximum number of timeout process

			var msg *nats.Msg
			// Infinite loop to listen for messages
			for messageCount < maxMessages {
				// Wait for a message
				msg, err = unSubscribeFinish.NextMsg(1 * time.Second) // Timeout after 5 seconds if no message
				if err != nil {
					if err == nats.ErrTimeout {
						fmt.Println("Timed out waiting for a message.", messageCount)
						messageCount++
						continue
					}
					log.Fatal(err)
				}

				// Process the received message
				fmt.Printf("Received message: %s\n", msg.Data)
				break
			}
			if err := ubsubcribeProcess.Unsubscribe(); err != nil {
				fmt.Printf("Error unsubscribing process: %v\n", err)
				continue
			}
			if err := unSubscribeFinish.Unsubscribe(); err != nil {
				fmt.Printf("Error unsubscribing done: %v\n", err)
				continue
			}
			if msg != nil {
				resData := map[string]interface{}{}
				err := json.Unmarshal(msg.Data, &resData)
				if err != nil {
					fmt.Println("json.Unmarshal(msg.Data, &mesData) :: ", err)
					continue
				}
				returnData := resData["return"].(map[string]interface{})
				filePath := "./storage/rdtr_files/" + (string(payload["uuid"].(string))) + ".json"
				downloadFileJSON(filePath, string(returnData["download"].(string)))
				var groupProperty map[string]interface{}
				err = json.Unmarshal([]byte(rdtrFileItem.Rdtr_group.Properties), &groupProperty)
				if err != nil {
					fmt.Println("Unmarshal([]byte(rdtrFileItem.Rdtr_group.Properties), &groupProperty) :: ", err)
				}
				generateGeometry(int(rdtrFileItem.Rdtr_id.Int64), int(rdtrFileItem.Rdtr_group_id.Int64), int(rdtrFileItem.Id), groupProperty, filePath)
			}

		}
	case "rtrw":
	case "zlp":

	}
	return nil
}

func downloadFileJSON(dir_path string, path string) {
	fmt.Println("downloadFileJSON", path)

	resp, err := http.Get(path)
	if err != nil {
		fmt.Println("Error fetching the file:", err)
		return
	}
	defer resp.Body.Close()
	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get the file. Status code: %d\n", resp.StatusCode)
		return
	}
	// Create a new file where you want to save the downloaded content
	out, err := os.Create(dir_path)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer out.Close()

	// Copy the content from the HTTP response to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error copying the content:", err)
		return
	}

	fmt.Println("File downloaded successfully!")
}

func generateGeometry(rdtr_id int, rdtr_group_id int, rdtr_file_id int, property map[string]interface{}, path string) {
	rdtrGeojsonService := service.RdtrGeojsonService{
		DB: model.DB,
	}
	readFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("os.ReadFile(path) :: ", err)
	}
	var dataArray []map[string]interface{}

	// Unmarshal the JSON content into the slice of maps
	if err := json.Unmarshal(readFile, &dataArray); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Print the slice to see the JSON content in Go format
	for a := 0; a < len(dataArray); a++ {
		arrItem := dataArray[a]
		if arrItem["properties"] == nil {
			fmt.Println("kosong")
		} else {
			arrProperty := arrItem["properties"].(map[string]interface{})
			// Merge map2 into map1
			arrProperty["fill"] = property["fill_color"]
			arrItem["properties"] = arrProperty
		}
		rdtrGeojsonPayloadAdd := rdtrGeojsonService.RdtrGeojsonAdd
		uuidValue, _ := uuid.NewRandom()
		rdtrGeojsonPayloadAdd.Uuid = uuidValue.String()
		propertiesByte, _ := json.Marshal(arrItem["properties"])
		rdtrGeojsonPayloadAdd.Properties = propertiesByte
		rdtrGeojsonPayloadAdd.Rdtr_file_id = int64(rdtr_file_id)
		rdtrGeojsonPayloadAdd.Rdtr_group_id = int64(rdtr_group_id)
		rdtrGeojsonPayloadAdd.Rdtr_id = int64(rdtr_id)
		geometryByte, _ := json.Marshal(arrItem["geometry"])
		rdtrGeojsonPayloadAdd.Geojson = string(geometryByte) // pq.StringArray{string(geometryByte)}
		resRdtrGeoJsonAdd, err := rdtrGeojsonService.Add(rdtrGeojsonPayloadAdd)
		if err != nil {
			fmt.Println("rdtrGeojsonService.Add(rdtrGeojsonPayloadAdd) :: ", err)
		}
		fmt.Println("resRdtrGeoJsonAdd", resRdtrGeoJsonAdd)
	}

}
