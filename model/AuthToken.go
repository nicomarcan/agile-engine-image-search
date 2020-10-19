package model

type AuthToken struct {
	Auth  bool   `json:"auth"`
	Token string `json:"token"`
}
