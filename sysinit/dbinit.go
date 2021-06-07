package sysinit

import "github.com/astaxie/beego"
import "github.com/astaxie/beego/orm"
import _ "github.com/go-sql-driver/mysql"

// which db start
func dbinit(aliases ...string)  {

	if len(aliases) > 0 {
		for _,alias := range aliases {
			registerDatabase(alias)
			if "w" == alias {
				orm.RunSyncdb("default",false,true)
			}
		}
	} else {
		registerDatabase("w")
		orm.RunSyncdb("default",false,true)
	}

	isDev := ("dev" == beego.AppConfig.String("runmode"))

	if isDev {
		orm.Debug = isDev
	}
}

func registerDatabase(alias string) {
	if len(alias) <= 0 {
		return
	}
	dbAlias := alias // default
	if "w" == alias || "default" == alias || len(alias) <= 0 {
		dbAlias = "default"
		alias = "w"
	}

	dbName := beego.AppConfig.String("db_"+alias+"_database")
	dbUser := beego.AppConfig.String("db_"+alias+"_username")
	dbPwd := beego.AppConfig.String("db_"+alias+"_password")
	dbPort := beego.AppConfig.String("db_"+alias+"_port")
	dbHost := beego.AppConfig.String("db_"+alias+"_host")

	orm.RegisterDataBase(dbAlias,"mysql",dbUser+":"+dbPwd+"@tcp("+
		dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")


}