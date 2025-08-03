package models

import "github.com/google/uuid"

type UserCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Fullname string    `json:"fullname"`
	Password string    `json:"-"`
}

type UserListResponse struct {
	Data       []User `json:"data"`
	Pagination Pages  `json:"pagination"`
}
