package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*

-- Paste Python communication stuff here

*/

//[][]int, peopleList []*Person, fireStartPos []int

//TODO func for input from pipe
//TODO runtime single simulation
//TODO runtime multiple simulation
//TODO firestuff

func toPipe(list [][]int) {
	bytes, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	s := string(bytes[:])
	fmt.Println(s)

}

func sendToPipe(exitStatus *int, posList [][]int, fireList [][]int, a *int, b *int) {

	for *exitStatus == 0 {
		if !(a == b) {
			//TODO Copy list to pipe
			posCopy := posList
			fireCopy := fireList
			*b++
			toPipe(posCopy)
			toPipe(fireCopy)
		}

	}

}

func fromPipe() ([][]int, [][]int) {
	//Import map
	b, err := ioutil.ReadFile("../src/mapfile.txt")
	if err != nil {
		panic(err)
	}
	var m = [][]int{}
	err2 := json.Unmarshal(b, &m)
	if err2 != nil {
		panic(err2)
	}

	//Import people
	c, err3 := ioutil.ReadFile("../src/playerfile.txt")
	if err3 != nil {
		panic(err3)
	}
	var p = [][]int{}
	err4 := json.Unmarshal(c, &p)
	if err4 != nil {
		panic(err4)
	}
	//TODO: Get fire start position
	return m, p
}

func singleSimulation(fireStartPos [][]int) {
	mapList, peopleList := fromPipe()
	//TODO: create lsit for positions
	//TODO: implement spinlock in gameloop
	a := 0 //pointers?
	b := 0
	exitStatus := 0
	//TODO: create function to copy list and send to python through pipe
	//TODO: implenet sem lock + spinlock t ensure wait for all people to move
	//TODO: implement that both gameloop and copy func tries to run concurrently, spinlock continously spins
	/**
	  size := len(peopleList)
	  statsList := make([][]int, size)
	  for i := range statsList {
	    statsList[i] = make([]int, 3)
	  }*/
	posList := StatsInit(len(peopleList))
	fireList := StatsInit(10)
	//

	go GameLoop(mapList, peopleList, fireStartPos, posList, &a, &b, &exitStatus)
	go sendToPipe(&exitStatus, posList, fireList, &a, &b)
	//
}

func main() {
	fmt.Println("Bye World!!")
}
