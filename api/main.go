package main

import (
	"fmt"
	"net"
	"strings"
)

type Connection struct {
	conn net.PacketConn
}

var formatToWorkerAddress map[string]string = map[string]string{
	"json":     "json-worker:8000",
	"msgpack":  "msgpack-worker:8000",
	"yaml":     "yaml-worker:8000",
	"avro":     "avro-worker:8000",
	"protobuf": "protobuf-worker:8000",
}

func callWorker(address string) string {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)

	if err != nil {
		panic(err)
	}

	return string(buffer)
}

func getSupportedFormats() []string {
	var formats []string = make([]string, 0, len(formatToWorkerAddress))
	for k := range formatToWorkerAddress {
		formats = append(formats, k)
	}
	return formats
}

func handleMessage(msg string) string {
	tokens := strings.Split(strings.Trim(msg, "\n"), " ")

	addExample := func(msg string) string {
		supportedFormats := getSupportedFormats()
		return fmt.Sprintf("%s\nUsage: get_result format\nAll supported formats: %s",
			msg, strings.Join(supportedFormats, ", "))
	}

	if tokens[0] != "get_result" {
		return addExample("Unknown command: did you mean \"get_result\"?")
	}
	if len(tokens) > 2 {
		return addExample("Too many arguments: get_result accepts only one argument")
	}
	if len(tokens) < 2 {
		return addExample("Too few arguments: get_result requires an argument")
	}

	address, ok := formatToWorkerAddress[tokens[1]]

	if !ok {
		return addExample("Unknown format")
	}

	return callWorker(address)
}

func main() {

	conn, err := net.ListenPacket("udp", ":2000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		go func(sender net.Addr, msg string) {
			result := handleMessage(msg)
			conn.WriteTo([]byte(result+"\n"), addr)
		}(addr, string(buffer[:n]))
	}

}
