package main

import (
	"math/rand"
)

func getRandomName() string {
	names := make([]string, 0)
	names = append(names, "Ahyoka", "Akena", "Aslan", "Eywa", "Enya",
		"Inouk", "Stone", "Boomer", "Sitka", "Diego")
	return names[rand.Intn(9)]
}
