package routes

import (
	"fmt"
	"net/http"

	radio "github.com/galenchurch/bt-ui/radio"
	"github.com/labstack/echo"
)

var rad radio.Radio
var p radio.Pair

func InitHandler(cxt echo.Context) error {

	rad.InitPort("/dev/ttyUSB1")

	rad.GetSerialLineTime(3000)

	rad.ReadTimeout(1000)
	fmt.Printf("buffer: %q\n", rad.UartBuf)
	rad.PopTilReady()

	return cxt.String(http.StatusOK, string(rad.UartBuf))
}
func ClosePortHandler(cxt echo.Context) error {
	rad.ClosePort()
	return cxt.String(http.StatusOK, "Port Closed")
}

func ScoCloseHandler(cxt echo.Context) error {

	p.SCOClose(rad)
	return cxt.String(http.StatusOK, "SCO Close")
}

func ScoOpenHandler(cxt echo.Context) error {

	p.SCOOpen(rad)
	return cxt.String(http.StatusOK, "SCO Open")
}

func GetPairHandler(cxt echo.Context) error {
	rad.SendLn("SET")

	rad.ReadTimeout(1000)
	p = rad.GetPair()
	return cxt.String(http.StatusOK, "pair")
}

func KillHandler(cxt echo.Context) error {
	p.Kill(rad)
	return cxt.String(http.StatusOK, "Kill")
}

func PurgeHandler(cxt echo.Context) error {
	rad.BufPurge()
	return cxt.String(http.StatusOK, "Purge")
}

func ReadHandler(cxt echo.Context) error {
	rad.ReadTimeout(1000)
	return cxt.String(http.StatusOK, string(rad.UartBuf))
}

func HSPHander(cxt echo.Context) error {
	p.ConnectHSP(rad)
	return cxt.String(http.StatusOK, "HSP")
}

func A2DPHander(cxt echo.Context) error {
	p.ConnectA2DP(rad)
	return cxt.String(http.StatusOK, "A2DP")
}

func BufferHandler(cxt echo.Context) error {

	return cxt.String(http.StatusOK, string(rad.UartBuf))
}
