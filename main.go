package main

import (
	"fmt"
	"gis_map_info/app"
	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
	Model "gis_map_info/app/model"
	"gis_map_info/app/pubsub"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	Task "gis_map_info/task"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	Model.ConnectDatabase()

	app.Install()

	// Initialize Gin's default router
	router := gin.Default()

	router.Use(gin.Logger())

	router.Use(CORSMiddleware())

	// Define a route handler
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	api := router.Group("/api")
	{
		// Admin side
		admin := api.Group("/admin")
		{
			dashboard := AdminController.Dashboard()
			admin.GET("/dashboard/rdtr-periode", gin.HandlerFunc(dashboard.GetRdtrPeriode))
			admin.GET("/dashboard/rtrw-periode", gin.HandlerFunc(dashboard.GetRtrwPeriode))

			rdtrController := &AdminController.RdtrController{}
			admin.GET("/zone_rdtr/get/:id/view", rdtrController.GetRdtrById)
			admin.GET("/zone_rdtr/gets/paginate", rdtrController.GetRdtrsPaginate)
			admin.GET("/zone_rdtr/gets", rdtrController.GetRdtrs)
			admin.POST("/zone_rdtr/add", rdtrController.AddRdtr)
			admin.POST("/zone_rdtr/update", rdtrController.UpdateRdtr)
			admin.POST("/zone_rdtr/delete", rdtrController.DeleteRdtr)
			admin.POST("/zone_rdtr/validate", rdtrController.Validate)

			rdtrFileController := &AdminController.RdtrFileController{}
			admin.POST("/rdtr_file/add", rdtrFileController.Add)
			admin.GET("/rdtr_file/get/:uuid", rdtrFileController.GetByUUID)
			admin.POST("/rdtr_file/delete", rdtrController.DeleteRdtr)
			admin.Static("/rdtr_file/assets", "./storage/rdtr_files")

			rtrwController := &AdminController.RtrwController{}
			admin.GET("/zone_rtrw/rtrws", rtrwController.GetRtrws)
			admin.GET("/zone_rtrw/:uuid/view", rtrwController.GetRtrwByUUId)
			admin.POST("/zone_rtrw/add", rtrwController.AddRtrw)
			admin.POST("/zone_rtrw/update", rtrwController.UpdateRtrw)
			admin.POST("/zone_rtrw/delete", rtrwController.DeleteRtrw)

			zlpController := &AdminController.ZLPController{}
			admin.GET("/zone_land_price/zlps", zlpController.GetZLPs)
			admin.GET("/zone_land_price/:uuid/view", zlpController.GetZLPByUUId)
			admin.POST("/zone_land_price/new", zlpController.AddZLP)
			admin.POST("/zone_land_price/update", zlpController.UpdateZLP)
			admin.POST("/zone_land_price/delete", zlpController.DeleteZLP)

			regLocationController := &AdminController.RegLocationController{}
			admin.GET("/reg_location/province/provinces", regLocationController.GetProvinces)
			admin.GET("/reg_location/regency/regencies", regLocationController.GetRegencies)
			admin.GET("/reg_location/district/districts", regLocationController.GetDistricts)
			admin.GET("/reg_location/village/villages", regLocationController.GetVillages)

		}

		// Front side

		frontRdtrController := FrontController.RdtrController()
		api.GET("rdtr/rdtrs", frontRdtrController.Gets)
		api.GET("rdtr/:id/view", frontRdtrController.GetByUUID)

		frontRtrwController := FrontController.RtrwController()
		api.GET("rtrw/rtrws", frontRtrwController.Gets)
		api.GET("rtrw/:id/view", frontRtrwController.GetByUUID)
	}

	InitNats()
	TestingRunningTask()

	go func() {
		// Start the server
		router.Run(":8080")
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Shutdown Gin server (Gin does not provide a built-in shutdown method, so you may handle this as needed)
	// For example, you can use http.Server for more control over the Gin server and shutdown it gracefully.

	log.Println("Servers gracefully stopped.")
}

func TestingRunningTask() {

	fmt.Println("Running")
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort, Password: redisPassword})
	defer client.Close()

	// This is real case validate kml
	task, err := Task.NewValidateKmlTask(161, "rdtr")
	if err != nil {
		log.Fatalf("Could not schedule task : %v", err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not scheudle task: %v", err)
	}

	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	// task, err = Task.NewEmailDeliveryTask(42, "some:template:id")
	// if err != nil {
	// 	log.Fatalf("could not schedule task: %v", err)
	// }

	// info, err = client.Enqueue(task)
	// if err != nil {
	// 	log.Fatalf("could not schedule task: %v", err)
	// }
	// log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	// info, err = client.Enqueue(task, asynq.ProcessIn(20*time.Second))
	// if err != nil {
	// 	log.Fatalf("could not schedule task: %v", err)
	// }
	// log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	// task, err = Task.NewIMageResizeTask("https://example.com/myassets/image.jpg")
	// if err != nil {
	// 	log.Fatalf("could not create task: %v", err)
	// }
	// info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(1*time.Minute))
	// if err != nil {
	// 	log.Fatalf("could not enqueue task: %v", err)
	// }
	// log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort, Password: redisPassword},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
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
	mux.HandleFunc(Task.TypeEmailDelivery, Task.HandleEmailDeliveryTask)
	mux.Handle(Task.TypeImageResize, Task.NewImageProcessor())
	mux.HandleFunc(Task.TypeValidateKml, Task.HandleValidateKmlTask)
	// ...register other handlers...
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()
}

func InitNats() {
	_, err := pubsub.ConnectPubSub()
	if err != nil {
		fmt.Println("Nats error :: ", err.Error())
	} else {
		// Simple Async Subscriber
		pubsub.NATS.Subscribe("foo", func(m *nats.Msg) {
			fmt.Printf("\nReceived a message: %s\n", string(m.Data))
		})
		go func() {
			timer := time.After(5 * time.Second)
			<-timer
			// Simple Publisher
			pubsub.NATS.Publish("foo", []byte("Hello World NATS"))
		}()
	}

}

func InitSocketClient() {
	// uri := "http://127.0.0.1:8000"

	go func() {
		timer := time.After(2 * time.Second)
		// Wait for the signal
		fmt.Println("Waiting...")
		<-timer
		fmt.Println("Done!!!")
		// client, _ := socketio.NewClient(uri, nil)

		// // Handle an incoming event
		// client.OnEvent("reply", func(s socketio.Conn, msg string) {
		// 	log.Println("Receive Message /reply: ", "reply", msg)
		// })

		// client.Connect()
		// client.Emit("notice", "hello")
		// client.Close()
	}()
}
