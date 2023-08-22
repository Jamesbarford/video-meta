package server

import (
	"fmt"

	"github.com/Jamesbarford/video-meta/server/database"
	"github.com/Jamesbarford/video-meta/server/middleware"
	"github.com/Jamesbarford/video-meta/server/video"
	"github.com/labstack/echo/v4"
)

const API_KEY = "SUPER_SECRET"

func ServerMain(port string) {
	db, err := database.NewDbConnection(database.DbConfigFromEnvironment())
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	e := echo.New()
	e.Use(middleware.APIKeyMiddleware(API_KEY))

	repository := video.NewVideoMetaRepository(db)
	service := video.NewVideoMetaService(repository)
	handler := video.NewVideoMetaHandler(service)

	e.POST("/:videoId", handler.CreateVideoMeta)
	e.GET("/:videoId", handler.ReadVideoMeta)
	e.PUT("/:metaId", handler.UpdateVideoMeta)
	e.DELETE("/:metaId", handler.DeleteVideoMeta)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
