package main

import (
	"fmt"
	"log"
	//"strconv"
	"encoding/json"
	"os"
	"math"
)

func readStats(peopleArray []*Person, inmap [][]tile) {
	
	// or create
  pplfile, err := os.Create("peopleStats.txt") //[alive, deaths, injured],[average exit time], [average health], [died from smoke, died from fire]
  if err != nil {
    log.Fatal("Cannot create file, ppl")
  }
	defer pplfile.Close()
//	toPpl := [][]float32{}
	toPpl := [][]float32{PeopleStats(peopleArray)}
//	smoke, fire := smokeVSFireDmg(peopleArray)
	toPpl = append(toPpl, []float32{averageExitTime(peopleArray)}, []float32{averageExitHealth(peopleArray)}, smokeVSFireDmg(peopleArray))
	// 
/*	
  pplStats := PeopleStats(peopleArray)
  bytes2, err2 := json.Marshal(pplStats)
  if err2 != nil {
    panic(err2)
  }
  s := string(bytes2[:])
  //alive dead injured
 // fmt.Fprintf(pplfile, s)
	toPpl = append(toPpl, s)
  exitTime := []float32{averageExitTime(peopleArray)} 
  
  bytes2, err2 = json.Marshal(exitTime)
  if err2 != nil {
    panic(err2)
  }
  s = string(bytes2[:])
  //average exit time
 // fmt.Fprintf(pplfile, s)
	toPpl = append(toPpl, s)
	
  exitHealth := []int{averageExitHealth(peopleArray)} 
  
  bytes2, err2 = json.Marshal(exitHealth)
  if err2 != nil {
    panic(err2)
  }
  s = string(bytes2[:])
  //average exit health
 // fmt.Fprintf(pplfile, s)
	toPpl = append(toPpl, s)
*/	
  bytes2, err2 := json.Marshal(toPpl)
  if err2 != nil {
    panic(err2)
  }
  s := string(bytes2[:])
	
	fmt.Fprintf(pplfile, s)

// mapstats
  mapfile, err2 := os.Create("mapStats.txt")
  //[how many tiles are on fire?], [amount exited per door], [exit-doorcordinates]
  if err2 != nil {
    log.Fatal("Cannot create file, map")
  }
  defer mapfile.Close()

	toMap := [][][]int{}
	ds1, ds2 := doorStats(peopleArray, inmap)
	toMap = append(toMap, [][]int{MapStats(inmap)}, [][]int{ds1}, ds2) 

	bytes2, err2 = json.Marshal(toMap)
	if err2 != nil {
		panic(err2)
	}
	s = string(bytes2[:])
	fmt.Fprintf(mapfile, s)
/*	
  //burning tiles
  burningTiles := MapStats(inmap)

  bytes2, err2 = json.Marshal(burningTiles)
  if err2 != nil {
    panic(err2)
  }
  str := string(bytes2[:])
  fmt.Fprintf(mapfile, str)

  //most used exitsdoors
  noUsed, doorCoord := doorStats(peopleArray, inmap)

  bytes2, err2 = json.Marshal(noUsed)
  if err2 != nil {
    panic(err2)
  }

  str2 := string(bytes2[:])
 // fmt.Println(str2)
  fmt.Fprintf(mapfile, str2)
  //fmt.Print(str2)

  bytes2, err2 = json.Marshal(doorCoord)
  if err2 != nil {
    panic(err2)
  }

  str3 := string(bytes2[:])
 // fmt.Println(str3)
  fmt.Fprintf(mapfile, str3)
 // fmt.Print(str3)
  //mapfile.Close() */

	// escaped/died times
  timefile, err2 := os.Create("timeStats.txt")
  //[how many tiles are on fire?], [amount exited per door], [exit-doorcordinates]
  if err2 != nil {
    log.Fatal("Cannot create file, map")
  }
  defer timefile.Close()

//	toTime := [][]int{}
	escaped, died := exitTimes(peopleArray)
	toTime := [][]int{escaped, died}

	bytes2, err2 = json.Marshal(toTime)
	if err2 != nil {
		panic(err2)
	}
	s = string(bytes2[:])
	fmt.Fprintf(timefile, s)	
	
}

//TODO most used exit door
func doorStats(peopleArray []*Person, inmap [][]tile) ([]int, [][]int) {

  doors := DoorCoord(inmap)
  numberOfExits := make([]int, len(doors))

  var tmpx int
  var tmpy int
  for  i := 0; i < (len(peopleArray)); i++ {

    index := (len(peopleArray[i].path) - 2) //vi behöver näst sista kordinaten
    tmpx = peopleArray[i].path[index].xCoord
    tmpy = peopleArray[i].path[index].yCoord
    for j := 0; j < (len(doors)); j++ {
      if tmpx == doors[j][0] && tmpy == doors[j][1] {
        numberOfExits[j] = numberOfExits[j] + 1
      }
    }
  }
  return numberOfExits, doors
}

//TODO average escape time
func averageExitTime(peopleArray []*Person) float32 {

  var totalTime float32

  size := len(peopleArray)
  for i, p := range peopleArray {
    if (peopleArray[i].alive == true) {
      totalTime = totalTime + p.time
    }
  }
  if size != 0 {
    return (totalTime/float32(size))
  }else {return 0}
}


//TODO average health impact
func averageExitHealth(peopleArray []*Person) float32 {

  var totalHealth int
  var alive int

  for i, _ := range peopleArray {
     if(peopleArray[i].alive == true) {
      alive ++
    }
  }
  for i, p := range peopleArray {
    if(peopleArray[i].alive == true) {
      totalHealth = totalHealth + p.hp
    }
  }
  if alive != 0 {
	  return (float32(totalHealth/alive))
  }else {return 0}
}

func smokeVSFireDmg(peopleArray []*Person) []float32{
	smoke := float32(0)
	fire := float32(0)
	for _, p := range peopleArray {
		if !p.alive {
			if p.smokedmg > (100 - p.hp)/2 {
				smoke++} else {fire++}
		}
	}
	return []float32{smoke, fire}
}

func exitTimes(peopleArray []*Person) ([]int, []int){  // finishedtimes?? I dunno
	var escaped = make([]int, int(step+10))
	var died = make([]int, int(step+10))  
	for _, p := range peopleArray {          //Obs not 100% on the indexing here, in case of bug, break the glass.
		//	if p.time > step {panic(fmt.Sprintf("step: %v, time: %v", step, p.time ))}
		ind := int(math.Floor(float64(p.time))) 
		if ind > len(escaped) {
			for i := len(escaped); i <= ind; i++ {
				escaped = append(escaped, 0)
				died = append(died, 0)
			}
		}
		if p.alive {
			escaped[ind]++ //= escaped[int(math.Floor(float64(p.time)))] + 1
		} else {
			died[ind]++} // = died[int(math.Floor(float64(p.time)))] + 1 }
	}
	return escaped, died
}

//TODO death y [....] per time x [....] in file
//TODO average time spent waiting
//TODO took most damage from smoke/fire
//relevant? on individual lvl? 
/*
func main() {

  matrix := [][]int{
    {1, 1, 3, 1, 1, 1, 1},
    {1, 0, 0, 0, 0, 0, 0},
    {1, 0, 0, 1, 0, 0, 0},
    {1, 0, 0, 1, 0, 0, 2},
    {1, 0, 0, 1, 0, 0, 0},
    {3, 0, 0, 0, 0, 0, 0},
    {1, 0, 0, 0, 0, 0, 0}}

    ppl := [][]int{
      {1, 1},
      {0, 6},
      {0, 5},
      {2, 4}}

      testmap := TileConvert(matrix)
      pplArray:= PeopleInit(testmap,ppl)

      SetFire(&testmap[3][4])
	(&testmap[3][4]).smoke = 40
	//SetFire(&testmap[3][6])
     for i := 0 ; i < 7 ; i++ {
//	for CheckFinish(pplArray){
	Run(&testmap, pplArray, &ppl)
      }
      readStats(pplArray, testmap)
	averageExitTime(pplArray)

	smoke, fire := smokeVSFireDmg(pplArray)
	println("Smoke:", smoke)
	println("Fire:", fire)

      for i, p := range pplArray {
        fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
      }
    }
*/
