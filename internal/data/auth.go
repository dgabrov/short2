package data

type LoginAuthData struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Login  string   `json:"login"`
	Rights []string `json:"rights"`
}

type LoginData struct {
	Login    string
	Password string
}

type Person struct {
	ID         string
	Login      string
	FullName   string
	ProvidedID string
}
