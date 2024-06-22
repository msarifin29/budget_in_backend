package model

type Category struct {
	CategoryId float64 `json:"category_id"`
	Id         float64 `json:"id"`
	Title      string  `json:"title"`
	UserId     string  `json:"user_id"`
}
type CategoryReponse struct {
	CategoryId float64 `json:"category_id"`
	Id         float64 `json:"t_id"`
	Title      string  `json:"title"`
	UserId     string  `json:"user_id"`
}

func NewCategory(c Category) *Category {
	return &Category{
		CategoryId: c.CategoryId,
		Id:         c.Id,
		Title:      c.Title,
		UserId:     c.UserId,
	}
}
