package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	_ "mbook/routers"
	_ "mbook/sysinit"
	"mbook/utils/pagecache"
)

func main() {
	task := toolbox.NewTask("clear_expired_cache", "2 * * * * *", func() error {
		pagecache.ClearExpireFile()
		return nil
	})
	toolbox.AddTask("mbook_task", task)
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()
}
