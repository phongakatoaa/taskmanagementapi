package auth

import (
	"testing"
	"time"
)

func TestUser_IsEmployee(t *testing.T) {
	user := User{
		Role: RoleEmployee,
	}
	if !user.IsEmployee() {
		t.Error("User should be employee")
	}

	user.Role = RoleEmployer
	if user.IsEmployee() {
		t.Error("User should not be employee")
	}
}

func TestUser_IsEmployer(t *testing.T) {
	user := User{
		Role: RoleEmployer,
	}
	if !user.IsEmployer() {
		t.Error("User should be employer")
	}

	user.Role = RoleEmployee
	if user.IsEmployer() {
		t.Error("User should not be employer")
	}
}

func TestUser_EncryptPassword(t *testing.T) {
	user := User{
		Password: "password3",
	}
	err := user.EncryptPassword()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user.EncryptedPassword)

	t.Log(time.Now().Format(time.RFC3339))
}
