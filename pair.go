package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Pair struct {
	addr []byte
	mode string
	link int
}

func (p *Pair) kill(r Radio) {
	s := fmt.Sprintf("KILL %s", p.addr)
	fmt.Println(s)
	r.sendLn(s)
}

func (p *Pair) connectHSP(r Radio) {
	s := fmt.Sprintf("CALL %s 1108 HSP-AG", p.addr)
	fmt.Println(s)
	r.sendLn(s)
	var link []byte

	r.readTimeout(1000)

	regline := regexp.MustCompile("CALL [[:digit:]]\r\n")
	regnum := regexp.MustCompile("[[:digit:]]")
	for {
		ln := r.bufPopLine()
		if regline.Match(ln) {
			link = regnum.Find(ln)
			break
		}
	}

	p.link, _ = strconv.Atoi(string(link[0]))

	fmt.Printf("Link: %d", p.link)

	r.sendLn(fmt.Sprintf("SCO OPEN %d", p.link))
	r.readTimeout(1000)

}
