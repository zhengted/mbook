package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"time"
)

//拼接返回到接口的图书信息
type BookData struct {
	BookId         int       `json:"book_id"`
	BookName       string    `json:"book_name"`
	Identify       string    `json:"identify"`
	OrderIndex     int       `json:"order_index"`
	Description    string    `json:"description"`
	PrivatelyOwned int       `json:"privately_owned"`
	PrivateToken   string    `json:"private_token"`
	DocCount       int       `json:"doc_count"`
	CommentCount   int       `json:"comment_count"`
	CreateTime     time.Time `json:"create_time"`
	CreateName     string    `json:"create_name"`
	ModifyTime     time.Time `json:"modify_time"`
	Cover          string    `json:"cover"`
	MemberId       int       `json:"member_id"`
	Username       int       `json:"user_name"`
	Editor         string    `json:"editor"`
	RelationshipId int       `json:"relationship_id"`
	RoleId         int       `json:"role_id"`
	RoleName       string    `json:"role_name"`
	Status         int
	Vcnt           int    `json:"vcnt"`
	Collection     int    `json:"star"`
	Score          int    `json:"score"`
	CntComment     int    `json:"cnt_comment"`
	CntScore       int    `json:"cnt_score"`
	ScoreFloat     string `json:"score_float"`
	LastModifyText string `json:"last_modify_text"`
	Author         string `json:"author"`
	AuthorURL      string `json:"author_url"`
}

func NewBookData() *BookData {
	return &BookData{}
}

// TODO:COMPLETE ME!
func (m *BookData) SelectByIdentify(identify string, memberId int) (result *BookData, err error) {
	if identify == "" || memberId <= 0 {
		return nil, errors.New("Invalid parameter")
	}
	book := NewBook()
	o := orm.NewOrm()
	err = o.QueryTable(TNBook()).Filter("identify", identify).One(book)
	if err != nil {
		return nil, err
	}

	// check authentic
	relationship := NewRelationship()
	err = o.QueryTable(TNRelationship()).Filter("book_id", book.BookId).Filter("role_id", 0).One(relationship)
	if err != nil {
		return result, err
	}

}
