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

func (p *Pair) SCOOpen(r Radio) {
	s := "SCO OPEN"
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) ConnectHSP(r Radio) {
	s := fmt.Sprintf("CALL %s 1108 HSP-AG", p.addr)
	fmt.Println(s)
	r.SendLn(s)
	r.GetSerialLineTime(1000)
	r.SendLn(fmt.Sprint("SCO OPEN"))
	r.ReadTimeout(1000)

}

func (p *Pair) ConnectA2DP(r Radio) {
	s := fmt.Sprintf("CALL %s 19 A2DP", p.addr)
	fmt.Println(s)
	r.SendLn(s)

	r.ReadTimeout(1000)

}
