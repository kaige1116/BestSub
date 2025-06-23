package main

import "github.com/bestruirui/bestsub/internal/start"

func main() {
	if err := start.InitConfig(); err != nil {
		panic(err)
	}
	if err := start.InitDatabase(); err != nil {
		panic(err)
	}
}
