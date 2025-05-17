package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/rs/zerolog/log"
)

var sublogger = log.With().Str("module", "dnsparser").Logger()

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

type DnsRecord struct {
	name   []byte
	rType  int
	rClass int
	ttl    int
	data   []byte
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
	// id := 0x8298
	id := rand.Intn(65536)
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

func ParseHeaderFromBuf(buf []byte) (*DnsHeader, error) {
	if len(buf) < 12 {
		return nil, errors.New("buffer too short for DNS header")
	}

	header := &DnsHeader{
		id:             int(binary.BigEndian.Uint16(buf[0:2])),
		flags:          int(binary.BigEndian.Uint16(buf[2:4])),
		numQuestions:   int(binary.BigEndian.Uint16(buf[4:6])),
		numAnswers:     int(binary.BigEndian.Uint16(buf[6:8])),
		numAuthorities: int(binary.BigEndian.Uint16(buf[8:10])),
		numAdditionals: int(binary.BigEndian.Uint16(buf[10:12])),
	}

	return header, nil
}

func DecodeDnsName(buf []byte) (string, error) {
	parts := []string{}
	i := 0
	sublogger.Info().Bytes("buf", buf).Msg("Bug")
	for length := int(buf[i]); length != 0; length = int(buf[i]) {
		part := string(buf[i : i+length+2])
		sublogger.Info().
			Int("length", length).
			Int("i", i).
			Str("part", part).Msg("DecodeDnsName")
		parts = append(parts, []string{part}...)
		i = i + length + 1
	}
	return strings.Join(parts, "."), nil
}

func DecodeDnsQuestion(buf []byte) (*DnsQuestion, error) {
	name, err := DecodeDnsName(buf[:len(buf)-4])
	if err != nil {
		return nil, err
	}

	qType := int(binary.BigEndian.Uint16(buf[len(name)+1 : len(name)+3]))
	qClass := int(binary.BigEndian.Uint16(buf[len(name)+3 : len(name)+5]))

	return &DnsQuestion{
		name:   []byte(name),
		qType:  qType,
		qClass: qClass,
	}, nil
}

func DecodeDnsRecord(buf []byte) (*DnsRecord, error) {
	name, err := DecodeDnsName(buf[:len(buf)-10])
	if err != nil {
		return nil, err
	}

	rType := int(binary.BigEndian.Uint16(buf[len(name)+1 : len(name)+3]))
	rClass := int(binary.BigEndian.Uint16(buf[len(name)+3 : len(name)+5]))
	ttl := int(binary.BigEndian.Uint32(buf[len(name)+5 : len(name)+9]))
	data := buf[len(name)+9:]

	return &DnsRecord{
		name:   []byte(name),
		rType:  rType,
		rClass: rClass,
		ttl:    ttl,
		data:   data,
	}, nil
}
