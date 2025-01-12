package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"slashbase.com/backend/src/config"
	"slashbase.com/backend/src/controllers"
	"slashbase.com/backend/src/middlewares"
)

// NewRouter return a gin router for server
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowCredentials = true
	if config.IsLive() || config.IsDevelopment() {
		corsConfig.AllowOrigins = []string{config.GetAppHost()}
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	}
	router.Use(cors.New(corsConfig))
	api := router.Group("/api/v1")
	{
		userGroup := api.Group("user")
		{
			userController := new(controllers.UserController)
			userGroup.POST("/login", userController.LoginUser)
			userGroup.Use(middlewares.FindUserMiddleware())
			userGroup.Use(middlewares.AuthUserMiddleware())
			userGroup.POST("/edit", userController.EditAccount)
			userGroup.POST("/add", userController.AddUser)
			userGroup.GET("/all", userController.GetUsers)
			userGroup.GET("/logout", userController.Logout)
		}
		projectGroup := api.Group("project")
		{
			projectController := new(controllers.ProjectController)
			projectGroup.Use(middlewares.FindUserMiddleware())
			projectGroup.Use(middlewares.AuthUserMiddleware())
			projectGroup.POST("/create", projectController.CreateProject)
			projectGroup.GET("/all", projectController.GetProjects)
			projectGroup.POST("/:projectId/members/create", projectController.AddProjectMembers)
			projectGroup.GET("/:projectId/members", projectController.GetProjectMembers)
		}
		dbConnGroup := api.Group("dbconnection")
		{
			dbConnController := new(controllers.DBConnectionController)
			dbConnGroup.Use(middlewares.FindUserMiddleware())
			dbConnGroup.Use(middlewares.AuthUserMiddleware())
			dbConnGroup.POST("/create", dbConnController.CreateDBConnection)
			dbConnGroup.GET("/all", dbConnController.GetDBConnections)
			dbConnGroup.GET("/project/:projectId", dbConnController.GetDBConnectionsByProject)
			dbConnGroup.GET("/:dbConnId", dbConnController.GetSingleDBConnection)
		}
		queryGroup := api.Group("query")
		{
			queryController := new(controllers.QueryController)
			queryGroup.Use(middlewares.FindUserMiddleware())
			queryGroup.Use(middlewares.AuthUserMiddleware())
			queryGroup.POST("/run", queryController.RunQuery)
			queryGroup.POST("/save/:dbConnId", queryController.SaveDBQuery)
			queryGroup.GET("/getall/:dbConnId", queryController.GetDBQueriesInDBConnection)
			queryGroup.GET("/get/:queryId", queryController.GetSingleDBQuery)
			queryGroup.GET("/history/:dbConnId", queryController.GetQueryHistoryInDBConnection)
			dataGroup := queryGroup.Group("data")
			{
				dataGroup.GET("/:dbConnId", queryController.GetData)
				dataGroup.POST("/:dbConnId/single", queryController.UpdateSingleData)
				dataGroup.POST("/:dbConnId/add", queryController.AddData)
				dataGroup.POST("/:dbConnId/delete", queryController.DeleteData)
			}
			dataModelGroup := queryGroup.Group("datamodel")
			{
				dataModelGroup.GET("/all/:dbConnId", queryController.GetDataModels)
				dataModelGroup.GET("/single/:dbConnId", queryController.GetSingleDataModel)
			}
		}
	}
	return router

}
