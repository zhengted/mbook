package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

// 评论表
type Comments struct {
	Id         int
	Uid        int       `orm:"index"` // 用户ID
	BookId     int       `orm:"index"` // 文档项目ID
	Content    string    // 评论内容
	TimeCreate time.Time // 评论时间
}

func (m *Comments) TableName() string {
	return TNComments()
}

type BookCommentsResult struct {
	Uid        int       `json:"uid"`
	Score      int       `json:"score"`
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	Content    string    `json:"content"`
	TimeCreate time.Time `json:"time_create"` //评论时间
}

func (m *Comments) AddComments(uid, bookId int, content string) (err error) {
	var comment Comments
	// 限制评论频率
	second := 10
	sql := `select id from ` + TNComments() + ` where uid=? and time_create>? order by id desc`
	o := orm.NewOrm()
	o.Raw(sql, uid, time.Now().Add(-time.Duration(second)*time.Second)).QueryRow(&comment)
	if comment.Id > 0 {
		return errors.New(fmt.Sprintf("您距离上次发表评论时间小于 %v 秒，请稍后再发", second))
	}
	// 插入评论数据
	sql = `insert into ` + TNComments() + `(uid,book_id,content,time_create) values(?,?,?,?)`
	_, err = o.Raw(sql, uid, bookId, content, time.Now()).Exec()
	if err != nil {
		beego.Error(err.Error())
		err = errors.New("发表评论失败")
		return
	}
	// 评论数+1
	sql = `update ` + TNBook() + ` set cnt_comment=cnt_comment+1 where book_id=?`
	o.Raw(sql, bookId).Exec()
	return
}

//评论内容
func (m *Comments) BookComments(page, size, bookId int) (comments []BookCommentsResult, err error) {

	sql := `select book_id,uid,content,time_create from md_comments where book_id=? limit %v offset %v`
	sql = fmt.Sprintf(sql, size, (page-1)*size)
	o := orm.NewOrm()
	_, err = o.Raw(sql, bookId).QueryRows(&comments)
	if nil != err {
		return
	}

	// 头像昵称
	uids := []string{}
	for _, v := range comments {
		uids = append(uids, strconv.Itoa(v.Uid))
	}
	uidStr := strings.Join(uids, ",")
	sql = `select avatar,nickname from md_members where member_id in(` + uidStr + `)`
	members := []Member{}
	_, err = o.Raw(sql).QueryRows(&members)
	if err != nil {
		fmt.Println("[error] get nickname and avatar err ", err)
		return
	}
	memberMap := make(map[int]Member)
	for _, member := range members {
		memberMap[member.MemberId] = member
	}
	for _, comment := range comments {
		comment.Avatar = memberMap[comment.Uid].Avatar
		comment.Nickname = memberMap[comment.Uid].Nickname
	}

	// 评分信息
	sql = `select uid,score from md_score where book_id=? and uid in(` + uidStr + `)`
	scores := []Score{}
	_, err = o.Raw(sql, bookId).QueryRows(&scores)
	if err != nil {
		fmt.Println("[error] get score err ", err)
		return
	}
	scoreMap := make(map[int]Score)
	for _, score := range scores {
		scoreMap[score.Uid] = score
	}
	for _, comment := range comments {
		comment.Score = scoreMap[comment.Uid].Score
	}

	return
}

type Score struct {
	Id         int
	BookId     int
	Uid        int
	Score      int
	TimeCreate time.Time
}

func (m *Score) TableName() string {
	return TNScore()
}

func (m *Score) TotalUnique() [][]string {
	return [][]string{
		[]string{"Uid", "BookId"},
	}
}

//评分内容
type BookScoresResult struct {
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	Score      string    `json:"score"`
	TimeCreate time.Time `json:"time_create"` //评论时间
}

// 获取评分内容
func (m *Score) BookScore(p, listRows, bookId int) (scores []BookScoresResult, err error) {
	sql := `select s.score,s.time_create,m.avatar,m.nickname from ` + TNScore() + ` s left join ` + TNMembers() + ` m on m.member_id=s.uid where s.book_id=? order by s.id desc limit %v offset %v`
	sql = fmt.Sprintf(sql, listRows, (p-1)*listRows)
	_, err = orm.NewOrm().Raw(sql, bookId).QueryRows(&scores)
	return
}

// 查询用户对文档的评分
func (m *Score) BookScoreByUid(uid, bookId interface{}) int {
	var score Score
	orm.NewOrm().QueryTable(TNScore()).Filter("uid", uid).Filter("book_id", bookId).One(&score, "score")
	return score.Score
}

// 添加评分
func (m *Score) AddScore(uid, bookId, score int) (err error) {
	// 查询评分是否已存在
	o := orm.NewOrm()
	var scoreObj = Score{
		Uid:    uid,
		BookId: bookId,
	}
	o.Read(&scoreObj, "uid", "book_id")
	if scoreObj.Id > 0 {
		err = errors.New("您已经给当前文档打过分了")
		return
	}

	// 评分不存在，添加评分记录
	score = score * 10
	scoreObj.Score = score
	scoreObj.TimeCreate = time.Now()
	o.Insert(&scoreObj)
	if scoreObj.Id > 0 {
		// 评分添加成功，评分人数+1
		var book = Book{BookId: bookId}
		o.Read(&book, "book_id")
		if book.CntScore == 0 {
			book.CntScore = 1
			book.Score = 0
		} else {
			book.CntScore = book.CntScore + 1
		}
		book.Score = (book.Score*(book.CntScore-1) + score) / book.CntScore
		_, err = o.Update(&book, "cnt_score", "score")
		if err != nil {
			beego.Error(err.Error())
			err = errors.New("评分失败，内部错误")
		}

	}
	return
}
