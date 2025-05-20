package models

type APIResult struct {
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
	Status  int         `json:"status"`
}
