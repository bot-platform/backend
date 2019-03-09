package main

import (
	"github.com/labstack/echo"
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
	team := Team{
		Name: c.QueryParam("name"),
	}
	err := api.repo.InsertTeam(&team)
	if err != nil {
		return err
	}
	return c.JSON(201, team)
}

func NewAPI(repo *Repository) *API {
	return &API{
		repo: repo,
	}
}
