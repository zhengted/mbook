package controllers

import (
	"mbook/models"
	"mbook/utils/dynamiccache"
)

type CachedUserController struct {
	BaseController
	UcenterMember models.Member
}

func (c *CachedUserController) Prepare() {
	c.BaseController.Prepare()

	username := c.GetString(":username")
	// 从缓存读取用户信息
	cachekeyUser := "dynamiccache_user:" + username
	err := dynamiccache.ReadStruct(cachekeyUser, &c.UcenterMember)
	if err != nil {
		c.UcenterMember, _ = new(models.Member).GetByUsername(username)
		dynamiccache.WriteStruct(cachekeyUser, c.UcenterMember)
	}

	if c.UcenterMember.MemberId == 0 {
		c.Abort("404")
		return
	}
	c.Data["IsSelf"] = c.UcenterMember.MemberId == c.BaseController.Member.MemberId
	c.Data["User"] = c.UcenterMember
	c.Data["Tab"] = "share"
}

func (c *CachedUserController) Index() {

}

func (c *CachedUserController) Collection() {

}

func (c *CachedUserController) Follow() {

}

func (c *CachedUserController) Fans() {

}
