package app

import (
	"fmt"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/controllers"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/middleware"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
)


type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Server	string
}

func (a *App) Initialize(projectEnv string) {

	e := c.LoadDotenv(projectEnv)
	if e != nil {
		fmt.Println(e)
	}

	models.InitDB(projectEnv)
	a.DB = models.GetDB()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8000"
	}
	serverAddress := os.Getenv("SERVER_ADDRESS")
	a.Server = serverAddress + ":" + serverPort

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	jwtAuth :=  os.Getenv("JWT_AUTH")
	if jwtAuth == "enabled" {
		a.Router.Use(middleware.JwtAuthentication)
	}

}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/user/register", controllers.CreateUser).Methods("POST")
	a.Router.HandleFunc("/api/user/login", controllers.AuthenticateUser).Methods("POST")
	a.Router.HandleFunc("/api/token/refresh", controllers.RefreshToken).Methods("POST")
	a.Router.HandleFunc("/api/assets/favorites", controllers.GetAssetsFavor).Methods("GET")
	a.Router.HandleFunc("/api/assets/favorites", controllers.SetAssetsFavor).Methods("POST")
	a.Router.HandleFunc("/api/assets/favorites", controllers.UnsetAssetsFavor).Methods("DELETE")
	a.Router.HandleFunc("/api/assets", controllers.GetAssets).Methods("GET")
	a.Router.HandleFunc("/api/assets/{id}", controllers.GetAsset).Methods("GET")
	a.Router.HandleFunc("/api/assets/{id}", controllers.UpdateAsset).Methods("PUT")
	a.Router.HandleFunc("/api/assets/{id}", controllers.DeleteAsset).Methods("DELETE")
	a.Router.HandleFunc("/populate/assets", controllers.PopulateAssets).Methods("POST")
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.Server, a.Router))
}
