package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/tarm/serial"
	sr "go.bug.st/serial.v1"
)

type Radio struct {
	UartBuf []byte
	index   int
}

func (r *Radio) bufAddLine(l []byte) int {
	r.UartBuf = append(r.UartBuf, l...)
	return 1
}

func (r *Radio) bufPopLine() []byte {

	lend := bytes.Index(r.UartBuf, []byte("\r\n"))
	if lend < 0 {
		return nil
	} else {
		ret := r.UartBuf[:lend+2]
		fmt.Printf("ret:%s", ret)
		r.UartBuf = r.UartBuf[lend+2:]
		return ret
	}
}

func (r *Radio) readTimeout(p *serial.Port, t int64) bool {
	cur := time.Now()

	for true {
		what := r.getSerialLineTime(p, t)
		fmt.Printf("What: %v\n", what)
		if what == 0 {

			bdel := time.Now().Sub(cur).Nanoseconds()
			th, _ := time.ParseDuration(fmt.Sprintf("%dns", t*1000000))
			fmt.Printf("big del: %v\n", bdel)

			if bdel > th.Nanoseconds() {
				fmt.Println("full timeout")
				return true
			}
		} else {
			fmt.Println("new read")
			cur = time.Now()
		}

		fmt.Println("again, again")
	}
	fmt.Println("Exit")

	return false

}

func (r *Radio) getSerialLine(p *serial.Port) int {

	// cha := make(chan string)
	graphLine := regexp.MustCompile(`[[:print:]]*\r\n`)

	buf := make([]byte, 128)
	var line []byte

	for {
		n, err := p.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		line = append(line, buf[:n]...)

		//fmt.Println(graphLine.FindAllIndex(line, -1))
		search := graphLine.FindAll(line, -1)
		if search != nil {
			fmt.Printf("Search: %q\n", search)
			for _, item := range search {
				r.bufAddLine(item)
				fmt.Printf("Added: %s\n", item)
			}
			// cha <- "tests"
			fmt.Println("Return")
			return 1
		}
	}
}

func (r *Radio) getSerialLineTime(p *serial.Port, t int64) int {

	// cha := make(chan string)
	graphLine := regexp.MustCompile(`[[:print:]]*\r\n`)

	buf := make([]byte, 128)
	var line []byte

	timerout := time.Now()

	for {
		fmt.Println("before read")
		n, err := p.Read(buf)
		fmt.Println("After Read")
		if err != nil {
			log.Fatal(err)
		}
		line = append(line, buf[:n]...)
		fmt.Printf("Line: %s -- n: %d\n", line, n)

		search := graphLine.FindAll(line, -1)
		if search != nil {
			//reset timer
			timerout = time.Now()

			fmt.Printf("Search: %q\n", search)
			for _, item := range search {
				r.bufAddLine(item)
				fmt.Printf("Added: %s\n", item)
			}
			fmt.Println("Return")
			return 1
		} else {
			del := time.Now().Sub(timerout).Nanoseconds()
			tr, _ := time.ParseDuration(fmt.Sprintf("%dns", t*1000000))
			fmt.Printf("del: %v\n", del)
			if del > tr.Nanoseconds() {
				fmt.Println("line timeout")
				return 0
			}

		}
	}
}

type Pair struct {
	addr []byte
}

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

	c := &serial.Config{Name: portArg[0], Baud: 115200}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// rad.readAllTimeout(port, 1000)
	//rad.readTimeout(port, 1000)
	// rad.getSerialLine(port)
	// rad.getSerialLine(port)
	// rad.getSerialLine(port)
	// rad.getSerialLine(port)

	rad.getSerialLineTime(port, 3000)
	// rad.getSerialLineTime(port, 1000)
	// rad.getSerialLineTime(port, 1000)
	// rad.getSerialLineTime(port, 1000)

	rad.readTimeout(port, 1000)
	fmt.Printf("buffer: %q\n", rad.UartBuf)

	//getSerialLineTimeout(port, 1)
	//getSerialLine(port)
	//getSerialLine(port)
	//getSerialLineTimeout(port, 1)

	s := "SET\r\n"
	w, err_w := port.Write([]byte(s))
	if err_w != nil {
		log.Fatalf("port.Write: %v", err_w)
	}

	//time.Sleep(time.Second)
	fmt.Println("Wrote", w, "bytes: ", s)
	rad.readTimeout(port, 1000)
	//	getSerialLineTimeout(port, 1)

	// Make sure to close it later.
	fmt.Printf("buffer: %q\n", rad.UartBuf)

	defer port.Close()
}
