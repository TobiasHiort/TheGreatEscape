package main

import (
	//"os"
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

			//TODO Copy list to pipe
	posCopy := posList
	fireCopy := fireList
		
	toPipe(posCopy)
	toPipe(fireCopy)
}

func fromPipe() ([][]int, [][]int) {
	//TODO: Get fire start position
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

	GameLoop(mapList, peopleList, fireStartPos)
	//sendToPipe(&exitStatus, &posList, fireList, &a, &b)
}

func main() {
	var fireStartPos [2]int
	fireStartPos[0] = 1
	fireStartPos[1] = 1
	singleSimulation(fireStartPos)
}
