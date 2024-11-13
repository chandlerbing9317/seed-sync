package common

type PageRequest[T any] struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type PageResponse[T any] struct {
	Page    int   `json:"page" form:"page"`
	Size    int   `json:"size" form:"size"`
	Total   int64 `json:"total" form:"total"`
	Records []T   `json:"records" form:"records"`
}
