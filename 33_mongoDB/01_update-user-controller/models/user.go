package models

type User struct {
	Name   string
	Gender string
	Age    int
	ID     string
}

type JsonErr struct {
	Type   string
	Status int
	Error  string
}
