package asynq_support

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client

type Task struct {
}

func Init() *asynq.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	redisDB := os.Getenv("REDIS_DB")

	// Create asynq config for asynq cli
	configCLi := map[string]interface{}{
		"uri":      fmt.Sprint(redisHost, ":", redisPort),
		"db":       redisDB,
		"password": redisPassword,
	}

	f, _ := os.Create(".asynq.json")
	configString, _ := json.Marshal(configCLi)
	f.WriteString(string(configString))
	f.Close()

	// Convert string to int
	DB, _ := strconv.Atoi(redisDB)

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisHost + ":" + redisPort, Password: redisPassword,
		DB: DB,
	})

	Client = client

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort, Password: redisPassword},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 1,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// mux.HandleFunc(task.TypeEmailDelivery, task.HandleEmailDeliveryTask)
	// mux.Handle(task.TypeImageResize, task.NewImageProcessor())
	mux.HandleFunc(TypeValidateKml, HandleValidateKmlTask)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	return client
}
