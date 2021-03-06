package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	radio "github.com/galenchurch/bt-ui/radio"
	"github.com/labstack/echo"
	sr "go.bug.st/serial.v1"
)

var rad radio.Radio
var p radio.Pair

type pts struct {
	Port []string
}

// FindPortsHandler returns list of avaible serial ports
// Using Serial.v1 (rather then tram/serial
func FindPortsHandler(cxt echo.Context) error {
	p, err := sr.GetPortsList()

	var port pts
	if err != nil {
		log.Fatal(err)
	}
	if len(p) == 0 {
		log.Fatal("No serial ports found!")
	}
	// Print the list of detected ports
	for _, p := range p {
		fmt.Printf("Found port: %v\n", p)
		port.Port = append(port.Port, p)
	}
	return cxt.JSON(http.StatusOK, port)
}

// InitHandler opens initial serial connection
// on port passed as form param
// It pops and reads serial buffer until READY\r\n or timeout
func InitHandler(cxt echo.Context) error {

	//p := cxt.Param("port")
	p := cxt.FormValue("port")
	log.Printf("Port: %v\n", p)

	rad.InitPort(p)

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
	l := cxt.FormValue("link")
	p.SCOClose(rad, l)

	return cxt.String(http.StatusOK, "SCO Close")
}

func ScoOpenHandler(cxt echo.Context) error {

	l := cxt.FormValue("link")
	p.SCOOpen(rad, l)

	return cxt.String(http.StatusOK, "SCO Open")
}

func ScoSwitchHandler(cxt echo.Context) error {
	la := cxt.FormValue("linka")
	// lb := cxt.FormValue("linkb")
	t := cxt.FormValue("time")
	tint, _ := strconv.Atoi(t)

	for index := 0; index < 50; index++ {
		p.SCOClose(rad, la)
		//time.Sleep(time.Duration(tint) * time.Millisecond)
		p.SCOOpen(rad, la)
		time.Sleep(time.Duration(tint) * time.Millisecond)
		// p.SCOClose(rad, lb)
		// //time.Sleep(time.Duration(tint) * time.Millisecond)
		// p.SCOOpen(rad, la)
		// time.Sleep(time.Duration(tint) * time.Millisecond)
	}

	return cxt.String(http.StatusOK, "SCO Switch")
}

func GetPairHandler(cxt echo.Context) error {
	rad.SendLn("SET")

	rad.ReadTimeout(1000)
	p = rad.GetPair()
	return cxt.String(http.StatusOK, "pair")
}

func KillHandler(cxt echo.Context) error {

	a := cxt.FormValue("address")
	log.Printf("Address: %v\n", a)
	pr := radio.Pair{Addr: []byte(a)}

	pr.Kill(rad)
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

	a := cxt.FormValue("address")
	log.Printf("Address: %v\n", a)
	pr := radio.Pair{Addr: []byte(a)}

	pr.ConnectHSP(rad)
	return cxt.String(http.StatusOK, "HSP")
}

func A2DPHander(cxt echo.Context) error {
	a := cxt.FormValue("address")
	log.Printf("Address: %v\n", a)
	pr := radio.Pair{Addr: []byte(a)}

	pr.ConnectA2DP(rad)
	return cxt.String(http.StatusOK, "A2DP")
}

func InquiryHandler(cxt echo.Context) error {
	dev := rad.Inquiry(10)
	return cxt.JSON(http.StatusOK, dev)
}

func ListHandler(cxt echo.Context) error {
	rad.List()
	return cxt.String(http.StatusOK, "List")
}

func PurgePairHandler(cxt echo.Context) error {
	rad.PurgePairs()
	return cxt.String(http.StatusOK, "Purge Pairs")
}

func ListPairsHandler(cxt echo.Context) error {
	rad.ListPairs()
	return cxt.String(http.StatusOK, "List Paris")
}

// func PairHandler(cxt echo.Context) error {

// 	a := cxt.FormValue("address")

// }

func BufferHandler(cxt echo.Context) error {

	return cxt.String(http.StatusOK, string(rad.UartBuf))
}
