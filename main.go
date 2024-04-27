package main

import (
	"fmt"
	"gis_map_info/app"
	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
	"gis_map_info/app/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gis_map_info/support/asynq_support"
	"gis_map_info/support/gorm_support"
	"gis_map_info/support/log_support"
	"gis_map_info/support/nats_support"
	"gis_map_info/support/redis_support"

	"github.com/gin-gonic/gin"
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

	// Init log
	f := log_support.Init()
	defer f.Close()

	gorm_support.ConnectDatabase()

	app.Install()

	// Init nats pub sub
	InitNats()

	// Set stopped the last bad asynq on database
	asynqJobService := service.AsynqJobService{
		DB: gorm_support.DB,
	}
	asynqJobService.StopLastProcess()

	// Define asynq client
	client := asynq_support.Init()
	defer client.Close()

	// Define Job Manager
	// jobm_support.InitJobManager()

	// Define redis client
	redisClient := redis_support.ConnectRedis()
	defer redisClient.Close()

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
			admin.POST("/zone_rdtr/validate_kml", rdtrController.ValidateKml)
			admin.POST("/zone_rdtr/validate_mbtile", rdtrController.ValidateMbtile)
			admin.GET("/zone_rdtr/validate/ws", rdtrController.HandleWS)

			rdtrFileController := &AdminController.RdtrFileController{}
			admin.POST("/rdtr_file/add", rdtrFileController.Add)
			admin.GET("/rdtr_file/get/:uuid", rdtrFileController.GetByUUID)
			admin.Static("/rdtr_file/assets", "./storage/rdtr_files")

			rdtrMbtileController := &AdminController.RdtrMbtileController{}
			admin.POST("/rdtr_mbtile/add", rdtrMbtileController.Add)
			admin.GET("/rdtr_mbtile/get/:uuid", rdtrMbtileController.GetbyUUID)
			admin.GET("/rdtr_mbtile/martin_config", rdtrMbtileController.GetMartinConfig)

			rtrwController := &AdminController.RtrwController{}
			admin.GET("/zone_rtrw/get/:id/view", rtrwController.GetRtrwById)
			admin.GET("/zone_rtrw/gets/paginate", rtrwController.GetRtrwsPaginate)
			admin.GET("/zone_rtrw/gets", rtrwController.GetRtrws)
			admin.POST("/zone_rtrw/add", rtrwController.AddRtrw)
			admin.POST("/zone_rtrw/update", rtrwController.UpdateRtrw)
			admin.POST("/zone_rtrw/delete", rtrwController.DeleteRtrw)
			admin.POST("/zone_rtrw/validate_kml", rtrwController.ValidateKml)
			admin.POST("/zone_rtrw/validate_mbtile", rtrwController.ValidateMbtile)
			admin.GET("/zone_rtrw/validate/ws", rtrwController.HandleWS)

			rtrwFileController := &AdminController.RtrwFileController{}
			admin.POST("/rtrw_file/add", rtrwFileController.Add)
			admin.GET("/rtrw_file/get/:uuid", rtrwFileController.GetByUUID)
			admin.Static("/rtrw_file/assets", "./storage/rtrw_files")

			rtrwMbtileController := &AdminController.RtrwMbtileController{}
			admin.POST("/rtrw_mbtile/add", rtrwMbtileController.Add)
			admin.GET("/rtrw_mbtile/get/:uuid", rtrwMbtileController.GetbyUUID)
			admin.GET("/rtrw_mbtile/martin_config", rtrwMbtileController.GetMartinConfig)

			zlpController := &AdminController.ZlpController{}
			admin.GET("/zone_zlp/get/:id/view", zlpController.GetZlpById)
			admin.GET("/zone_zlp/gets/paginate", zlpController.GetZlpsPaginate)
			admin.GET("/zone_zlp/gets", zlpController.GetZlps)
			admin.POST("/zone_zlp/add", zlpController.AddZlp)
			admin.POST("/zone_zlp/update", zlpController.UpdateZlp)
			admin.POST("/zone_zlp/delete", zlpController.DeleteZlp)
			admin.POST("/zone_zlp/validate_kml", zlpController.ValidateKml)
			admin.POST("/zone_zlp/validate_mbtile", zlpController.ValidateMbtile)
			admin.GET("/zone_zlp/validate/ws", zlpController.HandleWS)

			zlpFileController := &AdminController.ZlpFileController{}
			admin.POST("/zlp_file/add", zlpFileController.Add)
			admin.GET("/zlp_file/get/:uuid", zlpFileController.GetByUUID)
			admin.Static("/zlp_file/assets", "./storage/zlp_files")

			zlpMbtileController := &AdminController.ZlpMbtileController{}
			admin.POST("/zlp_mbtile/add", zlpMbtileController.Add)
			admin.GET("/zlp_mbtile/get/:uuid", zlpMbtileController.GetbyUUID)
			admin.GET("/zlp_mbtile/martin_config", zlpMbtileController.GetMartinConfig)

			regLocationController := &AdminController.RegLocationController{}
			admin.GET("/reg_location/province/provinces", regLocationController.GetProvinces)
			admin.GET("/reg_location/regency/regencies", regLocationController.GetRegencies)
			admin.GET("/reg_location/district/districts", regLocationController.GetDistricts)
			admin.GET("/reg_location/village/villages", regLocationController.GetVillages)

			asynqJobController := AdminController.AsynqJobController{}
			admin.GET("/asynq_job/asynq_jobs", asynqJobController.GetsAsynqJob)
			admin.GET("/asynq_job/:uuid/app_uuid", asynqJobController.GetAsynqJobByAppUuid)
			admin.POST("/asynq_job/delete", asynqJobController.DeleteAsynqJobByUUIDS)

			// Front Job log Controller
			frontJobLogController := AdminController.MessageLogControllerConstruct()
			admin.POST("/message_log/message_logs", frontJobLogController.Gets)

		}

		// Front side

		frontRdtrController := FrontController.RdtrController()
		api.GET("rdtr/rdtrs", frontRdtrController.Gets)
		api.GET("rdtr/:id/view", frontRdtrController.GetByUUID)
		api.GET("rdtr/position/:latlng", frontRdtrController.GetByPosition)
		api.GET("rdtr/regencies/:province_id", frontRdtrController.GetRegenciesByProvinceId)

		frontRtrwController := FrontController.RtrwController()
		api.GET("rtrw/rtrws", frontRtrwController.Gets)
		api.GET("rtrw/:id/view", frontRtrwController.GetByUUID)
		api.GET("rtrw/position/:latlng", frontRtrwController.GetByPosition)
		api.GET("rtrw/regencies/:province_id", frontRtrwController.GetRegenciesByProvinceId)

		frontZlpController := FrontController.ZlpController()
		api.GET("zlp/zlps", frontZlpController.Gets)
		api.GET("zlp/:id/view", frontZlpController.GetByUUID)
		api.GET("zlp/position/:latlng", frontZlpController.GetByPosition)
		api.GET("zlp/regencies/:province_id", frontZlpController.GetRegenciesByProvinceId)
		api.GET("zlp_group/zlp_groups", frontZlpController.GetsByZlpGroup)
		api.GET("zlp_group/position/:latlng", frontZlpController.GetPositionByZlpGroup)

		locationController := FrontController.LocationController{}
		api.GET("location/provinces/exists", locationController.GetsProvinceDistincExist)

		api.GET("test/sse", FrontController.Testsse)

	}

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

func InitNats() {
	_, err := nats_support.ConnectPubSub()
	if err != nil {
		fmt.Println("Nats error :: ", err.Error())
	} else {
		// Simple Async Subscriber
		nats_support.NATS.Subscribe("foo", func(m *nats.Msg) {
			fmt.Printf("\nReceived a message: %s\n", string(m.Data))
		})
		go func() {
			timer := time.After(5 * time.Second)
			<-timer
			// Simple Publisher
			nats_support.NATS.Publish("foo", []byte("Hello World NATS"))
		}()
	}

}
