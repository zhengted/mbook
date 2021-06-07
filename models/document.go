package models

import "time"

type Document struct {
	DocumentId   int           `orm:"pk;auto;column(document_id)" json:"doc_id"`
	DocumentName string        `orm:"column(document_name);size(500)" json:"doc_name"`
	Identify     string        `orm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId       int           `orm:"column(book_id);type(int)" json:"book_id"`
	ParentId     int           `orm:"column(parent_id);type(int);default(0)" json:"parent_id"`
	OrderSort    int           `orm:"column(order_sort);default(0);type(int)" json:"order_sort"`
	Release      string        `orm:"column(release);type(text);null" json:"release"`
	CreateTime   time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId     int           `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time     `orm:"column(modify_time);type(datetime);default(null);auto_now" json:"modify_time"`
	ModifyAt     int           `orm:"column(modify_at);type(int)" json:"-"`
	Version      int64         `orm:"type(bigint);column(version)" json:"version"`
	AttachList   []*Attachment `orm:"-" json:"attach"`
	Vcnt         int           `orm:"column(vcnt);default(0)" json:"vcnt"`
	Markdown     string        `orm:"-" json:"markdown"`
}

func (m *Document) TableName() string {
	return TNDocuments()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}
