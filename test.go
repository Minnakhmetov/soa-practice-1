package main

import (
	"encoding/json"
	"log"
)

type Test interface {
	GetName() string
	Serialize(TestStruct) []byte
	Deserialize([]byte, *TestStruct)
}

type JsonTest struct{}

func MakeJsonTest() *JsonTest {
	return &JsonTest{}
}

func (t *JsonTest) GetName() string {
	return "json"
}

func (t *JsonTest) Serialize(testStruct TestStruct) []byte {
	b, err := json.Marshal(testStruct)

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func (t *JsonTest) Deserialize(b []byte, testStruct *TestStruct) {
	json.Unmarshal(b, testStruct)
}
