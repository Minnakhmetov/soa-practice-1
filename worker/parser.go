package main

import (
	"encoding/json"
	"os"

	"github.com/hamba/avro/v2"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

type Parser interface {
	GetName() string
	Serialize(TestStruct) []byte
	Deserialize([]byte, *TestStruct)
}

type BaseParser struct {
	Name      string
	Marshal   func(TestStruct) ([]byte, error)
	Unmarshal func([]byte, *TestStruct) error
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
			Marshal:   func(a TestStruct) ([]byte, error) { return json.Marshal(a) },
			Unmarshal: func(b []byte, a *TestStruct) error { return json.Unmarshal(b, a) },
		},
	}
}

type MessagePackParser struct {
	BaseParser
}

func MakeMessagePackParser() *MessagePackParser {
	return &MessagePackParser{
		BaseParser{
			Name:      "msgpack",
			Marshal:   func(a TestStruct) ([]byte, error) { return msgpack.Marshal(a) },
			Unmarshal: func(b []byte, a *TestStruct) error { return msgpack.Unmarshal(b, a) },
		},
	}
}

type YAMLParser struct {
	BaseParser
}

func MakeYAMLParser() *YAMLParser {
	return &YAMLParser{
		BaseParser{
			Name:      "yaml",
			Marshal:   func(a TestStruct) ([]byte, error) { return yaml.Marshal(a) },
			Unmarshal: func(b []byte, a *TestStruct) error { return yaml.Unmarshal(b, a) },
		},
	}
}

type AvroParser struct {
	BaseParser
}

func MakeAvroParser() *AvroParser {
	b, err := os.ReadFile("./avro_schema.json")
	if err != nil {
		panic(err)
	}

	schema, err := avro.Parse(string(b))
	if err != nil {
		panic(err)
	}

	return &AvroParser{
		BaseParser{
			Name:      "avro",
			Marshal:   func(a TestStruct) ([]byte, error) { return avro.Marshal(schema, a) },
			Unmarshal: func(b []byte, a *TestStruct) error { return avro.Unmarshal(schema, b, a) },
		},
	}
}

type ProtobufParser struct {
	BaseParser
}

func MakeProtobufParser() *ProtobufParser {
	return &ProtobufParser{
		BaseParser{
			Name: "protobuf",
			Marshal: func(a TestStruct) ([]byte, error) {
				testStructPb := TestStructPb{}
				testStructPb.Str = a.Str
				testStructPb.Integer = a.Integer
				testStructPb.Float64 = a.Float
				testStructPb.ArrayInt = a.ArrayInt
				testStructPb.ArrayFloat = a.ArrayFloat
				testStructPb.ArrayString = a.ArrayString
				testStructPb.StringToInt = a.StringToInt
				return proto.Marshal(&testStructPb)
			},
			Unmarshal: func(b []byte, ts *TestStruct) error {
				testStructPb := TestStructPb{}
				err := proto.Unmarshal(b, &testStructPb)
				if err != nil {
					return err
				}
				ts.Str = testStructPb.Str
				ts.Integer = testStructPb.Integer
				ts.Float = testStructPb.Float64
				ts.ArrayInt = testStructPb.ArrayInt
				ts.ArrayFloat = testStructPb.ArrayFloat
				ts.ArrayString = testStructPb.ArrayString
				ts.StringToInt = testStructPb.StringToInt
				return nil
			},
		},
	}
}
