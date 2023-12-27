package main

import (
	"github.com/gin-gonic/gin"

	"gis_map_info/app"
	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
)

func main() {

	app.Install()

	return
	// Initialize Gin's default router
	router := gin.Default()

	router.Use(gin.Logger())

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
			admin.GET("/rdtr/rdtrs", rdtrController.GetRdtrs)
			admin.GET("/rdtr/:uuid/view", rdtrController.GetRdtrByUUId)
			admin.POST("/rdtr/new", rdtrController.AddRdtr)
			admin.POST("/rdtr/update", rdtrController.UpdateRdtr)
			admin.POST("/rdtr/delete", rdtrController.DeleteRdtr)

			rtrwController := &AdminController.RtrwController{}
			admin.GET("/rtrw/rtrws", rtrwController.GetRtrws)
			admin.GET("/rtrw/:uuid/view", rtrwController.GetRtrwByUUId)
			admin.POST("/rtrw/new", rtrwController.AddRtrw)
			admin.POST("/rtrw/update", rtrwController.UpdateRtrw)
			admin.POST("/rtrw/delete", rtrwController.DeleteRtrw)

			zlpController := &AdminController.ZLPController{}
			admin.GET("/zlp/zlps", zlpController.GetZLPs)
			admin.GET("/zlp/:uuid/view", zlpController.GetZLPByUUId)
			admin.POST("/zlp/new", zlpController.AddZLP)
			admin.POST("/zlp/update", zlpController.UpdateZLP)
			admin.POST("/zlp/delete", zlpController.DeleteZLP)

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
