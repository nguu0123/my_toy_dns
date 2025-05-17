package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	domain_name := os.Args[1]
	query, err := BuildQuery(domain_name, 1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// fmt.Printf("%b\n", query)
	serverAddr := "8.8.8.8:53"

	// Create a UDP connection
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DNS server: %w", err))
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write(query)
	if err != nil {
		panic(fmt.Errorf("failed to send DNS query: %w", err))
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(fmt.Errorf("failed to read DNS response: %w", err))
	}

	response := buf[:n]
	fmt.Printf("Received %d bytes:\n% x\n", n, response)

	dnsHeader, err := ParseHeaderFromBuf(response[:12])
	if err != nil {
		panic(fmt.Errorf("failed to parse DNS header: %w", err))
	}

	fmt.Println(dnsHeader)

	dnsQuestion, err := DecodeDnsQuestion(response[12:])
	fmt.Println(dnsQuestion)

	dnsRecord, err := DecodeDnsRecord(response[12+len(dnsQuestion.name)+4:])
	fmt.Println(dnsRecord)
	// name := []byte{0x03, 'w', 'w', 'w', 0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, 0x00, 0x01, 0x00, 0x01}
	//
	// question := DnsQuestion{name: name, qType: 1, qClass: 1}

	// data, err := QuestionToBytes(question)
	// if err != nil {
	// 	panic(err)
	// }
}
