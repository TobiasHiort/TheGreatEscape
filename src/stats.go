package main

import (
  "fmt"
  "log"
  //"strconv"
  "encoding/json"
  "os"
)

func readStats(peopleArray []*Person, inmap [][]tile) {

  pplfile, err := os.Create("peopleStats.txt")
  if err != nil {
    log.Fatal("Cannot create file, ppl")
  }
  defer pplfile.Close()

  pplStats := PeopleStats(peopleArray)

  bytes2, err2 := json.Marshal(pplStats)
  if err2 != nil {
    panic(err2)
  }
  s := string(bytes2[:])
  //alive dead injured
  fmt.Fprintf(pplfile, s)


  mapfile, err2 := os.Create("mapStats.txt")
  if err2 != nil {
    log.Fatal("Cannot create file, map")
  }
  defer mapfile.Close()

  //burning tiles
  mapStats := MapStats(inmap)

  bytes2, err2 = json.Marshal(mapStats)
  if err2 != nil {
    panic(err2)
  }
  str := string(bytes2[:])
  fmt.Fprintf(mapfile, str)

  //mapfile.Close()
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
//TODO write average escapetime to file
//TODO call this func at end of simuation
func averageExitTime(peopleArray[] *Person) float32 {

  var totalTime float32
  size := len(peopleArray)
  for i := 0; i < size; i++ {
    totalTime += peopleArray[i].time 
  }
  if size != 0 {
    fmt.Print(totalTime/float32(size))
    return (totalTime/float32(size))
  }else {return 0}
}


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
  //averageExitTime(pplArray)

}
