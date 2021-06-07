package models

type BookCategory struct {
	Id         int //自增主键
	BookId     int //书籍id
	CategoryId int //分类id
}

func (m *BookCategory) TableName() string {
	return TNBookCategory()
}
