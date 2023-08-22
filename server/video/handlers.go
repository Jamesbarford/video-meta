package video

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type VideoMetaHandler interface {
	CreateVideoMeta(c echo.Context) error
	ReadVideoMeta(c echo.Context) error
	UpdateVideoMeta(c echo.Context) error
	DeleteVideoMeta(c echo.Context) error
}

type VideoMetaHandlerImpl struct {
	service VideoMetaService
}

func NewVideoMetaHandler(service VideoMetaService) VideoMetaHandler {
	return &VideoMetaHandlerImpl{
		service: service,
	}
}

func (handler *VideoMetaHandlerImpl) CreateVideoMeta(c echo.Context) error {
	meta := new([]VideoMetaDataPayload)
	if err := c.Bind(&meta); err != nil {
		log.Printf("Failed to create video meta %s\n", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create metadata: "+err.Error())
	}

	videoId, err := strconv.Atoi(c.Param("videoId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create metadata videoId must be an int")
	}

	if err := handler.service.CreateVideoMeta(videoId, meta); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create metadata: "+err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}

func (handler *VideoMetaHandlerImpl) ReadVideoMeta(c echo.Context) error {
	videoId, err := strconv.Atoi(c.Param("videoId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read metadata, videoId must be an int")
	}

	metaData, err := handler.service.ReadVideoMeta(videoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request: "+err.Error())
	}

	if len(*metaData) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Failed to read metadata: Not found")
	}

	return c.JSON(http.StatusOK, metaData)
}

func (handler *VideoMetaHandlerImpl) UpdateVideoMeta(c echo.Context) error {
	meta := new(VideoMetaData)
	if err := c.Bind(meta); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to update metadata invalid payload"+err.Error())
	}

	metaId, err := strconv.Atoi(c.Param("metaId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to update metadata, metaId must be an int")
	}

	updatedMeta, err := handler.service.UpdateVideoMeta(metaId, meta)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to update metadata: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedMeta)
}

func (handler *VideoMetaHandlerImpl) DeleteVideoMeta(c echo.Context) error {
	metaId, err := strconv.Atoi(c.Param("metaId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to delete metadata, metaId must be an int")
	}

	if err := handler.service.DeleteVideoMeta(metaId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to delete metadata: "+err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
