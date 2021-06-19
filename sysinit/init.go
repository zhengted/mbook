package sysinit

func init() {
	sysinit()
	dbinit("w")
	dbinit("r")
	//dbinit("uaw","uar")
}
