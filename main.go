package main

import (
	"github.com/astaxie/beego"
	_ "mbook/routers"
	_ "mbook/sysinit"
)

func main() {
	beego.Run()
}
