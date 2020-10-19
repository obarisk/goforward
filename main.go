package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"io"
	"net"
	"os/exec"
	"strings"
)

func getWslIp() (string, error) {
	bsh := "ip address show dev eth0| grep 'inet '| awk '{print $2}'| awk -F'/' '{print $1}'"
	cmd := exec.Command("bash", "-c", bsh)
	out, err := cmd.Output()
	return strings.TrimSuffix(string(out), "\n"), err
}

func forward(conn net.Conn, ippt string, ctype string) {
	prxy, err := net.Dial(ctype, ippt)
	if err != nil {
		panic(err)
	}
	go copyIO(conn, prxy)
	go copyIO(prxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}

func main() {
	usage := `goforward.
  
  Usage:
		goforward wsl tcp -l 0.0.0.0:30522 -p 22
		goforward     tcp -l 0.0.0.0:30522 -f 127.0.0.1:22
		goforward wsl udp -l 0.0.0.0:30522 -p 30522

	Options:
	  -l --listen=listen            listen
		-p --port=port                port
		-f --forward=forward          forward
	`
	opts, _ := docopt.ParseDoc(usage)
	lppt, _ := opts.String("--listen")
	ctyp := "tcp"
	if udp, _ := opts.Bool("udp"); udp {
		ctyp = "udp"
	}
	var ippt string
	if wsl, _ := opts.Bool("wsl"); wsl {
		var wip string
		var err error
		wip, err = getWslIp()
		if err != nil {
			fmt.Printf("%e", err)
		}
		wpt, _ := opts.String("--port")
		ippt = wip + ":" + wpt
	} else {
		ippt, _ = opts.String("--forward")
	}
	fmt.Println(ippt)
	ln, err := net.Listen(ctyp, lppt)
	if err != nil {
		panic(err)
	}
	for {
		con, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go forward(con, ippt, ctyp)
	}
}
