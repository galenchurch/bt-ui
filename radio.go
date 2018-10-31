package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/tarm/serial"
)

type Radio struct {
	UartBuf []byte
	index   int
	port    *serial.Port
}

func (r *Radio) initPort(pname string) {
	c := &serial.Config{Name: pname, Baud: 115200, ReadTimeout: time.Second * 2}
	port, err := serial.OpenPort(c)
	if err != nil {
		//log.Fatal(err)
		log.Printf("TIMEOUT: %s", err)
	}
	r.port = port
}

func (r *Radio) closePort() {
	r.port.Close()
}
func (r *Radio) bufPurge() {
	r.UartBuf = nil
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
		r.UartBuf = r.UartBuf[lend+2:]
		return ret
	}
}

func (r *Radio) readTimeout(t int64) bool {
	cur := time.Now()

	for {
		what := r.getSerialLineTime(t)
		if what == 0 {
			bdel := time.Now().Sub(cur).Nanoseconds()
			th, _ := time.ParseDuration(fmt.Sprintf("%dns", t*1000000))
			fmt.Printf("big del: %v\n", bdel)

			if bdel > th.Nanoseconds() {
				return true
			}
		} else {
			cur = time.Now()
		}
	}
	return false

}

func (r *Radio) popTilReady() bool {
	for {
		cur := r.bufPopLine()
		comp, _ := regexp.Match("READY.\r\n", cur)
		if cur == nil {
			return false
		} else if comp {
			log.Printf("System is Ready\n")
			return true
		}

	}
}

func (r *Radio) getPair() Pair {

	addrForm := regexp.MustCompile("[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}")

	for {
		cur := r.bufPopLine()
		comp, _ := regexp.Match("SET BT PAIR [[:graph:]]+ [[:alnum:]]+\r\n", cur)
		fmt.Printf("Cur: %q, comp: %v\n", cur, comp)
		if cur == nil {
			return Pair{}
		} else if comp {
			addr := addrForm.Find(cur)
			fmt.Printf("Address: %s\n", addr)
			return Pair{addr: addr}
		}
	}
}

func (r *Radio) getSerialLine() int {

	// cha := make(chan string)
	graphLine := regexp.MustCompile(`[[:print:]]*\r\n`)

	buf := make([]byte, 128)
	var line []byte

	for {
		n, err := r.port.Read(buf)
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
			return 1
		}
	}
}

func (r *Radio) sendLn(s string) bool {

	w, errw := r.port.Write([]byte(fmt.Sprintf("%s\r\n", s)))
	if errw != nil {
		log.Printf("port.Write: %v", errw)
	}

	//time.Sleep(time.Second)
	fmt.Println("Wrote", w, "bytes: ", s)
	return true
}

func (r *Radio) getSerialLineTime(t int64) int {

	// cha := make(chan string)
	graphLine := regexp.MustCompile(`[[:print:]]*\r\n`)

	buf := make([]byte, 128)
	var line []byte

	timerout := time.Now()

	for {
		n, err := r.port.Read(buf) ///Will hang here without timeout
		if err != nil {
			log.Printf("TIMEOUT: %s", err)
		}
		line = append(line, buf[:n]...)

		search := graphLine.FindAll(line, -1)
		if search != nil {
			//reset timer
			timerout = time.Now()

			for _, item := range search {
				r.bufAddLine(item)
				fmt.Printf("Added: %q\n", item)
			}
			return 1
		} else {
			del := time.Now().Sub(timerout).Nanoseconds()
			tr, _ := time.ParseDuration(fmt.Sprintf("%dns", t*1000000))

			if del > tr.Nanoseconds() {
				log.Printf("line timeout, del=%v\n", del)
				return 0
			}

		}
	}
}
