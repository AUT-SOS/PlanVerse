package routes

import (
	"PlanVerse/controllers"
	"PlanVerse/helpers"
	"PlanVerse/middlewares"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func Routes(server *echo.Echo, customRegistry *prometheus.Registry) {
	//middlewares
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	server.Use(middlewares.AuthMiddleware)
	server.Use(middlewares.VerifyMiddleware)

	//user api //means done
	server.POST("/register", controllers.RegisterHandler) //
	server.POST("/verify", controllers.VerifyHandler)     //
	server.POST("/refresh", controllers.RefreshHandler)
	server.POST("/login", controllers.LoginHandler)               //
	server.POST("/resend-email", controllers.ResendEmailHandler)  //
	server.GET("/get-user/:user-id", controllers.GetUserHandler)  //
	server.GET("/get-my-user", controllers.GetUserIDHandler)      //
	server.POST("/edit-profile", controllers.EditUserHandler)     //
	server.POST("/delete-account", controllers.DeleteUserHandler) //

	//project api
	server.GET("/list-project", controllers.ProjectListHandler)                                                    //
	server.POST("/create-project", controllers.CreateProjectHandler)                                               //
	server.POST("/share-link/:project-id", controllers.ShareProjectHandler)                                        //
	server.POST("/show-project", controllers.ShowProjectHandler)                                                   //
	server.POST("/join-project/:project-id", controllers.JoinProjectHandler)                                       //
	server.POST("/promote/:project-id/:user-id", controllers.ChangeRoleMemberHandler, middlewares.AdminMiddleware) //
	server.POST("/demote/:project-id/:user-id", controllers.ChangeRoleAdminHandler)                                //
	server.GET("/get-project/:project-id", controllers.GetProjectHandler)                                          //
	server.GET("/get-project-members/:project-id", controllers.GetProjectMembersHandler)                           //
	server.POST("/leave-project/:project-id", controllers.LeaveProjectHandler)                                     //
	server.POST("/edit-project/:project-id", controllers.EditProjectHandler, middlewares.AdminMiddleware)          //
	server.POST("/delete-project/:project-id", controllers.DeleteProjectHandler)                                   //

	//state api
	server.GET("/list-state/:project-id", controllers.StateListHandler)                                             //
	server.POST("/create-state/:project-id", controllers.CreateStateHandler, middlewares.AdminMiddleware)           //
	server.POST("/edit-state/:project-id/:state-id", controllers.EditStateHandler, middlewares.AdminMiddleware)     //removed ID from request body. instead added state-id in query-param
	server.POST("/delete-state/:project-id/:state-id", controllers.DeleteStateHandler, middlewares.AdminMiddleware) //
	server.GET("/get-state/:project-id/:state-id", controllers.GetStateHandler)                                     //

	//task api
	server.POST("/create-task/:project-id/:state-id", controllers.CreateTaskHandler, middlewares.AdminMiddleware)          //removed stateID from request body. instead added state-id in query-param
	server.POST("/change-state/:project-id/:task-id", controllers.ChangeTaskStateHandler)                                  //removed taskID from request body. instead added task-id in query-param
	server.POST("/add-performer/:project-id/:task-id", controllers.AddPerformerHandler, middlewares.AdminMiddleware)       //removed taskID from request body. instead added task-id in query-param
	server.POST("/remove-performer/:project-id/:task-id", controllers.RemovePerformerHandler, middlewares.AdminMiddleware) //removed taskID from request body. instead added task-id in query-param
	server.POST("/edit-task/:project-id/:task-id", controllers.EditTaskHandler, middlewares.AdminMiddleware)               //removed taskID from request body. instead added task-id in query-param
	server.POST("/delete-task/:project-id/:task-id", controllers.DeleteTaskHandler, middlewares.AdminMiddleware)           //
	server.GET("/get-task/:task-id", controllers.GetTaskHandler)                                                           //

	//create metrics endpoint for prometheus
	server.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: customRegistry}))

	//create websocket connection
	server.GET("/ws", controllers.CreateWSConnection)
	go helpers.CheckAliveness()

	//start server
	log.Fatal(server.Start("0.0.0.0:8080"))
}
