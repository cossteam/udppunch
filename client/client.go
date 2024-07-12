package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/yinheli/udppunch"
	"github.com/yinheli/udppunch/client/netx"
	"github.com/yinheli/udppunch/client/wg"
)

var (
	l        = log.New(os.Stdout, "", log.LstdFlags)
	iface    = flag.String("iface", "wg0", "wireguard interface")
	server   = flag.String("server", "", "udp punch server")
	interval = flag.Int("interval", 5, "interval time, 0 means not continuous")
	version  = flag.Bool("version", false, "show version")
)

func main() {
	if flag.Parse(); !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Println(udppunch.Version)
		os.Exit(0)
	}

	if *server == "" {
		l.Fatal("server is empty")
	}

	if *iface == "" {
		l.Fatal("iface is empty")
	}

	rAddr, err := net.ResolveUDPAddr("udp", *server)

	if err != nil {
		l.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		l.Fatal(err)
	}

	// handshake
	go handshake(*rAddr)

	// wait for handshake
	time.Sleep(time.Second * 2)

	// resolve
	var bak string
	for {
		clients, err := wg.GetEndpoints(*iface)
		if err != nil {
			bak = ""
			l.Print("get endpoints error:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		data := make([]byte, 0, 1+len(clients)*32)
		data = append(data, udppunch.ResolveType)

		var arr []string
		for client := range clients {
			data = append(data, client[:]...)
			arr = append(arr, client.String())
		}

		if bak != strings.Join(arr, "++") {
			bak = strings.Join(arr, "++")
			_, _ = conn.Write(data)
		}

		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			l.Print("read error: ", err)
			time.Sleep(time.Second * 10)
			continue
		}

		if n < 38 {
			time.Sleep(time.Second * 5)
			continue
		}

		peers := udppunch.ParsePeers(buf[:n])
		for _, peer := range peers {
			key, addr := peer.Parse()
			if clients[key] == addr {
				l.Print("already resolve ", key, " ", addr)
				continue
			}
			l.Print("new resolve ", key, " ", addr)
			if err = wg.SetPeerEndpoint(*iface, key, addr); err != nil {
				l.Printf("set peer (%v) endpoint error: %v", key, err)
				break
			}
		}

		if *interval == 0 {
			break
		} else {
			time.Sleep(time.Second * time.Duration(*interval))
		}
	}
}

func handshake(rAddr net.UDPAddr) {
	defer func() {
		if x := recover(); x != nil {
			l.Print("handshake error:", x)
			time.Sleep(time.Second * 10)
			handshake(rAddr)
		}
	}()

	for {
		port, err := wg.GetIfaceListenPort(*iface)
		if err != nil {
			l.Print("get interface listen-port:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		pubKey, err := wg.GetIfacePubKey(*iface)
		if err != nil {
			l.Print("get interface public key:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		doHandshake(rAddr.IP, port, uint16(rAddr.Port), pubKey)

		time.Sleep(time.Second * 25)
	}
}

func doHandshake(ip net.IP, srcPort uint16, dstPort uint16, pubKey udppunch.Key) {
	conn, err := netx.Dial(ip, srcPort, dstPort)
	if err != nil {
		l.Print("handshake dial error:", err)
		time.Sleep(time.Second * 10)
		return
	}

	defer func(conn *netx.UDPConn) {
		_ = conn.Close()
	}(conn)

	data := make([]byte, 0, 32+1)
	data = append(data, udppunch.HandshakeType)
	data = append(data, pubKey[:]...)

	_, _ = conn.Write(data)
}
