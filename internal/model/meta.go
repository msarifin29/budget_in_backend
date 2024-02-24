package model

type MetaResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MetaErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
