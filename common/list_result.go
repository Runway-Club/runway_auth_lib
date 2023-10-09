package common

type ListResult[T any] struct {
	Data    []T `json:"data"`
	EndPage int `json:"end_page"`
}
