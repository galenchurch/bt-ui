package radio

import (
	"fmt"
)

type Pair struct {
	addr []byte
	mode string
	link int
}

func (p *Pair) Kill(r Radio) {
	s := fmt.Sprintf("KILL %s", p.addr)
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) SCOClose(r Radio) {
	s := "SCO CLOSE"
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) ConnectHSP(r Radio) {
	s := fmt.Sprintf("CALL %s 1108 HSP-AG", p.addr)
	fmt.Println(s)
	r.SendLn(s)
	// var link []byte

	// r.ReadTimeout(1000)

	// regline := regexp.MustCompile("CALL [[:digit:]]\r\n")
	// regnum := regexp.MustCompile("[[:digit:]]")
	// for {
	// 	ln := r.BufPopLine()
	// 	if regline.Match(ln) {
	// 		link = regnum.Find(ln)
	// 		break
	// 	}
	// }

	// p.link, _ = strconv.Atoi(string(link[0]))

	// fmt.Printf("Link: %d", p.link)

	r.SendLn(fmt.Sprint("SCO OPEN"))
	r.ReadTimeout(1000)

}

func (p *Pair) ConnectA2DP(r Radio) {
	s := fmt.Sprintf("CALL %s 19 A2DP", p.addr)
	fmt.Println(s)
	r.SendLn(s)

	r.ReadTimeout(1000)

}
