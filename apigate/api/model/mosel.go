package model

type User struct {
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
}

type ResponseError struct {
	Error interface{} `json:"error"`
}

type ServerError struct {
	Status string `json:"status"`
	Message string `json:"message"`
}