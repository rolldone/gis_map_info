package main

import (
	"github.com/gin-gonic/gin"

	AdminController "gis_map_info/app/http/admin"
	FrontController "gis_map_info/app/http/front"
)

func main() {
	// Initialize Gin's default router
	router := gin.Default()

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

			admin.GET("/rdtr/rdtrs", AdminController.GetRdtrs)
			admin.GET("/rdtr/:uuid/view", AdminController.GetRdtrByUUId)
			admin.POST("/rdtr/new", AdminController.AddRdtr)
			admin.POST("/rdtr/update", AdminController.UpdateRdtr)
			admin.POST("/rdtr/delete", AdminController.DeleteRdtr)

			admin.GET("/rtrw/rtrws", AdminController.GetRtrws)
			admin.GET("/rtrw/:uuid/view", AdminController.GetRtrwByUUId)
			admin.POST("/rtrw/new", AdminController.AddRtrw)
			admin.POST("/rtrw/update", AdminController.UpdateRtrw)
			admin.POST("/rtrw/delete", AdminController.DeleteRtrw)

			admin.GET("/zlp/zlps", AdminController.GetZLPs)
			admin.GET("/zlp/:uuid/view", AdminController.GetZLPByUUId)
			admin.POST("/zlp/new", AdminController.AddZLP)
			admin.POST("/zlp/update", AdminController.UpdateZLP)
			admin.POST("/zlp/delete", AdminController.DeleteZLP)

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
