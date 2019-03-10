package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type API struct {
	repo *Repository
}

func (api *API) getBots(c echo.Context) error {
	bots, err := api.repo.FindBots()
	if err != nil {
		return err
	}
	return c.JSON(200, bots)
}

func (api *API) getTournaments(c echo.Context) error {
	tournaments, err := api.repo.FindTournaments()
	if err != nil {
		return err
	}
	return c.JSON(200, tournaments)
}

func (api *API) createTeam(c echo.Context) error {
	var json CreateTeamRequest
	if err := c.Bind(&json); err != nil {
		return err
	}
	if err := c.Validate(json); err != nil {
		return c.JSON(400, ResponseError{Error: err.Error()})
	}
	team := Team{
		Name: json.Name,
	}
	if err := api.repo.InsertTeam(&team); err != nil {
		return err
	}
	return c.JSON(201, team)
}

func (api *API) addBot(c echo.Context) error {
	var json AddBotRequest
	if err := c.Bind(&json); err != nil {
		return err
	}
	if err := c.Validate(json); err != nil {
		return c.JSON(400, ResponseError{Error: err.Error()})
	}
	_, err := http.Get(json.Address)
	if err != nil {
		return c.JSON(400, ResponseError{Error: "cant connect bot server"})
	}
	bot := Bot{
		Name:    json.Name,
		Address: json.Address,
	}
	if err := api.repo.InsertBot(&bot); err != nil {
		return err
	}
	return c.JSON(201, bot)
}

func (api *API) deleteBot(c echo.Context) error {
	id := c.Param("id")
	bot, err := api.repo.GetBotByID(id)
	if err != nil {
		return err
	}
	if bot == nil {
		return c.JSON(404, ResponseError{Error: "bot does not exist"})
	}
	if err := api.repo.DeleteBot(bot); err != nil {
		return err
	}
	return c.JSON(204, id)
}

func (api *API) login(c echo.Context) error {
	var json LoginRequest
	if err := c.Bind(&json); err != nil {
		return err
	}
	if err := c.Validate(json); err != nil {
		return c.JSON(400, ResponseError{Error: err.Error()})
	}
	user, err := api.repo.GetUserByLogin(json.Login)
	if err != nil {
		return err
	}
	if user == nil {
		return c.JSON(404, ResponseError{Error: "неверный логин или пароль"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password)); err != nil {
		return err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	//claims["name"] = "Jon Snow"
	//claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"token": t,
	})
}

func (api *API) register(c echo.Context) error {
	var json RegisterRequest
	if err := c.Bind(&json); err != nil {
		return err
	}
	if err := c.Validate(json); err != nil {
		return c.JSON(400, ResponseError{Error: err.Error()})
	}
	{
		user, err := api.repo.GetUserByLogin(json.Login)
		if err != nil {
			return err
		}
		if user != nil {
			return c.JSON(400, ResponseError{Error: "логин занят"})
		}
	}

	password := random.String(6)
	password = json.Login
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Login:    json.Login,
		Password: string(bytes),
	}

	if err := api.repo.InsertUser(&user); err != nil {
		return err
	}
	return c.JSON(201, map[string]interface{}{
		"password": password,
	})
}

func NewAPI(repo *Repository) *API {
	return &API{
		repo: repo,
	}
}
