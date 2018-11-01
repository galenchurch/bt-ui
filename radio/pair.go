package radio

import (
	"fmt"
)

type Pair struct {
	Addr []byte
	mode string
	link int
}

func (p *Pair) Kill(r Radio) {
	s := fmt.Sprintf("KILL %s", p.Addr)
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) SCOClose(r Radio, l string) {
	s := fmt.Sprintf("SCO CLOSE %s", l)
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) SCOOpen(r Radio, l string) {
	s := fmt.Sprintf("SCO OPEN %s", l)
	fmt.Println(s)
	r.SendLn(s)
}

func (p *Pair) ConnectHSP(r Radio) {
	s := fmt.Sprintf("CALL %s 1108 HSP-AG", p.Addr)
	fmt.Println(s)
	r.SendLn(s)
	r.GetSerialLineTime(1000)
	r.SendLn(fmt.Sprint("SCO OPEN"))
	r.ReadTimeout(1000)

}

func (p *Pair) ConnectA2DP(r Radio) {
	s := fmt.Sprintf("CALL %s 19 A2DP", p.Addr)
	fmt.Println(s)
	r.SendLn(s)

	r.ReadTimeout(1000)

}
