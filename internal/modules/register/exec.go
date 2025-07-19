package register

func Exec(m string, i instance) {
	register("exec", m, i)
}

func GetExec(m string, c string) (instance, error) {
	return get("exec", m, c)
}

func GetExecList() []string {
	return getList("exec")
}

func GetExecInfoMap() map[string][]desc {
	return getInfoMap("exec")
}
