package main

import "math/rand"

func getRandomChar(r *rand.Rand) string {
	return string(rune(int('a') + r.Intn(26)))
}

func getRandomString(r *rand.Rand, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result = result + string(getRandomChar(r))
	}
	return result
}

func fillRandomly[T any](r *rand.Rand, getRandomElement func(r *rand.Rand) T, array []T) {
	for i := 0; i < len(array); i++ {
		array[i] = getRandomElement(r)
	}
}
