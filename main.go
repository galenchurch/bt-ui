package main

import (
	"fmt"
	"log"
	"os"

	sr "go.bug.st/serial.v1"
)

func main() {
	var rad Radio

	ports, err := sr.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	// Print the list of detected ports
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	portArg := os.Args[1:2]

	fmt.Printf("Passed Port: %v\n", portArg[0])

	rad.initPort(portArg[0])

	rad.getSerialLineTime(3000)

	rad.readTimeout(1000)
	fmt.Printf("buffer: %q\n", rad.UartBuf)
	rad.popTilReady()

	rad.sendLn("SET")

	rad.readTimeout(1000)
	p := rad.getPair()

	rad.bufPurge()
	p.kill(rad)
	p.connectHSP(rad)

	defer rad.closePort()
}
