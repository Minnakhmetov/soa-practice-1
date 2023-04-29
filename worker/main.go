package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	testStructFieldLength = 2
	stringLength          = 2
)

type TestStruct struct {
	Str         string
	Integer     int64
	Float       float64
	ArrayInt    [testStructFieldLength]int64
	ArrayFloat  [testStructFieldLength]float64
	ArrayString [testStructFieldLength]string
	StringToInt map[string]int64
}

func getTestStruct() TestStruct {
	r := rand.New(rand.NewSource(424242))
	testStruct := TestStruct{}
	testStruct.Integer = r.Int63()
	testStruct.Str = getRandomString(r, 10)
	testStruct.Float = r.Float64()

	fillRandomly(
		r,
		func(r *rand.Rand) int64 { return r.Int63() },
		testStruct.ArrayInt[:],
	)
	fillRandomly(
		r,
		func(r *rand.Rand) float64 { return r.Float64() },
		testStruct.ArrayFloat[:],
	)
	fillRandomly(
		r,
		func(r *rand.Rand) string { return getRandomString(r, stringLength) },
		testStruct.ArrayString[:],
	)

	testStruct.StringToInt = map[string]int64{}

	for i := 0; i < testStructFieldLength; i++ {
		testStruct.StringToInt[getRandomString(r, stringLength)] = r.Int63()
	}

	return testStruct
}

func getAverageExecutionTimeInMs(f func()) int {
	start := time.Now()
	iterationCount := 0
	sum := 0
	for time.Since(start) < time.Second {
		iterationStart := time.Now()
		f()
		sum += int(time.Since(iterationStart).Microseconds())
		iterationCount += 1
	}
	return sum / iterationCount
}

func runTest(test Parser) string {
	averageSerializationTime := getAverageExecutionTimeInMs(
		func() { test.Serialize(getTestStruct()) },
	)
	serialized := test.Serialize(getTestStruct())
	testStruct := TestStruct{}
	averageDeserializationTime := getAverageExecutionTimeInMs(
		func() {
			test.Deserialize(serialized, &testStruct)
		},
	)
	return fmt.Sprintf(
		"%s-%d-%dms-%dms",
		test.GetName(),
		len(serialized),
		averageSerializationTime,
		averageDeserializationTime,
	)
}

func handleConnnection(makeParser func() Parser, conn net.Conn) {
	defer conn.Close()
	var parser Parser = makeParser()
	conn.Write([]byte(runTest(parser)))
}

func main() {
	var format = flag.String("format", "json", "serialization format")
	flag.Parse()

	conn, err := net.Listen("tcp", ":8000")

	makeParser := func() Parser {
		switch *format {
		case "json":
			return MakeJsonParser()
		case "msgpack":
			return MakeMessagePackParser()
		default:
			panic(fmt.Sprintf("Unknown format: %s", *format))
		}
	}

	if err != nil {
		panic(err)
	}

	for {
		conn, err := conn.Accept()

		if err != nil {
			panic(err)
		}

		go handleConnnection(makeParser, conn)
	}
}
