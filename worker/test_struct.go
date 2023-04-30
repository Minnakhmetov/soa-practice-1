package main

import "math/rand"

const (
	testStructArrayLength = 10
	stringLength          = 7
)

type TestStruct struct {
	Str         string           `avro:"str"`
	Integer     int64            `avro:"integer"`
	Float       float64          `avro:"float"`
	ArrayInt    []int64          `avro:"array_int"`
	ArrayFloat  []float64        `avro:"array_float"`
	ArrayString []string         `avro:"array_string"`
	StringToInt map[string]int64 `avro:"string_to_int"`
}

func MakeTestStruct() TestStruct {
	r := rand.New(rand.NewSource(424242))
	testStruct := TestStruct{}
	testStruct.Integer = r.Int63()
	testStruct.Str = getRandomString(r, 10)
	testStruct.Float = r.Float64()
	testStruct.ArrayInt = make([]int64, testStructArrayLength)
	testStruct.ArrayFloat = make([]float64, testStructArrayLength)
	testStruct.ArrayString = make([]string, testStructArrayLength)

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

	for i := 0; i < testStructArrayLength; i++ {
		testStruct.StringToInt[getRandomString(r, stringLength)] = r.Int63()
	}

	return testStruct
}
