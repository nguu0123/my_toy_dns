package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	TYPE_A   = 1
	CLASS_IN = 1
)

type DnsHeader struct {
	id             int
	flags          int
	numQuestions   int
	numAnswers     int
	numAuthorities int
	numAdditionals int
}

type DnsQuestion struct {
	name   []byte
	qType  int
	qClass int
}

func HeaderToBytes(header DnsHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint16(header.id))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(header.flags))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(header.numQuestions))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(header.numAnswers))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(header.numAuthorities))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(header.numAdditionals))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func QuestionToBytes(question DnsQuestion) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint16(question.qType))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, uint16(question.qClass))
	if err != nil {
		return nil, err
	}

	// fmt.Println("question.name", question.name)
	return append(question.name, buf.Bytes()...), nil
}

func EncodeDnsName(domain_name string) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, part := range strings.Split(domain_name, ".") {
		if len(part) > 63 {
			return nil, fmt.Errorf("label too long: %s", part)
		}
		buf.WriteByte(byte(len(part)))
		_, err := buf.Write([]byte(part))
		if err != nil {
			return nil, err
		}
	}
	buf.WriteByte(0x00)
	return buf.Bytes(), nil
}

func BuildQuery(domain_name string, record_type int) ([]byte, error) {
	name, err := EncodeDnsName(domain_name)
	if err != nil {
		return nil, err
	}
	id := 0x8298
	RECURSION_DESIRED := 1 << 8
	header := DnsHeader{id: id, numQuestions: 1, flags: RECURSION_DESIRED}
	question := DnsQuestion{name: name, qType: record_type, qClass: CLASS_IN}
	headerBytes, err := HeaderToBytes(header)
	if err != nil {
		return nil, err
	}
	questionBytes, err := QuestionToBytes(question)
	if err != nil {
		return nil, err
	}

	return append(headerBytes, questionBytes...), err
}

func main() {
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
}
