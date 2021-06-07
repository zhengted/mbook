package sysinit

import (
	"github.com/astaxie/beego"
	"mbook/utils"
	"path/filepath"
	"strings"
)

func sysinit() {
	uploads := filepath.Join(
		"./", "uploads")
	beego.BConfig.WebConfig.StaticDir["/uploads"] = uploads

	// 注册前端使用函数
	registerFunctions()
}

func registerFunctions() {
	beego.AddFuncMap("cdnjs", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdnjs", "")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
	beego.AddFuncMap("cdncss", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdncss", "")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})

	beego.AddFuncMap("inMap", utils.InMap)
	// 临时空函数
	beego.AddFuncMap("getUsernameByUid", func() string {
		return ""
	})
	beego.AddFuncMap("getNicknameByUid", func() string {
		return ""
	})
	beego.AddFuncMap("IsFollow", func(id int) bool {
		return true
	})
	beego.AddFuncMap("doesCollection", func(id int) int {
		return id
	})

	//	//用户是否收藏了文档
	//beego.AddFuncMap("doesCollection", new(models.Collection).DoesCollection)
	//	beego.AddFuncMap("scoreFloat", utils.ScoreFloat)
	beego.AddFuncMap("showImg", utils.ShowImg)
	//beego.AddFuncMap("IsFollow", new(models.Fans).Relation)

	beego.AddFuncMap("isubstr", utils.Substr)
}
