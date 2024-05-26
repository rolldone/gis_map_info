package main

import (
	"fmt"
	"gis_map_info/app"
	"gis_map_info/app/cli"
	"gis_map_info/app/http/admin"
	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
	"gis_map_info/app/http/middleware"
	"gis_map_info/app/model"
	"gis_map_info/app/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gis_map_info/support/app_support"
	"gis_map_info/support/asynq_support"
	"gis_map_info/support/gorm_support"
	"gis_map_info/support/log_support"
	"gis_map_info/support/nats_support"
	"gis_map_info/support/redis_support"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"

	"github.com/go-co-op/gocron/v2"
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

	// Init App Support
	app_support.Init()

	// Init log
	f := log_support.Init()
	defer f.Close()

	gorm_support.ConnectDatabase()

	app.Install()

	// Init the cli.
	// This sub app for run app as cli mode.
	bypass := cli.InitCli(gorm_support.DB)
	if !bypass {
		quit()
		return
	}
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
		adminRoute := api.Group("/admin")
		{
			adminTokenMiddleware := middleware.AdminTokenMiddleware()

			adminAuthController := admin.AuthController{}
			adminRoute.POST("/auth/login", adminAuthController.Login)
			adminRoute.POST("/auth/register", adminAuthController.Register)
			adminRoute.POST("/auth/forgot_password", adminAuthController.ForgotPassword)
			adminRoute.POST("/auth/recovery_password", adminAuthController.RecoveryPassword)
			adminRoute.GET("/auth/logout", adminTokenMiddleware, adminAuthController.Logout)
			adminRoute.GET("/auth", adminTokenMiddleware, adminAuthController.Auth)
			adminRoute.POST("/auth/profile/update", adminTokenMiddleware, adminAuthController.UpdateProfile)

			userController := admin.UserController{}
			adminRoute.POST("/user/users", userController.Gets)
			adminRoute.GET("/user/:uuid", userController.GetByUUID)
			adminRoute.POST("/user/add", userController.Add)
			adminRoute.POST("/user/update", userController.Update)
			adminRoute.POST("/user/delete", userController.Delete)

			dashboard := AdminController.Dashboard()
			adminRoute.GET("/dashboard/rdtr-periode", gin.HandlerFunc(dashboard.GetRdtrPeriode))
			adminRoute.GET("/dashboard/rtrw-periode", gin.HandlerFunc(dashboard.GetRtrwPeriode))

			rdtrController := &AdminController.RdtrController{}
			adminRoute.GET("/zone_rdtr/get/:id/view", rdtrController.GetRdtrById)
			adminRoute.GET("/zone_rdtr/gets/paginate", rdtrController.GetRdtrsPaginate)
			adminRoute.GET("/zone_rdtr/gets", rdtrController.GetRdtrs)
			adminRoute.POST("/zone_rdtr/add", rdtrController.AddRdtr)
			adminRoute.POST("/zone_rdtr/update", rdtrController.UpdateRdtr)
			adminRoute.POST("/zone_rdtr/delete", rdtrController.DeleteRdtr)
			adminRoute.POST("/zone_rdtr/validate_kml", rdtrController.ValidateKml)
			adminRoute.POST("/zone_rdtr/validate_mbtile", rdtrController.ValidateMbtile)
			adminRoute.GET("/zone_rdtr/validate/ws", rdtrController.HandleWS)

			rdtrFileController := &AdminController.RdtrFileController{}
			adminRoute.POST("/rdtr_file/add", rdtrFileController.Add)
			adminRoute.GET("/rdtr_file/get/:uuid", rdtrFileController.GetByUUID)
			adminRoute.Static("/rdtr_file/assets", "./storage/rdtr_files")

			rdtrMbtileController := &AdminController.RdtrMbtileController{}
			adminRoute.POST("/rdtr_mbtile/add", rdtrMbtileController.Add)
			adminRoute.GET("/rdtr_mbtile/get/:uuid", rdtrMbtileController.GetbyUUID)
			adminRoute.GET("/rdtr_mbtile/martin_config", rdtrMbtileController.GetMartinConfig)

			rtrwController := &AdminController.RtrwController{}
			adminRoute.GET("/zone_rtrw/get/:id/view", rtrwController.GetRtrwById)
			adminRoute.GET("/zone_rtrw/gets/paginate", rtrwController.GetRtrwsPaginate)
			adminRoute.GET("/zone_rtrw/gets", rtrwController.GetRtrws)
			adminRoute.POST("/zone_rtrw/add", rtrwController.AddRtrw)
			adminRoute.POST("/zone_rtrw/update", rtrwController.UpdateRtrw)
			adminRoute.POST("/zone_rtrw/delete", rtrwController.DeleteRtrw)
			adminRoute.POST("/zone_rtrw/validate_kml", rtrwController.ValidateKml)
			adminRoute.POST("/zone_rtrw/validate_mbtile", rtrwController.ValidateMbtile)
			adminRoute.GET("/zone_rtrw/validate/ws", rtrwController.HandleWS)

			rtrwFileController := &AdminController.RtrwFileController{}
			adminRoute.POST("/rtrw_file/add", rtrwFileController.Add)
			adminRoute.GET("/rtrw_file/get/:uuid", rtrwFileController.GetByUUID)
			adminRoute.Static("/rtrw_file/assets", "./storage/rtrw_files")

			rtrwMbtileController := &AdminController.RtrwMbtileController{}
			adminRoute.POST("/rtrw_mbtile/add", rtrwMbtileController.Add)
			adminRoute.GET("/rtrw_mbtile/get/:uuid", rtrwMbtileController.GetbyUUID)
			adminRoute.GET("/rtrw_mbtile/martin_config", rtrwMbtileController.GetMartinConfig)

			zlpController := &AdminController.ZlpController{}
			adminRoute.GET("/zone_zlp/get/:id/view", zlpController.GetZlpById)
			adminRoute.GET("/zone_zlp/gets/paginate", zlpController.GetZlpsPaginate)
			adminRoute.GET("/zone_zlp/gets", zlpController.GetZlps)
			adminRoute.POST("/zone_zlp/add", zlpController.AddZlp)
			adminRoute.POST("/zone_zlp/update", zlpController.UpdateZlp)
			adminRoute.POST("/zone_zlp/delete", zlpController.DeleteZlp)
			adminRoute.POST("/zone_zlp/validate_kml", zlpController.ValidateKml)
			adminRoute.POST("/zone_zlp/validate_mbtile", zlpController.ValidateMbtile)
			adminRoute.GET("/zone_zlp/validate/ws", zlpController.HandleWS)

			zlpFileController := &AdminController.ZlpFileController{}
			adminRoute.POST("/zlp_file/add", zlpFileController.Add)
			adminRoute.GET("/zlp_file/get/:uuid", zlpFileController.GetByUUID)
			adminRoute.Static("/zlp_file/assets", "./storage/zlp_files")

			zlpMbtileController := &AdminController.ZlpMbtileController{}
			adminRoute.POST("/zlp_mbtile/add", zlpMbtileController.Add)
			adminRoute.GET("/zlp_mbtile/get/:uuid", zlpMbtileController.GetbyUUID)
			adminRoute.GET("/zlp_mbtile/martin_config", zlpMbtileController.GetMartinConfig)

			regLocationController := &AdminController.RegLocationController{}
			adminRoute.GET("/reg_location/province/provinces", regLocationController.GetProvinces)
			adminRoute.GET("/reg_location/regency/regencies", regLocationController.GetRegencies)
			adminRoute.GET("/reg_location/district/districts", regLocationController.GetDistricts)
			adminRoute.GET("/reg_location/village/villages", regLocationController.GetVillages)

			asynqJobController := AdminController.AsynqJobController{}
			adminRoute.GET("/asynq_job/asynq_jobs", asynqJobController.GetsAsynqJob)
			adminRoute.GET("/asynq_job/:uuid/app_uuid", asynqJobController.GetAsynqJobByAppUuid)
			adminRoute.POST("/asynq_job/delete", asynqJobController.DeleteAsynqJobByUUIDS)

			// Front Job log Controller
			frontJobLogController := AdminController.MessageLogControllerConstruct()
			adminRoute.POST("/message_log/message_logs", frontJobLogController.Gets)

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

	go func() {
		// Start schedule delete data
		InitDeleteUnusedData_WithSchedule()
		ValidateMbtileOnStart()
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

func ValidateMbtileOnStart() {
	fmt.Println("Validate mbtile on start")
	func() {
		var mbtile_datas = []model.RdtrMbtile{}
		// Load RdtrMbtileService
		rdtrMbtileService := service.RdtrMbtileService{
			DB: gorm_support.DB,
		}

		// Next get all check datas
		rdtrMbtileDB := rdtrMbtileService.Gets()
		err := rdtrMbtileDB.Where("checked_at IS NOT NULL").Find(&mbtile_datas).Error
		if err != nil {
			fmt.Println("err - 2342344 :: ", err.Error())
			return
		}

		martinConfig, err := os.ReadFile("./sub_app/martin/config.yaml")
		if err != nil {
			fmt.Println("err - 2562344 :: ", err.Error())
			return
		}

		var martinMap map[string]interface{}
		err = yaml.Unmarshal(martinConfig, &martinMap)
		if err != nil {
			fmt.Println("err - 2554344 :: ", err.Error())
			return
		}

		martin_mbtile_sources := map[string]interface{}{}
		martin_mbtiles := martinMap["mbtiles"].(map[string]interface{})
		martin_mbtile_sources_parse, ok := martin_mbtiles["sources"].(map[string]interface{})
		if ok {
			martin_mbtile_sources = martin_mbtile_sources_parse
		}

		// Iterate over the map
		for _, value2 := range mbtile_datas {
			martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/rdtr/", value2.UUID, ".mbtiles")
		}

		// Redefine again to martin config mbtiles
		martin_mbtiles["sources"] = martin_mbtile_sources
		martinMap["mbtiles"] = martin_mbtiles

		// The last save config martin
		martinConfig, err = yaml.Marshal(martinMap)
		if err != nil {
			fmt.Println("err - 2938854734 :: ", err.Error())
			return
		}

		err = os.WriteFile("./sub_app/martin/config.yaml", []byte(martinConfig), 0755)
		if err != nil {
			fmt.Println("err - 2938855534 :: ", err.Error())
			return
		}
	}()

	func() {
		var mbtile_datas = []model.RtrwMbtile{}
		// Load RdtrMbtileService
		rtrwMbtileService := service.RtrwMbtileService{
			DB: gorm_support.DB,
		}

		// Next get all check datas
		rtrwMbtileDB := rtrwMbtileService.Gets()
		err := rtrwMbtileDB.Where("checked_at IS NOT NULL").Find(&mbtile_datas).Error
		if err != nil {
			fmt.Println("err - 2342344 :: ", err.Error())
			return
		}

		martinConfig, err := os.ReadFile("./sub_app/martin/config.yaml")
		if err != nil {
			fmt.Println("err - 2562344 :: ", err.Error())
			return
		}

		var martinMap map[string]interface{}
		err = yaml.Unmarshal(martinConfig, &martinMap)
		if err != nil {
			fmt.Println("err - 2554344 :: ", err.Error())
			return
		}

		martin_mbtile_sources := map[string]interface{}{}
		martin_mbtiles := martinMap["mbtiles"].(map[string]interface{})
		martin_mbtile_sources_parse, ok := martin_mbtiles["sources"].(map[string]interface{})
		if ok {
			martin_mbtile_sources = martin_mbtile_sources_parse
		}

		// Iterate over the map
		for _, value2 := range mbtile_datas {
			martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/rtrw/", value2.UUID, ".mbtiles")
		}

		// Redefine again to martin config mbtiles
		martin_mbtiles["sources"] = martin_mbtile_sources
		martinMap["mbtiles"] = martin_mbtiles

		// The last save config martin
		martinConfig, err = yaml.Marshal(martinMap)
		if err != nil {
			fmt.Println("err - 2938854734 :: ", err.Error())
			return
		}

		err = os.WriteFile("./sub_app/martin/config.yaml", []byte(martinConfig), 0755)
		if err != nil {
			fmt.Println("err - 2938855534 :: ", err.Error())
			return
		}
	}()
	func() {
		var mbtile_datas = []model.ZlpMbtile{}
		// Load RdtrMbtileService
		zlpMbtileService := service.ZlpMbtileService{
			DB: gorm_support.DB,
		}

		// Next get all check datas
		zlpMbtileDB := zlpMbtileService.Gets()
		err := zlpMbtileDB.Where("checked_at IS NOT NULL").Find(&mbtile_datas).Error
		if err != nil {
			fmt.Println("err - 2342344 :: ", err.Error())
			return
		}

		martinConfig, err := os.ReadFile("./sub_app/martin/config.yaml")
		if err != nil {
			fmt.Println("err - 2562344 :: ", err.Error())
			return
		}

		var martinMap map[string]interface{}
		err = yaml.Unmarshal(martinConfig, &martinMap)
		if err != nil {
			fmt.Println("err - 2554344 :: ", err.Error())
			return
		}

		martin_mbtile_sources := map[string]interface{}{}
		martin_mbtiles := martinMap["mbtiles"].(map[string]interface{})
		martin_mbtile_sources_parse, ok := martin_mbtiles["sources"].(map[string]interface{})
		if ok {
			martin_mbtile_sources = martin_mbtile_sources_parse
		}

		// Iterate over the map
		for _, value2 := range mbtile_datas {
			martin_mbtile_sources[value2.UUID] = fmt.Sprint("/app/mbtiles/zlp/", value2.UUID, ".mbtiles")
		}

		// Redefine again to martin config mbtiles
		martin_mbtiles["sources"] = martin_mbtile_sources
		martinMap["mbtiles"] = martin_mbtiles

		// The last save config martin
		martinConfig, err = yaml.Marshal(martinMap)
		if err != nil {
			fmt.Println("err - 2938854734 :: ", err.Error())
			return
		}

		err = os.WriteFile("./sub_app/martin/config.yaml", []byte(martinConfig), 0755)
		if err != nil {
			fmt.Println("err - 2938855534 :: ", err.Error())
			return
		}
	}()
}

func InitDeleteUnusedData_WithSchedule() {

	// Calculate the date from three months ago
	monthAgo := time.Now().AddDate(0, -1, 0)

	// add a job to the scheduler
	location, _ := time.LoadLocation(app_support.App.Time_zone)

	// create a scheduler
	s, err := gocron.NewScheduler(
		gocron.WithLocation(location),
	)

	// Every day at 12:00 AM
	// contJobVal := fmt.Sprint("0 0 * * *")

	// Every 1 minuted
	contJobVal := fmt.Sprint("* * * * *")

	j, err := s.NewJob(
		gocron.CronJob(contJobVal, false),
		gocron.NewTask(
			func(props map[string]interface{}) {

				// Select data from zlp
				// Select mbtile from  zlp
				func() {
					zlpMbtileService := service.ZlpMbtileService{
						DB: gorm_support.DB,
					}
					zlpMbtileDB := zlpMbtileService.Gets()
					zlpMbtileDatas := []model.ZlpMbtile{}
					err := zlpMbtileDB.Where("zlp_id is NULL").Where("created_at < ?", monthAgo).Find(&zlpMbtileDatas).Error
					if err != nil {
						fmt.Println("Err 293874755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data zlp mbtiles :: ", len(zlpMbtileDatas))
					pathMbtile := app_support.App.Mbtile_zlp_path
					if _, err := os.Stat(pathMbtile); os.IsNotExist(err) {
						for a := 0; a < len(zlpMbtileDatas); a++ {
							err := os.Remove(pathMbtile + "/" + zlpMbtileDatas[a].UUID + ".mbtiles")
							if err != nil {
								log.Println(pathMbtile+"/"+zlpMbtileDatas[a].UUID, " not exist")
							}
						}
					}
				}()
				// Select files from  zlp
				func() {
					zlpFilesService := service.ZlpFileService{
						DB: gorm_support.DB,
					}
					zlpFileDB := zlpFilesService.Gets(map[string]interface{}{})
					zlpFileDatas := []model.ZlpFile{}
					err = zlpFileDB.Where("zlp_id is NULL").Where("created_at < ?", monthAgo).Find(&zlpFileDatas).Error
					if err != nil {
						fmt.Println("Err 293474755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data zlp files :: ", len(zlpFileDatas))
					pathFile := app_support.App.File_zlp_path
					if _, err := os.Stat(pathFile); os.IsNotExist(err) {
						for a := 0; a < len(zlpFileDatas); a++ {
							err := os.Remove(pathFile + "/" + zlpFileDatas[a].UUID + ".kml")
							if err != nil {
								log.Println(pathFile+"/"+zlpFileDatas[a].UUID, " not exist")
							}
							err = os.Remove(pathFile + "/" + zlpFileDatas[a].UUID + ".json")
							if err != nil {
								log.Println(pathFile+"/"+zlpFileDatas[a].UUID, " not exist")
							}
						}
					}
				}()

				// Select data from rdtr
				// Select mbtile from  rdtr
				func() {
					rdtrMbtileService := service.RdtrMbtileService{
						DB: gorm_support.DB,
					}
					rdtrMbtileDB := rdtrMbtileService.Gets()
					rdtrMbtileDatas := []model.RdtrMbtile{}
					err = rdtrMbtileDB.Where("rdtr_id is NULL").Where("created_at < ?", monthAgo).Find(&rdtrMbtileDatas).Error
					if err != nil {
						fmt.Println("Err 243874755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data Rdtr mbtiles :: ", len(rdtrMbtileDatas))
					pathMbtile := app_support.App.Mbtile_rdtr_path
					if _, err := os.Stat(pathMbtile); os.IsNotExist(err) {
						for a := 0; a < len(rdtrMbtileDatas); a++ {
							err := os.Remove(pathMbtile + "/" + rdtrMbtileDatas[a].UUID)
							if err != nil {
								log.Println(pathMbtile+"/"+rdtrMbtileDatas[a].UUID, " not exist")
							}
						}
					}
				}()
				// Select files from rtrw
				func() {
					rdtrFilesService := service.RdtrFileService{
						DB: gorm_support.DB,
					}
					rdtrFileDB := rdtrFilesService.Gets(map[string]interface{}{})
					rdtrFileDatas := []model.RdtrFile{}
					err = rdtrFileDB.Where("rdtr_id is NULL").Where("created_at < ?", monthAgo).Find(&rdtrFileDatas).Error
					if err != nil {
						fmt.Println("Err 293474755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data Rtrw files :: ", len(rdtrFileDatas))
					pathFile := app_support.App.File_rdtr_path
					if _, err := os.Stat(pathFile); os.IsNotExist(err) {
						for a := 0; a < len(rdtrFileDatas); a++ {
							err := os.Remove(pathFile + "/" + rdtrFileDatas[a].UUID + ".kml")
							if err != nil {
								log.Println(pathFile+"/"+rdtrFileDatas[a].UUID, " not exist")
							}
							err = os.Remove(pathFile + "/" + rdtrFileDatas[a].UUID + ".json")
							if err != nil {
								log.Println(pathFile+"/"+rdtrFileDatas[a].UUID, " not exist")
							}
						}
					}
				}()

				// Select data from rtrw
				// Select mbtile from  rtrw
				func() {
					rtrwMbtileService := service.RtrwMbtileService{
						DB: gorm_support.DB,
					}
					rtrwMbtileDB := rtrwMbtileService.Gets()
					rtrwMbtileDatas := []model.RtrwMbtile{}
					err = rtrwMbtileDB.Where("rtrw_id is NULL").Where("created_at < ?", monthAgo).Find(&rtrwMbtileDatas).Error
					if err != nil {
						fmt.Println("Err 243854755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data Rtrw mbtiles :: ", len(rtrwMbtileDatas))
					pathMbtile := app_support.App.Mbtile_rtrw_path
					if _, err := os.Stat(pathMbtile); os.IsNotExist(err) {
						for a := 0; a < len(rtrwMbtileDatas); a++ {
							err := os.Remove(pathMbtile + "/" + rtrwMbtileDatas[a].UUID)
							if err != nil {
								log.Println(pathMbtile+"/"+rtrwMbtileDatas[a].UUID, " not exist")
							}
						}
					}
				}()
				// Select files from rtrw
				func() {
					rtrwFilesService := service.RtrwFileService{
						DB: gorm_support.DB,
					}
					rtrwFileDB := rtrwFilesService.Gets(map[string]interface{}{})
					rtrwFileDatas := []model.RtrwFile{}
					err = rtrwFileDB.Where("rtrw_id is NULL").Where("created_at < ?", monthAgo).Find(&rtrwFileDatas).Error
					if err != nil {
						fmt.Println("Err 293474755 :: ", err)
						panic(1)
					}
					fmt.Println("Total data Rtrw files  :: ", len(rtrwFileDatas))
					pathFile := app_support.App.File_rtrw_path
					if _, err := os.Stat(pathFile); os.IsNotExist(err) {
						for a := 0; a < len(rtrwFileDatas); a++ {
							err := os.Remove(pathFile + "/" + rtrwFileDatas[a].UUID + ".kml")
							if err != nil {
								log.Println(pathFile+"/"+rtrwFileDatas[a].UUID, " not exist")
							}
							err = os.Remove(pathFile + "/" + rtrwFileDatas[a].UUID + ".json")
							if err != nil {
								log.Println(pathFile+"/"+rtrwFileDatas[a].UUID, " not exist")
							}
						}
					}
				}()

			},
			map[string]interface{}{},
		),
	)

	if err != nil {
		// handle error
		log.Println("handleClearJobData - err - 2343223244 :: ", err)
		return
	}

	fmt.Println("Crontjob task with id :: ", j.ID(), " Started")

	s.Start()
}

func quit() {
	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Shutdown Gin server (Gin does not provide a built-in shutdown method, so you may handle this as needed)
	// For example, you can use http.Server for more control over the Gin server and shutdown it gracefully.
	log.Println("Servers gracefully stopped.")
}
