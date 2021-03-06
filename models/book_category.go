package models

import (
	"strconv"
)

type BookCategory struct {
	Id         int //自增主键
	BookId     int //书籍id
	CategoryId int //分类id
}

func (m *BookCategory) TableName() string {
	return TNBookCategory()
}

func (m *BookCategory) SelectByBookId(book_id int) (cates []Category, rows int64, err error) {
	o := GetOrm("w")
	sql := "select c.* from " + TNCategory() + " c left join " + TNBookCategory() + " bc on c.id=bc.category_id where bc.book_id=?"
	rows, err = o.Raw(sql, book_id).QueryRows(&cates)
	return
}

func (m *BookCategory) SetBookCates(bookId int, cids []string) {
	if len(cids) == 0 {
		return
	}

	var (
		cates             []Category
		tableCategory     = TNCategory()
		tableBookCategory = TNBookCategory()
	)

	o := GetOrm("w")
	o.QueryTable(tableCategory).Filter("id__in", cids).All(&cates, "id", "pid")

	cidMap := make(map[string]bool)
	for _, cate := range cates {
		cidMap[strconv.Itoa(cate.Pid)] = true
		cidMap[strconv.Itoa(cate.Id)] = true
	}
	cids = []string{}
	for cid, _ := range cidMap {
		cids = append(cids, cid)
	}

	o.QueryTable(tableBookCategory).Filter("book_id", bookId).Delete()
	var bookCates []BookCategory
	for _, cid := range cids {
		cidNum, _ := strconv.Atoi(cid)
		bookCate := BookCategory{
			CategoryId: cidNum,
			BookId:     bookId,
		}
		bookCates = append(bookCates, bookCate)
	}

	if l := len(bookCates); l > 0 {
		o.InsertMulti(l, &bookCates)
	}
	go CountCategory()

}
