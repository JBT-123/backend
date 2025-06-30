package models

import "time"

//this is just a placeholder for model like structs so that it is centralised and not declared in the file itself
type UserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name`
	LastName  string `json:"last_name`
	Gender    string `json:"gender"`
	BirthDate time.Time `json:"birth_date`
	Bio string `json:"bio"`
	ProfilePicUrl string `json:"profile_pic_url"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Password is not included in JSON responses
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserPref struct {
	PrefGender string `json:"preferred_genders"`
	Interest string `json:"interests"`
}

type UserConnectionRes struct {
	UserInfo UserInfo `json:"user_info"`
	UserPref UserPref `json:"user_pref"`
}