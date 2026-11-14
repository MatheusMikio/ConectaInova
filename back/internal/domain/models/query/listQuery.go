package models

type PageQuery struct {
	Page uint `form:"page"`
	Size uint `form:"size"`
}

const (
	DefaultPage = 1
	DefaultSize = 12
	MaxSize     = 30
)

func NewPaginationQuery(page, size uint) PageQuery {
	if page < DefaultPage {
		page = DefaultPage
	}
	if size == 0 || size > MaxSize {
		size = DefaultSize
	}
	return PageQuery{Page: page, Size: size}
}
