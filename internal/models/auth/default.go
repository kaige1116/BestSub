package auth

var authConfigData = Data{
	UserName: "admin",
	Password: "admin",
}

func Default() Data {
	return authConfigData
}
