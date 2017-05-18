package main

import (
	//"os"
	//"bufio"
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

func SendToPipe(posList [][]int, fireList [][]int) {

	//for *exitStatus == 0 {
		//if !(a == b) {
			//TODO Copy list to pipe
		
			posCopy := posList
			fireCopy := fireList
		
			toPipe(posCopy)
			toPipe(fireCopy)
		//}

	//}

}

func fromPipe() ([][]int, [][]int) {
	//Import map
	/**b, err := ioutil.ReadFile("../src/mapfile.txt")
	if err != nil {
		panic(err)
	}
	var m = [][]int{}
	err2 := json.Unmarshal(b, &m)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("MAP LIST DONE")

	//Import people
	c, err3 := ioutil.ReadFile("../src/playerfile.txt")
	if err3 != nil {
		panic(err3)
	}
	fmt.Println("OPEN PEOPLE LIST")
	var p = [][]int{}
	err4 := json.Unmarshal(c, &p)
	if err4 != nil {
		panic(err4)
	}
	fmt.Println("PEOPLE LIST COPIED")
	//TODO: Get fire start position*/
	b, err3 := ioutil.ReadFile("../src/mapfile.txt")
    if err3 != nil{
        panic(err3)
    }

    var m = [][]int{}
    err := json.Unmarshal(b, &m)
    if err != nil{
        panic(err)
    }


    c, err4 := ioutil.ReadFile("../src/playerfile.txt")
    if err4 != nil{
        panic(err4)
    }

    var mm = [][]int{}
    err5 := json.Unmarshal(c, &mm)
    if err5 != nil{
        panic(err5)
    }
	return m, mm
}

func singleSimulation(fireStartPos [2]int) {
	mapList, peopleList := fromPipe()
	//TODO: create lsit for positions
	//TODO: implement spinlock in gameloop

	//TODO: create function to copy list and send to python through pipe
	//TODO: implenet sem lock + spinlock t ensure wait for all people to move
	//TODO: implement that both gameloop and copy func tries to run concurrently, spinlock continously spins
	/**
	  size := len(peopleList)
	  statsList := make([][]int, size)
	  for i := range statsList {
	    statsList[i] = make([]int, 3)
	  }*/
	//posList := StatsInit(len(peopleList))
//	posList := StartStats(peopleList)
	//fireList := StatsInit(10)


	GameLoop(mapList, peopleList, fireStartPos)
	//sendToPipe(&exitStatus, &posList, fireList, &a, &b)
	//
}

func main() {
	//fmt.Println("Started main")
	var fireStartPos [2]int
	fireStartPos[0] = 1
	fireStartPos[1] = 1
	//fmt.Println("Fire pos started")
	singleSimulation(fireStartPos)
}
