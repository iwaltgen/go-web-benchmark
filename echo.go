package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func serveEcho() *echo.Echo {
	e := echo.New()
	e.GET("/hello", func(c echo.Context) error {
		mockProcessing()
		return c.HTMLBlob(http.StatusOK, message)
	})

	addr := ":" + strconv.Itoa(port)
	if isHTTPS() {
		go e.StartTLS(addr, tls.cert, tls.key)
	} else {
		go e.Start(addr)
	}
	return e
}

func shutdownEcho(e *echo.Echo) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	e.Shutdown(ctx)
}
