package main

import (
  "fmt"
  "log"
  "strconv"
  "os"
)

func readStats(peopleArray []*Person, inmap [][]tile) {

  pplfile, err := os.Create("peopleStats.txt")
  if err != nil {
    log.Fatal("Cannot create file, ppl")
  }
  defer pplfile.Close()


  pplStats := PeopleStats(peopleArray)
  //alive dead injured
	fmt.Fprintf(pplfile, json.Marshal(pplStats))

  mapfile, err2 := os.Create("mapStats.txt")
  if err2 != nil {
    log.Fatal("Cannot create file, map")
  }
  defer mapfile.Close()

	//fmt.Fprintf(file, "YAY")
	//burning tiles
	mapStats := MapStats(inmap)
	fmt.Fprintf(mapfile, json.Marshal(mapStats))
	//mapfile.Close()

	//exit statserino


}

//TODO most used door
func exitStats(peopleArray []*Person, inmap [][]tile) []int {

  doors := doorCoord(inmap)
  var doorStats []int
  index := (len(peopleArray) - 2)

  for  i := 0; i < (len(peopleArray)); i++ {
    tmp := []int {peopleArray[i].path[index].xCoord, peopleArray[i].path[index].yCoord}
    for j := 0; j < (len(doors)); j++ {
      if tmp[0] == doors[j][0] && tmp[1] == doors[j][1] {
        doorStats[j] ++
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
