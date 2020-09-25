package main

import (
	ctls "crypto/tls"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func serveFiber() *fiber.App {
	app := fiber.New()
	app.Get("/hello", func(c *fiber.Ctx) error {
		mockProcessing()
		return c.Send(message)
	})

	addr := ":" + strconv.Itoa(port)
	if isHTTPS() {
		cer, err := ctls.LoadX509KeyPair(tls.cert, tls.key)
		if err != nil {
			log.Fatal(err)
		}

		config := &ctls.Config{
			Certificates: []ctls.Certificate{cer},
		}

		ln, err := ctls.Listen("tcp", addr, config)
		if err != nil {
			log.Fatal(err)
		}

		go app.Listener(ln)
	} else {
		go app.Listen(addr)
	}
	return app
}

func shutdownFiber(app *fiber.App) {
	app.Shutdown()
}
