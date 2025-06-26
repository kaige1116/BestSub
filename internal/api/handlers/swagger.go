package handlers

import (
	_ "github.com/bestruirui/bestsub/docs"
	"github.com/bestruirui/bestsub/internal/api/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	router.NewGroupRouter("/swagger").
		AddRoute(
			router.NewRoute("/*any", router.GET).
				Handle(ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json"))).
				WithDescription("Swagger documentation"),
		)
}
