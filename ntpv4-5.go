package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
  "log"
	"net"
	"os"
	"time"
)

const version = 5
const mode = 3

func main() {

	host := "localhost"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

  // udp set-up
  conn, err := net.Dial("udp", host+":123")
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
	defer conn.Close()
  time.Sleep(time.Millisecond * 10)
  conn.SetDeadline(time.Now().Add(2 * time.Second))

  // Biild NTPv4 Packet
	sendBuf := bytes.NewBuffer([]byte{})
	sendBuf.Write([]byte{
		0<<6 + version<<3 + mode, // flags: LI(2) Version(3) Mode(3)
		0,               // straum (8)
		6,               // interval (8)
		0xe9,            // precisiion (8)
		0, 0, 0, 0,      // Root Delay
		0, 0, 0, 1, // Root Dispersion
		0x79, 0x75, 0x6b, 0x69, // Reference ID "yuki"
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
