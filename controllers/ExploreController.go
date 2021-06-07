package controllers

import (
	"github.com/astaxie/beego"
	"math"
	"mbook/models"
	"mbook/utils"
	"strconv"
)

type ExploreController struct {
	BaseController
}

func (c *ExploreController) Index() {
	var (
		cid       int
		cate      models.Category
		urlPrefix = beego.URLFor("ExploreController.Index")
	)
	if cid, _ = c.GetInt("cid"); cid > 0 {
		CateModel := new(models.Category)
		cate = CateModel.Find(cid)
		c.Data["Cate"] = cate
	}

	c.Data["Cid"] = cid
	c.TplName = "explore/index.html"

	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 24

	books, totalCount, err := models.NewBook().HomeDate(pageIndex, pageSize, cid)
	if nil != err {
		beego.Error(err)
		c.Abort("404")
	}

	if totalCount > 0 {
		urlSuffix := ""
		if cid > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(cid)
		}
		html := utils.NewPaginations(4, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	c.Data["Lists"] = books
}
