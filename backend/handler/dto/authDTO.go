package dto

import (
	"astor/backend/entity"
	"time"
)

type RequestRegister struct {
	RoleID    uint      `json:"role_id" binding:"required"`
	Npm       string    `json:"npm" binding:"required,len=10,numeric"`
	Firstname string    `json:"firstname" binding:"required"`
	Lastname  string    `json:"lastname"`
	Birthday  time.Time `json:"birthday"`
	Address   string    `json:"address"`
	Username  string    `json:"username" binding:"required"`
	Password  string    `json:"password" binding:"required,min=6"`
}

type RequestLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResponseRole struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ResponseUserLogin struct {
	ID        uint         `json:"id"`
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Username  string       `json:"username"`
	Role      ResponseRole `json:"role"`
}
type ResponseLogin struct {
	Token string            `json:"token"`
	User  ResponseUserLogin `json:"user"`
}

type ResponseUserDetail struct {
	ID        uint         `json:"id"`
	Npm       string       `json:"npm"`
	Firstname string       `json:"firstname"`
	Lastname  string       `json:"lastname"`
	Birthday  string       `json:"birthday"`
	Address   string       `json:"address"`
	Username  string       `json:"username"`
	Role      ResponseRole `json:"role"`
}

func ToResponseLogin(user *entity.Users, token string) ResponseLogin {
	return ResponseLogin{
		Token: token,
		User: ResponseUserLogin{
			ID:        user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Role: ResponseRole{
				ID:   user.Role.ID,
				Name: user.Role.Name,
			},
		},
	}
}

func ToResponseUserList(users []entity.Users) []ResponseUserDetail {
	res := make([]ResponseUserDetail, len(users))
	for i, u := range users {
		bday := ""
		if !u.Birthday.IsZero() {
			bday = u.Birthday.Format("2006-01-02")
		}
		res[i] = ResponseUserDetail{
			ID:        u.ID,
			Npm:       u.Npm,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Birthday:  bday,
			Address:   u.Address,
			Username:  u.Username,
			Role: ResponseRole{
				ID:   u.Role.ID,
				Name: u.Role.Name,
			},
		}
	}
	return res
}

func (r *RequestRegister) ToEntity() *entity.Users {
	return &entity.Users{
		RoleID:    r.RoleID,
		Npm:       r.Npm,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Birthday:  r.Birthday,
		Address:   r.Address,
		Username:  r.Username,
		Password:  r.Password,
	}
}
