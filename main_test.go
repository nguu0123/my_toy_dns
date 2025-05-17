package main

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestHeaderToBytes(t *testing.T) {
	header := DnsHeader{id: 0x1314, flags: 0, numQuestions: 1, numAnswers: 0, numAuthorities: 0, numAdditionals: 0}
	expected := []byte{0x13, 0x14, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	result, err := HeaderToBytes(header)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestQuestionToBytes(t *testing.T) {
	name, err := hex.DecodeString("31")
	question := DnsQuestion{name: name, qType: 1, qClass: 1}

	expected, err := hex.DecodeString("3100010001")

	result, err := QuestionToBytes(question)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("\nexpected %v\n got %v", expected, result)
	}
}

func TestBuildQuery(t *testing.T) {
	domain := "example.com"
	recordType := TYPE_A
	expected := []byte{
		0x82, 0x98,
		0x01, 0x00,
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01,
		0x00, 0x01,
	}

	result, err := BuildQuery(domain, recordType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("\nexpected %v\ngot %v", expected, result)
	}
}
