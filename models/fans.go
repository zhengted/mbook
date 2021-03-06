package models

import (
	"fmt"
)

type Fans struct {
	Id       int
	MemberId int
	FansId   int `orm:"index"`
}

type FansData struct {
	MemberId int
	Nickname string
	Avatar   string
	Account  string
}

func (m *Fans) TableName() string {
	return TNFans()
}

// 查询粉丝
func (m *Fans) FansList(mid, page, pageSize int) (fans []FansData, total int64, err error) {
	o := GetOrm("w")
	total, _ = o.QueryTable(m.TableName()).Filter("member_id", mid).Count()
	if total > 0 {
		sql := fmt.Sprintf(
			"select m.member_id member_id,m.avatar,m.account,m.nickname from "+TNMembers()+" m left join "+TNFans()+" f on m.member_id=f.fans_id where f.member_id=?  order by f.id desc limit %v offset %v",
			pageSize, (page-1)*pageSize,
		)
		_, err = o.Raw(sql, mid).QueryRows(&fans)
	}
	return
}

func (m *Fans) FollowOrCancel(mid, fansId int) (cancel bool, err error) {
	var fans Fans
	o := GetOrm("w")
	qs := o.QueryTable(m.TableName()).Filter("member_id", mid).Filter("fans_id", fansId)
	qs.One(&fans)
	if fans.Id > 0 {
		_, err = qs.Delete()
		cancel = true
	} else {
		fans.MemberId = mid
		fans.FansId = fansId
		_, err = o.Insert(&fans)
	}
	return
}

// 查询是否存在关注关系
func (m *Fans) Relation(mid, fansId interface{}) (ok bool) {
	var fans Fans
	GetOrm("w").QueryTable(TNFans()).Filter("member_id", mid).Filter("fans_id", fansId).One(&fans)
	return fans.Id != 0
}

// 查询关注的人
func (m *Fans) FollowList(fansId, page, pageSize int) (fans []FansData, total int64, err error) {
	o := GetOrm("w")
	total, _ = o.QueryTable(TNFans()).Filter("fans_id", fansId).Count() // 关注总数
	if total > 0 {
		sql := fmt.Sprintf(
			"select m.member_id member_id,m.avatar,m.account,m.nickname from "+TNMembers()+" m left join "+TNFans()+" f on m.member_id=f.member_id where f.fans_id=?  order by f.id desc limit %v offset %v",
			pageSize, (page-1)*pageSize,
		)
		_, err = o.Raw(sql, fansId).QueryRows(&fans)
	}
	return
}
