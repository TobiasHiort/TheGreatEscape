package main

import (
  "fmt"
  "log"
  "strconv"
  "os"
)

func readStats(peopleArray []*Person, inmap [][]tile) {

  file, err := os.Create("stats.txt")
  if err != nil {
    log.Fatal("Cannot create file")
  }
  defer file.Close()

  pplStats := PeopleStats(peopleArray)
  //alive dead injured
  fmt.Fprintf(file, strconv.Itoa(pplStats[0]) + strconv.Itoa(pplStats[1]) + strconv.Itoa(pplStats[2]))

	mapStats := MapStats(inmap)
	fmt.Fprintf(file, strconv.Itoa(mapStats[0]))

}

//TODO most used door
func exitStats(peopleArray []*Person, inmap [][]tile) []int {

  doors := doorCoord(inmap)
  var doorStats []int
  index := (len(peopleArray) - 2)

  for  i := 0; i < (len(peopleArray)); i++ {
    tmp := []int {peopleArray[i].path[index].xCoord, peopleArray[i].path[index].yCoord}
    for j := 0; j < (len(doors)); j++ {
      if tmp == doors[j]   {
        doorStats[j] += 1
      }
    }
  }
  return doorStats
}


//TODO average escape time
//TODO average health impact
//TODO total time for people to get out
//TODO average time spent waiting
//TODO took most damage from smoke/fire


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
  readStats(pplArray, testmap)

}
