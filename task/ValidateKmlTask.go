package task

import (
	"context"
	"encoding/json"
	"fmt"
	"gis_map_info/app/model"
	"gis_map_info/app/service"

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
