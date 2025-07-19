package register

func Notify(m string, i instance) {
	register("notify", m, i)
}

func GetNotify(m string, c string) (instance, error) {
	return get("notify", m, c)
}

func GetNotifyList() []string {
	return getList("notify")
}

func GetNotifyInfoMap() map[string][]desc {
	return getInfoMap("notify")
}
