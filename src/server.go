package src

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ServerMain(port string) {
    e := echo.New()
    e.GET(".", func (c echo.Context) error {
        return c.String(http.StatusOK, "HelloWorld!")
    })

    e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
