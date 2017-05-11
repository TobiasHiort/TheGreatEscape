package main

import (
  "fmt"
  "log"
  "strconv"
  "os"
)

func readStats(peopleArray []*Person) {

  file, err := os.Create("stats.txt")
  if err != nil {
    log.Fatal("Cannot create file")
  }
  defer file.Close()

  pplStats := PeopleStats(peopleArray)
  fmt.Fprintf(file, strconv.Itoa(pplStats[0]) + strconv.Itoa(pplStats[1]) + strconv.Itoa(pplStats[2]))

	mapStats := MapStats(testmap)
	fmt.Fprintf(file, strconv.Itoa(mapStats[0]))

}

func main() {

	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	ppl := [][]int{
		{0, 0},
		{0, 6},
		{2, 4}}

  testmap := TileConvert(matrix)
  pplArray:= PeopleInit(testmap,ppl)
  readStats(pplArray)

}
