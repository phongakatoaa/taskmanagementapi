package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID                int    `json:"id"`
		Username          string `json:"username"`
		Password          string `json:"password"`
		EncryptedPassword string `json:"-"`
		Role              Role   `json:"role"`
		CreatedAt         string `json:"created_at"`
	}

	Role string
)

const (
	RoleEmployer = "EMPLOYER"
	RoleEmployee = "EMPLOYEE"
)

func (u *User) IsEmployer() bool {
	return u.Role == RoleEmployer
}

func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

func (u *User) EncryptPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(bytes)
	return nil
}

func (r Role) String() string {
	return string(r)
}
