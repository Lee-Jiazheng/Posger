package Posger

type Databaser interface {
	AddUser(user User) (error)
	DeleteUser(user User) (error)
	SelectUser(pairs map[string]string) ([]User, error)
}

type User struct {
	// user'name or oauth2 website username
	Username	string
	// oauth2 authenticated website
	Source 		string
}

