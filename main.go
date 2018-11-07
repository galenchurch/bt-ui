package main

import (
	"net/http"

	"github.com/galenchurch/bt-ui/routes"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func main() {
	app := echo.New()

	app.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	app.Use(mw.StaticWithConfig(mw.StaticConfig{
		Root: "public",
	}))

	app.GET("/ping", func(cxt echo.Context) error {
		return cxt.String(http.StatusOK, "pong")
	})

	app.File("/", "ui.html")

	app.GET("/init", routes.InitHandler)
	app.GET("/kill", routes.KillHandler)
	app.GET("/purge", routes.PurgeHandler)
	app.GET("/close", routes.ClosePortHandler)
	app.GET("/ports", routes.FindPortsHandler)
	app.GET("/inquiry", routes.InquiryHandler)
	// app.GET("/pair", routes.PairHandler)
	app.GET("/list", routes.ListHandler)

	app.GET("/listpairs", routes.ListPairsHandler)
	app.GET("/purgepairs", routes.PurgePairHandler)
	app.GET("/switch", routes.ScoSwitchHandler)

	app.GET("/scoclose", routes.ScoCloseHandler)
	app.GET("/scoopen", routes.ScoOpenHandler)
	app.GET("/paired", routes.GetPairHandler)
	app.GET("/hsp", routes.HSPHander)
	app.GET("/a2dp", routes.A2DPHander)

	app.GET("/buffer", routes.BufferHandler)
	app.GET("/read", routes.ReadHandler)

	err := app.Start(":9000")
	if err != nil {
		app.Logger.Fatal(err)
	}

}
