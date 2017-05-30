package main

import (
	//"os"
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	
	//"os"
	//"runtime/trace"
	//	"math"
)

//TODO func for input from pipe
//TODO runtime single simulation
//TODO runtime multiple simulation
//TODO firestuff

func toPipe(list *[][]int) {
	bytes, err := json.Marshal(*list)
	if err != nil {
		panic(err)
	}
	s := string(bytes[:])
	fmt.Println(s)
}

func SendToPipe(posList *[][]int, fireList *[][]int, smokeList *[][]int) {
	toPipe(posList)
	toPipe(fireList)
	toPipe(smokeList)

	if posList == nil || fireList == nil || smokeList == nil {panic("whyy?")}
}

func fromPipe() ([][]int, [][]int, [][]int, []float64) {
	b, err3 := ioutil.ReadFile("../tmp/mapfile.txt")
	if err3 != nil {
		panic(err3)
	}

	var m = [][]int{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	c, err4 := ioutil.ReadFile("../tmp/playerfile.txt")
	if err4 != nil {
		panic(err4)
	}

	var mm = [][]int{}
	err5 := json.Unmarshal(c, &mm)
	if err5 != nil {
		panic(err5)
	}
	d, err6 := ioutil.ReadFile("../tmp/firefile.txt")
	if err6 != nil {
		panic(err6)
	}

	var mmm = [][]int{}

	err7 := json.Unmarshal(d, &mmm)
	if err7 != nil {
		panic(err7)
	}
	
	e, err8 := ioutil.ReadFile("../tmp/velocitiesfile.txt")
	if err8 != nil {
		panic(err8)
	}

	var mmmm = []float64{}
	err9 := json.Unmarshal(e, &mmmm)
	if err9 != nil {
		panic(err9)
	}
	
	return m, mm, mmm, mmmm  /// ... Marabou
}

func singleSimulation() {
	mapList, peopleList, fireList, velocities := fromPipe()
//	toPipe(&mapList)
//	toPipe(&mapList)
	//TODO: implement spinlock in gameloop

	//TODO: implenet sem lock + spinlock t ensure wait for all people to move
	//TODO: implement that both gameloop and copy func tries to run concurrently, spinlock continously spins
//	fs := float64(2)
//	ps := float64(2)
	GameLoop(mapList, peopleList, fireList, velocities)
}

func main() {
/*	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()*/
	//	mainMap()
	singleSimulation()
}
