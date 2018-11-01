package radio

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

func (r *Radio) InitPort(pname string) {
	c := &serial.Config{Name: pname, Baud: 115200, ReadTimeout: time.Second * 2}
	port, err := serial.OpenPort(c)
	if err != nil {
		//log.Fatal(err)
		log.Printf("TIMEOUT: %s", err)
	}
	r.port = port
}

func (r *Radio) ClosePort() {
	r.port.Close()
}
func (r *Radio) BufPurge() {
	r.UartBuf = nil
}

func (r *Radio) BufAddLine(l []byte) int {
	r.UartBuf = append(r.UartBuf, l...)
	return 1
}

func (r *Radio) BufPopLine() []byte {

	lend := bytes.Index(r.UartBuf, []byte("\r\n"))
	if lend < 0 {
		return nil
	} else {
		ret := r.UartBuf[:lend+2]
		r.UartBuf = r.UartBuf[lend+2:]
		return ret
	}
}

func (r *Radio) ReadTimeout(t int64) bool {
	cur := time.Now()

	for {
		what := r.GetSerialLineTime(t)
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

func (r *Radio) PopTilReady() bool {
	for {
		cur := r.BufPopLine()
		comp, _ := regexp.Match("READY.\r\n", cur)
		if cur == nil {
			return false
		} else if comp {
			log.Printf("System is Ready\n")
			return true
		}

	}
}

func (r *Radio) GetPair() Pair {

	addrForm := regexp.MustCompile("[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}")

	for {
		cur := r.BufPopLine()
		comp, _ := regexp.Match("SET BT PAIR [[:graph:]]+ [[:alnum:]]+\r\n", cur)
		fmt.Printf("Cur: %q, comp: %v\n", cur, comp)
		if cur == nil {
			return Pair{}
		} else if comp {
			addr := addrForm.Find(cur)
			fmt.Printf("Address: %s\n", addr)
			return Pair{Addr: addr}
		}
	}
}
func (r *Radio) List() {
	r.BufPurge()
	r.SendLn("LIST")
	r.ReadTimeout(1000)
}

func (r *Radio) PurgePairs() {
	r.SendLn("SET BT PAIR *")
}

func (r *Radio) ListPairs() {
	r.BufPurge()
	r.SendLn("SET BT PAIR")
	r.ReadTimeout(1000)
}

func (r *Radio) Inquiry(t int64) []Device {
	r.BufPurge()
	r.SendLn(fmt.Sprintf("INQUIRY %d NAME", t))
	r.ReadTimeout((t + 1) * 1000)
	return r.ParseInquiry()
}

type Device struct {
	Name    string
	Address string
}

func (r *Radio) ParseInquiry() []Device {
	var dev []Device
	nameln := regexp.MustCompile(`^NAME [[:print:]]+\r\n`)
	addrForm := regexp.MustCompile("[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}:[[:alnum:]]{2}")
	nameForm := regexp.MustCompile(`\"[[:print:]]+\"\r\n`)

	i := 0
	for {
		ln := r.BufPopLine()
		if nameln.Match(ln) {
			d := Device{Name: string(nameForm.Find(ln)), Address: string(addrForm.Find(ln))}
			dev = append(dev, d)
			i++
		} else if ln == nil {
			return dev
		}
	}

}

func (r *Radio) GetSerialLine() int {

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
				r.BufAddLine(item)
				fmt.Printf("Added: %s\n", item)
			}
			return 1
		}
	}
}

func (r *Radio) SendLn(s string) bool {

	w, errw := r.port.Write([]byte(fmt.Sprintf("%s\r\n", s)))
	if errw != nil {
		log.Printf("port.Write: %v", errw)
	}

	//time.Sleep(time.Second)
	fmt.Println("Wrote", w, "bytes: ", s)
	return true
}

func (r *Radio) GetSerialLineTime(t int64) int {

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
				r.BufAddLine(item)
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
