package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang-echo-elasticsearch-rest-api-example/config"
	"golang-echo-elasticsearch-rest-api-example/controller"
	_ "golang-echo-elasticsearch-rest-api-example/docs"
	"golang-echo-elasticsearch-rest-api-example/handler"
	"golang-echo-elasticsearch-rest-api-example/repository"
	"golang-echo-elasticsearch-rest-api-example/routes"
	"golang-echo-elasticsearch-rest-api-example/util"
	"log"
)

var userController *controller.UserController

// @title Golang User REST API
// @description Provides access to the core features of Golang User REST API
// @version 1.0
// @termsOfService http://swagger.io/terms/
// license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func main() {
	e := echo.New()

	e.HTTPErrorHandler = handler.ErrorHandler
	e.Validator = util.NewValidationUtil()
	config.CORSConfig(e)

	routes.GetUserApiRoutes(e, userController)
	routes.GetSwaggerRoutes(e)

	// echo server 9000 de başlatıldı.
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.ServerPort)))
}

func init() {
	es, err := config.ElasticsearchConnection()
	if err != nil {
		log.Println("Error when connect elasticsearch : ", err.Error())
	}

	userRepository := repository.NewUserRepository(es)
	userController = controller.NewUserController(userRepository)
}
