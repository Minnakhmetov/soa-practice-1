package main

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

type Parser interface {
	GetName() string
	Serialize(TestStruct) []byte
	Deserialize([]byte, *TestStruct)
}

type BaseParser struct {
	Name      string
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

func (t *BaseParser) GetName() string {
	return t.Name
}

func (t *BaseParser) Serialize(testStruct TestStruct) []byte {
	b, err := t.Marshal(testStruct)

	if err != nil {
		panic(err)
	}

	return b
}

func (t *BaseParser) Deserialize(b []byte, testStruct *TestStruct) {
	err := t.Unmarshal(b, testStruct)

	if err != nil {
		panic(err)
	}
}

type JsonParser struct {
	BaseParser
}

func MakeJsonParser() *JsonParser {
	return &JsonParser{
		BaseParser{
			Name:      "json",
			Marshal:   func(a any) ([]byte, error) { return json.Marshal(a) },
			Unmarshal: func(b []byte, a any) error { return json.Unmarshal(b, a) },
		},
	}
}

type MessagePackParser struct {
	BaseParser
}

func MakeMessagePackParser() *JsonParser {
	return &JsonParser{
		BaseParser{
			Name:      "messagepack",
			Marshal:   func(a any) ([]byte, error) { return msgpack.Marshal(a) },
			Unmarshal: func(b []byte, a any) error { return msgpack.Unmarshal(b, a) },
		},
	}
}
