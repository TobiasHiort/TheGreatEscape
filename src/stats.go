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

  stats := CompileStats(peopleArray);
  fmt.Fprintf(file, strconv.Itoa(stats[0]) + strconv.Itoa(stats[1]) + strconv.Itoa(stats[2]))

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
