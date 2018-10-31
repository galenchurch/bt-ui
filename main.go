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
	app.GET("/paired", routes.GetPairHandler)
	app.GET("/hsp", routes.HSPHander)
	app.GET("/a2dp", routes.A2DPHander)

	err := app.Start(":9000")
	if err != nil {
		app.Logger.Fatal(err)
	}

}

// func test() {
// 	var rad Radio

// 	ports, err := sr.GetPortsList()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if len(ports) == 0 {
// 		log.Fatal("No serial ports found!")
// 	}
// 	// Print the list of detected ports
// 	for _, port := range ports {
// 		fmt.Printf("Found port: %v\n", port)
// 	}
// 	portArg := os.Args[1:2]

// 	fmt.Printf("Passed Port: %v\n", portArg[0])

// 	rad.initPort(portArg[0])

// 	rad.getSerialLineTime(3000)

// 	rad.readTimeout(1000)
// 	fmt.Printf("buffer: %q\n", rad.UartBuf)
// 	rad.popTilReady()

// 	rad.sendLn("SET")

// 	rad.readTimeout(1000)
// 	p := rad.getPair()

// 	rad.bufPurge()
// 	p.kill(rad)
// 	p.connectHSP(rad)

// 	defer rad.closePort()
// }
