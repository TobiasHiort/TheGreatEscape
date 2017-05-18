package main

import (
	//"os"
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"math"
)

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
	toPipe(posList)
	toPipe(fireList)
}

func fromPipe() ([][]int, [][]int, []int) {
	//TODO: Get fire start position*/
	b, err3 := ioutil.ReadFile("../src/mapfile.txt")
	if err3 != nil {
		panic(err3)
	}

	var m = [][]int{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	c, err4 := ioutil.ReadFile("../src/playerfile.txt")
	if err4 != nil {
		panic(err4)
	}

	var mm = [][]int{}
	err5 := json.Unmarshal(c, &mm)
	if err5 != nil {
		panic(err5)
	}
	d, err6 := ioutil.ReadFile("../src/firefile.txt")
	if err6 != nil {
		panic(err6)
	}

	var mmm = []int{}

	err7 := json.Unmarshal(d, &mmm)
	if err7 != nil {
		panic(err7)
	}

	if len(mmm) < 2 {
		//We are gonna do something drastic here!
	} else if len(mmm)%2 != 0 {
		mmm = mmm[:len(mmm)-1]
	}
	return m, mm, mmm
}

func singleSimulation(fire [2]int) {
	mapList, peopleList, fireList := fromPipe()

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
	//posList := StartStats(peopleList)
	//fireList := StatsInit(10)
	GameLoop(mapList, peopleList, fireList)
	//sendToPipe(&exitStatus, &posList, fireList, &a, &b)
}

func main() {
	var fire [2]int
	fire[0] = 1
	fire[1] = 1

	singleSimulation(fire)
}
