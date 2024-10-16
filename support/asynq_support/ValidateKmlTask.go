package asynq_support

import (
	"context"
	"encoding/json"
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/support/gorm_support"
	"gis_map_info/support/nats_support"
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
	Table_group_id int
	Segment_group  string
}

func NewValidateKmlTask(table_group_id int, segment_group string) (*asynq.Task, error) {
	payload, err := json.Marshal(ValidateKmlPayload{
		Table_group_id: table_group_id,
		Segment_group:  segment_group,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeValidateKml, payload), nil
}

func HandleValidateKmlTask(ctx context.Context, t *asynq.Task) error {

	taskID, _ := asynq.GetTaskID(ctx)

	asyncJObService := service.AsynqJobService{
		DB:         gorm_support.DB,
		Async_uuid: &taskID,
	}

	// Call construct
	asyncJObService.Construct(taskID)
	// Check the job is possible continue or not
	// If not set status stopped
	if !asyncJObService.IsPossibleContinue() {
		fmt.Println("Async UUid ", taskID, " is stopped")
		ctx.Done()
		return nil
	}

	// Record first status
	asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_PROCESS, "Job is processed\n")

	var p ValidateKmlPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// fmt.Printf("Validate is running :: %v", p.Rdtr_group_id)
	switch p.Segment_group {
	case "rdtr":
		tx := gorm_support.DB
		fileService := service.RdtrFileService{
			DB: tx,
		}

		geoJsonService := service.RdtrGeojsonService{
			DB: tx,
		}

		// Delete exist rdtr geojson datas
		err := geoJsonService.DeleteByRdtrGroupId(int64(p.Table_group_id))
		if err != nil {
			fmt.Println("DeleteByRdtrGroupId(int64(p.Rdtr_group_id)) :: ", err)
		}

		rdtrFileDB := fileService.Gets(nil)
		rdtrFIleDatas := []model.RdtrFile{}
		err = rdtrFileDB.Preload("Rdtr_group").Preload("Rdtr").Where("rdtr_group_id = ? ", p.Table_group_id).Find(&rdtrFIleDatas).Error
		if err != nil {
			fmt.Println("Validate rdtr :: ", err.Error())
		}

		for i := 0; i < len(rdtrFIleDatas); i++ {
			rdtrFileItem := rdtrFIleDatas[i]
			payload := map[string]interface{}{}
			payload["uuid"] = rdtrFileItem.UUID
			payload["download"] = os.Getenv("APP_HOST") + "/api/admin/rdtr_file/assets/" + string(payload["uuid"].(string)+".kml")
			payloadString, _ := json.Marshal(payload)
			// Publish the payload string
			nats_support.NATS.Publish("convert_kml_geojson", payloadString)
			uuidItem := payload["uuid"].(string)

			// And listen each event listener
			ubsubcribeProcess, err := nats_support.NATS.Subscribe(uuidItem+"_process", func(msg *nats.Msg) {
				fmt.Printf("\n ubsubcribeProcess : %s\n", string(msg.Data))
			})
			if err != nil {
				fmt.Println("Validate rdtr - nats_support.NATS.Subscribe(uuidItem_process :: ", err.Error())
				continue
			}

			// Listen syncronize for _done event listener
			unSubscribeFinish, err := nats_support.NATS.SubscribeSync(uuidItem + "_done")
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
						asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_WAITING, "Job is waiting \n")
						messageCount++
						continue
					}
					log.Fatal(err)
				}

				// Process the received message
				fmt.Printf("\n unSubscribeFinish: %s\n", msg.Data)
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

				isWarning := generateRdtrGeometry(func(message string) {
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_PROCESS, message)
				}, int(rdtrFileItem.Rdtr_id.Int64), int(rdtrFileItem.Rdtr_group_id.Int64), int(rdtrFileItem.Id), groupProperty, filePath)
				if !isWarning {
					now := time.Now()
					fmt.Println("Validated ", rdtrFileItem.Id)
					err := fileService.DB.Model(&model.RdtrFile{}).Where("id = ?", rdtrFileItem.Id).Update("validated_at", now).Error
					if err != nil {
						fmt.Println("Update(validated_at, now) :: ", err)
					}
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_COMPLETED, "Job is completed \n")
				} else {
					fmt.Println("Unvalidated", rdtrFileItem.Id)
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is failed \n")
				}
			} else {
				asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is timeout \n")
			}
		}
	case "rtrw":
		tx := gorm_support.DB
		fileService := service.RtrwFileService{
			DB: tx,
		}

		geoJsonService := service.RtrwGeojsonService{
			DB: tx,
		}

		// Delete exist rtrw geojson datas
		err := geoJsonService.DeleteByRtrwGroupId(int64(p.Table_group_id))
		if err != nil {
			fmt.Println("DeleteByRtrwGroupId(int64(p.Table_group_id)) :: ", err)
		}

		rtrwFileDB := fileService.Gets(nil)
		rtrwFIleDatas := []model.RtrwFile{}
		err = rtrwFileDB.Preload("Rtrw_group").Preload("Rtrw").Where("rtrw_group_id = ? ", p.Table_group_id).Find(&rtrwFIleDatas).Error
		if err != nil {
			fmt.Println("Validate rtrw :: ", err.Error())
		}

		for i := 0; i < len(rtrwFIleDatas); i++ {
			rtrwFileItem := rtrwFIleDatas[i]
			payload := map[string]interface{}{}
			payload["uuid"] = rtrwFileItem.UUID
			payload["download"] = os.Getenv("APP_HOST") + "/api/admin/rtrw_file/assets/" + string(payload["uuid"].(string)+".kml")
			payloadString, _ := json.Marshal(payload)
			// Publish the payload string
			nats_support.NATS.Publish("convert_kml_geojson", payloadString)
			uuidItem := payload["uuid"].(string)

			// And listen each event listener
			ubsubcribeProcess, err := nats_support.NATS.Subscribe(uuidItem+"_process", func(msg *nats.Msg) {
				fmt.Printf("\n ubsubcribeProcess : %s\n", string(msg.Data))
			})
			if err != nil {
				fmt.Println("Validate rtrw - nats_support.NATS.Subscribe(uuidItem_process :: ", err.Error())
				continue
			}

			// Listen syncronize for _done event listener
			unSubscribeFinish, err := nats_support.NATS.SubscribeSync(uuidItem + "_done")
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
						asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_WAITING, "Job is waiting \n")
						messageCount++
						continue
					}
					log.Fatal(err)
				}

				// Process the received message
				fmt.Printf("\n unSubscribeFinish: %s\n", msg.Data)
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
				filePath := "./storage/rtrw_files/" + (string(payload["uuid"].(string))) + ".json"
				downloadFileJSON(filePath, string(returnData["download"].(string)))
				var groupProperty map[string]interface{}
				err = json.Unmarshal([]byte(rtrwFileItem.Rtrw_group.Properties), &groupProperty)
				if err != nil {
					fmt.Println("Unmarshal([]byte(rtrwFileItem.Rtrw_group.Properties), &groupProperty) :: ", err)
				}

				isWarning := generateRtrwGeometry(func(message string) {
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_PROCESS, message)
				}, int(rtrwFileItem.Rtrw_id.Int64), int(rtrwFileItem.Rtrw_group_id.Int64), int(rtrwFileItem.Id), groupProperty, filePath)
				if !isWarning {
					now := time.Now()
					fmt.Println("Validated ", rtrwFileItem.Id)
					err := fileService.DB.Model(&model.RtrwFile{}).Where("id = ?", rtrwFileItem.Id).Update("validated_at", now).Error
					if err != nil {
						fmt.Println("Update(validated_at, now) :: ", err)
					}
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_COMPLETED, "Job is completed \n")
				} else {
					fmt.Println("Unvalidated", rtrwFileItem.Id)
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is failed \n")
				}
			} else {
				asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is timeout \n")
			}
		}
	case "zlp":
		tx := gorm_support.DB
		fileService := service.ZlpFileService{
			DB: tx,
		}

		geoJsonService := service.ZlpGeojsonService{
			DB: tx,
		}

		// Delete exist zlp geojson datas
		err := geoJsonService.DeleteByZlpGroupId(int64(p.Table_group_id))
		if err != nil {
			fmt.Println("DeleteByZlpGroupId(int64(p.Table_group_id)) :: ", err)
		}

		zlpFileDB := fileService.Gets(nil)
		zlpFIleDatas := []model.ZlpFile{}
		err = zlpFileDB.Preload("Zlp_group").Preload("Zlp").Where("zlp_id != 0").Where("zlp_group_id = ? ", p.Table_group_id).Find(&zlpFIleDatas).Error
		if err != nil {
			fmt.Println("Validate zlp :: ", err.Error())
		}

		for i := 0; i < len(zlpFIleDatas); i++ {
			zlpFileItem := zlpFIleDatas[i]
			payload := map[string]interface{}{}
			payload["uuid"] = zlpFileItem.UUID
			payload["download"] = os.Getenv("APP_HOST") + "/api/admin/zlp_file/assets/" + string(payload["uuid"].(string)+".kml")
			payloadString, _ := json.Marshal(payload)
			// Publish the payload string
			nats_support.NATS.Publish("convert_kml_geojson", payloadString)
			uuidItem := payload["uuid"].(string)

			// And listen each event listener
			ubsubcribeProcess, err := nats_support.NATS.Subscribe(uuidItem+"_process", func(msg *nats.Msg) {
				fmt.Printf("\n ubsubcribeProcess : %s\n", string(msg.Data))
			})
			if err != nil {
				fmt.Println("Validate zlp - nats_support.NATS.Subscribe(uuidItem_process :: ", err.Error())
				continue
			}

			// Listen syncronize for _done event listener
			unSubscribeFinish, err := nats_support.NATS.SubscribeSync(uuidItem + "_done")
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
						asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_WAITING, "Job is waiting \n")
						messageCount++
						continue
					}
					log.Fatal(err)
				}

				// Process the received message
				fmt.Printf("\n unSubscribeFinish: %s\n", msg.Data)
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
				filePath := "./storage/zlp_files/" + (string(payload["uuid"].(string))) + ".json"
				downloadFileJSON(filePath, string(returnData["download"].(string)))
				var groupProperty map[string]interface{}
				err = json.Unmarshal([]byte(zlpFileItem.Zlp_group.Properties), &groupProperty)
				if err != nil {
					fmt.Println("Unmarshal([]byte(zlpFileItem.Zlp_group.Properties), &groupProperty) :: ", err)
				}

				isWarning := generateZlpGeometry(func(message string) {
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_PROCESS, message)
				}, int(zlpFileItem.Zlp_id.Int64), int(zlpFileItem.Zlp_group_id.Int64), int(zlpFileItem.Id), groupProperty, filePath)
				if !isWarning {
					now := time.Now()
					fmt.Println("Validated ", zlpFileItem.Id)
					err := fileService.DB.Model(&model.ZlpFile{}).Where("id = ?", zlpFileItem.Id).Update("validated_at", now).Error
					if err != nil {
						fmt.Println("Update(validated_at, now) :: ", err)
					}
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_COMPLETED, "Job is completed \n")
				} else {
					fmt.Println("Unvalidated", zlpFileItem.Id)
					// Record status
					asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is failed \n")
				}
			} else {
				asyncJObService.UpdateByAsynqUUID(asyncJObService.GetStatus().STATUS_FAILED, "Job is timeout \n")
			}
		}
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

func generateRdtrGeometry(gg func(message string), rdtr_id int, rdtr_group_id int, rdtr_file_id int, property map[string]interface{}, path string) bool {
	rdtr_ids_warnings := false
	rdtrGeojsonService := service.RdtrGeojsonService{
		DB: gorm_support.DB,
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
	_total := len(dataArray)
	for a := 0; a < _total; a++ {
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
		time.Sleep(100 * time.Millisecond)
		currentTime := time.Now()
		UnixNano := currentTime.UnixNano()
		rdtrGeojsonPayloadAdd.Order_number = UnixNano + int64(a)
		rdtrGeojsonPayloadAdd.Properties = propertiesByte
		rdtrGeojsonPayloadAdd.Rdtr_file_id = int64(rdtr_file_id)
		rdtrGeojsonPayloadAdd.Rdtr_group_id = int64(rdtr_group_id)
		rdtrGeojsonPayloadAdd.Rdtr_id = int64(rdtr_id)
		geometryByte, _ := json.Marshal(arrItem["geometry"])
		rdtrGeojsonPayloadAdd.Geojson = string(geometryByte) // pq.StringArray{string(geometryByte)}
		resRdtrGeoJsonAdd, err := rdtrGeojsonService.Add(rdtrGeojsonPayloadAdd)
		if err != nil {
			fmt.Println("rdtrGeojsonService.Add(rdtrGeojsonPayloadAdd) :: ", err)
			gg("Job Failed Process")
			rdtr_ids_warnings = true
			continue
		}
		gg(fmt.Sprint("Job process :: ", path, " :: ", _total, " > ", a, "\n"))
		fmt.Println(path, " - ", _total, " - ", a)
		// Ignore it
		_ = resRdtrGeoJsonAdd

	}
	return rdtr_ids_warnings
}

func generateRtrwGeometry(gg func(message string), rtrw_id int, rtrw_group_id int, rtrw_file_id int, property map[string]interface{}, path string) bool {
	rtrw_ids_warnings := false
	rtrwGeojsonService := service.RtrwGeojsonService{
		DB: gorm_support.DB,
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
	_total := len(dataArray)
	for a := 0; a < _total; a++ {
		arrItem := dataArray[a]
		if arrItem["properties"] == nil {
			fmt.Println("kosong")
		} else {
			arrProperty := arrItem["properties"].(map[string]interface{})
			// Merge map2 into map1
			arrProperty["fill"] = property["fill_color"]
			arrItem["properties"] = arrProperty
		}
		rtrwGeojsonPayloadAdd := rtrwGeojsonService.RtrwGeojsonAdd
		uuidValue, _ := uuid.NewRandom()
		rtrwGeojsonPayloadAdd.Uuid = uuidValue.String()
		propertiesByte, _ := json.Marshal(arrItem["properties"])
		time.Sleep(100 * time.Millisecond)
		currentTime := time.Now()
		UnixNano := currentTime.UnixNano()
		rtrwGeojsonPayloadAdd.Order_number = UnixNano + int64(a)
		rtrwGeojsonPayloadAdd.Properties = propertiesByte
		rtrwGeojsonPayloadAdd.Rtrw_file_id = int64(rtrw_file_id)
		rtrwGeojsonPayloadAdd.Rtrw_group_id = int64(rtrw_group_id)
		rtrwGeojsonPayloadAdd.Rtrw_id = int64(rtrw_id)
		geometryByte, _ := json.Marshal(arrItem["geometry"])
		rtrwGeojsonPayloadAdd.Geojson = string(geometryByte) // pq.StringArray{string(geometryByte)}
		resRtrwGeoJsonAdd, err := rtrwGeojsonService.Add(rtrwGeojsonPayloadAdd)
		if err != nil {
			fmt.Println("rtrwGeojsonService.Add(rtrwGeojsonPayloadAdd) :: ", err)
			gg("Job Failed Process")
			rtrw_ids_warnings = true
			continue
		}
		gg(fmt.Sprint("Job process :: ", path, " :: ", _total, " > ", a, "\n"))
		fmt.Println(path, " - ", _total, " - ", a)
		// Ignore it
		_ = resRtrwGeoJsonAdd

	}
	return rtrw_ids_warnings
}

func generateZlpGeometry(gg func(message string), zlp_id int, zlp_group_id int, zlp_file_id int, property map[string]interface{}, path string) bool {
	zlp_ids_warnings := false
	zlpGeojsonService := service.ZlpGeojsonService{
		DB: gorm_support.DB,
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
	_total := len(dataArray)
	for a := 0; a < _total; a++ {
		arrItem := dataArray[a]
		if arrItem["properties"] == nil {
			fmt.Println("kosong")
		} else {
			arrProperty := arrItem["properties"].(map[string]interface{})
			// Merge map2 into map1
			arrProperty["fill"] = property["fill_color"]
			arrItem["properties"] = arrProperty
		}
		zlpGeojsonPayloadAdd := zlpGeojsonService.ZlpGeojsonAdd
		uuidValue, _ := uuid.NewRandom()
		zlpGeojsonPayloadAdd.Uuid = uuidValue.String()
		propertiesByte, _ := json.Marshal(arrItem["properties"])
		time.Sleep(100 * time.Millisecond)
		currentTime := time.Now()
		UnixNano := currentTime.UnixNano()
		zlpGeojsonPayloadAdd.Order_number = UnixNano + int64(a)
		zlpGeojsonPayloadAdd.Properties = propertiesByte
		zlpGeojsonPayloadAdd.Zlp_file_id = int64(zlp_file_id)
		zlpGeojsonPayloadAdd.Zlp_group_id = int64(zlp_group_id)
		zlpGeojsonPayloadAdd.Zlp_id = int64(zlp_id)
		geometryByte, _ := json.Marshal(arrItem["geometry"])
		zlpGeojsonPayloadAdd.Geojson = string(geometryByte) // pq.StringArray{string(geometryByte)}
		resZlpGeoJsonAdd, err := zlpGeojsonService.Add(zlpGeojsonPayloadAdd)
		if err != nil {
			fmt.Println("zlpGeojsonService.Add(zlpGeojsonPayloadAdd) :: ", err)
			gg("Job Failed Process")
			zlp_ids_warnings = true
			continue
		}
		gg(fmt.Sprint("Job process :: ", path, " :: ", _total, " > ", a, "\n"))
		fmt.Println(path, " - ", _total, " - ", a)
		// Ignore it
		_ = resZlpGeoJsonAdd

	}
	return zlp_ids_warnings
}
