package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"strconv"
)

func main() {

	cfg := NewConfig()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api")

	db, err := NewMongoDb(cfg.DbHost)
	if err != nil {
		log.Fatal(err)
	}

	repo := NewRepository(db, cfg.DbName)

	apiHandlers := NewAPI(repo)

	api.GET("/bots", apiHandlers.getBots)
	api.GET("/tournaments", apiHandlers.getTournaments)
	api.POST("/teams", apiHandlers.createTeam)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}
