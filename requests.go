package main

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required"`
}

type AddBotRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required,url"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Login string `json:"login" validate:"required"`
}
