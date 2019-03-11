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

	e.Validator = NewCustomValidator()

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db, err := NewMongoDb(cfg.DbHost)
	if err != nil {
		log.Fatal(err)
	}

	apiHandlers := NewAPI(NewRepository(db, cfg.DbName))

	jwtMiddleware := middleware.JWT(cfg.JWTSecret)

	api := e.Group("/api")

	api.POST("/login", apiHandlers.login)
	api.POST("/register", apiHandlers.register)

	api.GET("/bots", apiHandlers.getBots, jwtMiddleware)
	api.POST("/bots", apiHandlers.addBot, jwtMiddleware)
	api.DELETE("/bots/:id", apiHandlers.deleteBot, jwtMiddleware)
	api.GET("/tournaments", apiHandlers.getTournaments)
	api.POST("/teams", apiHandlers.createTeam, jwtMiddleware)

	ws := NewWebSocket()

	//ws.OnConnection = func(socket *Socket) {
	//	socket.On("ping", func(payload interface{}) {
	//		socket.Emit("ping_result", "pong")
	//	})
	//}

	e.GET("/io", ws.Handler)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}
