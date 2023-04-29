package main

import (
	"net"
	"strings"
)

type Connection struct {
	conn net.PacketConn
}

var formatToWorkerAddress map[string]string = map[string]string{
	"json":    "json-worker:8000",
	"msgpack": "msgpack-worker:8000",
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

func handleMessage(msg string) string {
	tokens := strings.Split(strings.Trim(msg, "\n"), " ")

	if tokens[0] != "get_result" {
		return "Unknown command: did you mean \"get_result\"?"
	}
	if len(tokens) > 2 {
		return "Too many arguments: get_result accepts only one argument"
	}
	if len(tokens) < 2 {
		return "Too few arguments: get_result requires an argument"
	}

	address, ok := formatToWorkerAddress[tokens[1]]

	if !ok {
		return "Unknown argument"
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
