package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func flagUsage() {
	usageText := `go run ./mock-client.go COMMAND host

COMMAND
- v4    : Normal NTPv4 Packet
- v4-ue : NTPv4 packet with unknown extension
- v4-5  : NTPv4 packet that version field is 5
- v5    : NTPv5 packet (draft-mlichvar-ntp-ntpv5-00)

HOST: target host
`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {
	flag.Usage = flagUsage

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	host := "localhost"
	if len(os.Args) > 2 {
		host = os.Args[2]
	}

	// udp set-up
	conn, err := net.Dial("udp", host+":123")
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
	defer conn.Close()
	time.Sleep(time.Millisecond * 10)
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Biild NTPv4 Packet
	sendBuf := bytes.NewBuffer([]byte{})
	switch os.Args[1] {
	case "v4", "v4-ue", "v4-5":
		version := byte(4)
		mode := byte(3)

		if os.Args[1] == "v4-5" {
			version = byte(5)
		}

		sendBuf.Write([]byte{
			0<<6 + version<<3 + mode, // flags: LI(2) Version(3) Mode(3)
			0,                        // straum (8)
			6,                        // interval (8)
			0xe9,                     // precisiion (8)
			0, 0, 0, 0,               // Root Delay
			0, 0, 0, 1, // Root Dispersion
      //0x79, 0x75, 0x6b, 0x69, // Reference ID "yuki"
      0x0, 0x1, 0x2, 0x0, // Reference ID "yuki"
			0, 0, 0, 0, 0, 0, 0, 0, // Reference Timestamp
			0, 0, 0, 0, 0, 0, 0, 0, // Origin Timestamp
			0, 0, 0, 0, 0, 0, 0, 0, // Receive Timestamp
		})

		// Write Transmit Timestamp
		utc, _ := time.LoadLocation("UTC")
		t1 := time.Date(1900, 1, 1, 0, 0, 0, 0, utc)
		t2 := time.Now().UTC()

		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(t2.Sub(t1).Seconds()))
		sendBuf.Write(b)
		sendBuf.Write([]byte{0, 0, 0, 0})

		if os.Args[1] == "v4-ue" {
			sendBuf.Write([]byte{0x00, 0x02, 0, 16}) // type(16), length(16)
			sendBuf.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8}) // value
		}

	case "v5":
		version := byte(5)
		mode := byte(3)
		sendBuf.Write([]byte{
			0<<6 + version<<3 + mode, // LI(2) Version(3) Mode(3)
			0,                        // scale (4) straum(4)
			6,                        // interval (8)
			0,                        // precisiion (8)
			0,                        // flags (8)
			0x80, 0,                  // era (16)
			0, 0, 0, 0, // Root Delay
			0, 0, 0, 1, // Root Dispersion
			0, 0, 0, 0, 0, 0, 0, 0, // Server Cookie
			0, 0, 0, 0, 0, 0, 0, 0, // Client Cookie
			0, 0, 0, 0, 0, 0, 0, 0, // Receive Timestamp
		})

		// Write Transmit Timestamp
		utc, _ := time.LoadLocation("UTC")
		t1 := time.Date(1900, 1, 1, 0, 0, 0, 0, utc)
		t2 := time.Now().UTC()

		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(t2.Sub(t1).Seconds()))
		sendBuf.Write(b)
		sendBuf.Write([]byte{0, 0, 0, 0})

		// Add Dummy Extension
	  sendBuf.Write([]byte{0, 10, 0, 8, 0, 0, 0, 1})
	default:
		flag.Usage()
	}

	// Send
	_, err = conn.Write(sendBuf.Bytes())
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}

	recvBuf := make([]byte, 1024)
	_, err = conn.Read(recvBuf)
	if err != nil {
		if err.(net.Error).Timeout() {
			fmt.Printf("%s response version: timeout\n", host)
		}
		os.Exit(0)
	}

	fmt.Printf("%s response version: %d\n", host, (int(recvBuf[0])>>3)&7)
}
