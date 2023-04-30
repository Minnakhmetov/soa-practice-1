package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func getAverageExecutionTimeInMs(f func()) int {
	start := time.Now()
	iterationCount := 0
	sum := 0
	for time.Since(start) < time.Second {
		iterationStart := time.Now()
		f()
		sum += int(time.Since(iterationStart).Nanoseconds())
		iterationCount += 1
	}
	return sum / iterationCount
}

func runTest(test Parser) string {
	averageSerializationTime := getAverageExecutionTimeInMs(
		func() { test.Serialize(MakeTestStruct()) },
	)
	serialized := test.Serialize(MakeTestStruct())
	testStruct := TestStruct{}
	averageDeserializationTime := getAverageExecutionTimeInMs(
		func() {
			test.Deserialize(serialized, &testStruct)
		},
	)
	return fmt.Sprintf(
		"%s-%d-%dus-%dus",
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

	log.Printf("Got format=%s", *format)
	log.Print("Ready to run tests")

	makeParser := func() Parser {
		switch *format {
		case "json":
			return MakeJsonParser()
		case "msgpack":
			return MakeMessagePackParser()
		case "yaml":
			return MakeYAMLParser()
		case "avro":
			return MakeAvroParser()
		case "protobuf":
			return MakeProtobufParser()
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

		log.Print("Got new connection")

		go handleConnnection(makeParser, conn)
	}
}
