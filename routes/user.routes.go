package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-gorm-postgres/controllers"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (rc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("user")

	router.GET("/check", rc.userController.GetUsers)
	router.POST("/create", rc.userController.CreateUser)
	router.GET("/get", rc.userController.GetUserByID)
	router.GET("/records", rc.userController.GetRecords)
}
