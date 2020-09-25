package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func serveGin() *http.Server {
	fmt.Println(`
       .__
  ____ |__| ____
 / ___\|  |/    \
/ /_/  >  |   |  \
\___  /|__|___|  /
/_____/         \/ v1.6.3
The fastest full-featured web framework for Go. Crystal clear.
https://gin-gonic.com/
____________________________________O/_______
                                    O\       `)

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.GET("/hello", func(c *gin.Context) {
		mockProcessing()
		c.Data(http.StatusOK, "text/html;charset=UTF-8", message)
	})

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: e,
	}

	if isHTTPS() {
		go srv.ListenAndServeTLS(tls.cert, tls.key)
	} else {
		go srv.ListenAndServe()
	}
	return srv
}

func shutdownGin(s *http.Server) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	s.Shutdown(ctx)
}
