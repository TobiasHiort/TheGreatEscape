package main

import (
	"fmt"
	"log"
	//"strconv"
	"encoding/json"
	"os"
)

func readStats(peopleArray []*Person, inmap [][]tile) {

	pplfile, err := os.Create("peopleStats.txt") //alive deaths injured
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


	mapfile, err2 := os.Create("mapStats.txt") //how many tiles are on fire?, most used exit
	if err2 != nil {
		log.Fatal("Cannot create file, map")
	}
	defer mapfile.Close()

	//burning tiles
	burningTiles := MapStats(inmap)
	//fmt.Print(mapStats)

	bytes2, err2 = json.Marshal(burningTiles)
	if err2 != nil {
		panic(err2)
	}
	str := string(bytes2[:])
	fmt.Fprintf(mapfile, str)

	/*
	most used exitsdoors
	mostUsedDoors := doorStats(peopleArray, inmap)

	fmt.Print(mostUsedDoors)
	bytes2, err2 = json.Marshal(mostUsedDoors)
	if err2 != nil {
		panic(err2)
	}
	*/
	str2 := string(bytes2[:])
	fmt.Fprintf(mapfile, str2)
	fmt.Print(str2)

	//mapfile.Close()
}
/*
//TODO most used exit door
//TODO dis function doesnt wörk, no exit is found
func doorStats(peopleArray []*Person, inmap [][]tile) []int {

	doors := DoorCoord(inmap)
	//debug
	fmt.Print(doors)
	var statistics []int
	index := (len(peopleArray) - 2) //vi behöver näst sista kordinaten

	var tmpx int
	var tmpy int
	for  i := 0; i < (len(peopleArray)); i++ {
		tmpx = peopleArray[i].path[index].xCoord
		tmpy = peopleArray[i].path[index].yCoord
		for j := 0; j < (len(doors)); j++ {
			if tmpx == doors[j][0] && tmpy == doors[j][1] {
				doorStats[j] ++
			}
		}
	}
	return statistics
}
*/


//TODO average escape time
//TODO write average escapetime to file
//TODO call this func at end of simuation
func averageExitTime(peopleArray []*Person) float32 {

	var totalTime float32

	size := len(peopleArray)
	for _, p := range peopleArray {
		totalTime = totalTime + p.time
	}
	if size != 0 {
		// fmt.Print(totalTime/float32(size))
		return (totalTime/float32(size))
	}else {return 0}
}


//TODO average health impact
//TODO write average healthexit to file
func averageExitHealth(peopleArray []*Person) int {

	var totalHealth int

	size := len(peopleArray)
	for _, p := range peopleArray {
		totalHealth = totalHealth + p.hp
	}
	if size != 0 {
		fmt.Print(totalHealth/size)
		return (totalHealth/size)
	}else {return 0}
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
		{1, 0, 0, 1, 0, 0, 0},
		{1, 0, 0, 1, 0, 0, 0},
		{3, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 2}}

		ppl := [][]int{
			{0, 0},
			{0, 6},
			{2, 4}}

			testmap := TileConvert(matrix)
			pplArray:= PeopleInit(testmap,ppl)

			SetFire(&testmap[1][1])
			for i := 0 ; i < 100 ; i++ {
				Run(&testmap, pplArray, &ppl)
			}
			readStats(pplArray, testmap)
			averageExitTime(pplArray)

			for i, p := range pplArray {
				fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
			}

			doorstat := doorStats(pplArray, testmap)
			for i := 0; i < len(doorstat); i++ {
				fmt.Print(doorstat[i])
			}
}
*/
