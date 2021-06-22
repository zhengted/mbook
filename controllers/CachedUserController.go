package controllers

import (
	"github.com/astaxie/beego"
	"mbook/common"
	"mbook/models"
	"mbook/utils"
	"mbook/utils/dynamiccache"
	"strconv"
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
	page, _ := c.GetInt("page")
	pageSize := 10
	if page < 1 {
		// 防止用户输入-1
		page = 1
	}

	// 动态缓存读取c.Data["Books"]信息
	var books []*models.BookData
	cachekeyBookList := "dynamiccache_userbook_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	totalCount, err := dynamiccache.ReadList(cachekeyBookList, &books)
	if err != nil {
		books, totalCount, _ = models.NewBook().SelectPage(page, pageSize, c.UcenterMember.MemberId, 0)
		dynamiccache.WriteList(cachekeyBookList, books, totalCount)
	}

	c.Data["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, totalCount, pageSize, page, beego.URLFor("CachedUserController.Index", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.TplName = "user/index.html"

}

func (c *CachedUserController) Collection() {
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize := 10
	var books []models.CollectionData
	cacheKeyCollectionData := "dynamiccache_usercollection_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamiccache.ReadList(cacheKeyCollectionData, &books)
	totalCount := int64(total)
	if err != nil {
		totalCount, books, _ = new(models.Collection).List(c.UcenterMember.MemberId, page, pageSize)
		dynamiccache.WriteList(cacheKeyCollectionData, books, int(totalCount))
	}

	c.Data["Books"] = books
	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Collection", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.Data["Tab"] = "collection"
	c.TplName = "user/collection.html"
}

func (c *CachedUserController) Follow() {
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize := 10
	var fans []models.FansData
	cacheKey := "dynamiccache_userfollow_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamiccache.ReadList(cacheKey, fans)
	totalCount := int64(total)
	if err != nil {
		fans, totalCount, _ = new(models.Fans).FollowList(c.UcenterMember.MemberId, page, pageSize)
		dynamiccache.WriteList(cacheKey, fans, int(totalCount))
	}

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Follow", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "follow"
	c.TplName = "user/fans.html"
}

func (c *CachedUserController) Fans() {
	page, _ := c.GetInt("page")
	pageSize := 18
	if page < 1 {
		page = 1
	}
	// fans, totalCount, _ = new(models.Fans).FansList(c.UcenterMember.MemberId, page, pageSize)
	var fans []models.FansData
	var totalCount int64
	cachekeyFansList := "dynamcache_userfans_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamiccache.ReadList(cachekeyFansList, &fans)
	totalCount = int64(total)
	if nil != err {
		fans, totalCount, _ = new(models.Fans).FansList(c.UcenterMember.MemberId, page, pageSize)
		dynamiccache.WriteList(cachekeyFansList, fans, int(totalCount))
	}
	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Fans", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "fans"
	c.TplName = "user/fans.html"
}
