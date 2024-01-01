package main

import (
	"gis_map_info/app"
	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
	Model "gis_map_info/app/model"

	"github.com/gin-gonic/gin"
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

			rdtrFileController := &AdminController.RdtrFileController{}
			admin.POST("/rdtr_file/add", rdtrFileController.Add)
			admin.GET("/rdtr_file/get/:uuid", rdtrFileController.GetByUUID)

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

	// Start the server
	router.Run(":8080")
}
