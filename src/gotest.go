package main

import (
    //"fmt"
    "os"
    "bufio"
    "fmt"
	"encoding/json"
	//"time"
    "io/ioutil"
    //"log"
)

func main1() {
	//fmt.Println("GO STARTED2")

    bio := bufio.NewReader(os.Stdin)
    line, _, _ := bio.ReadLine()

    var p = [][]float32{}

    err := json.Unmarshal(line, &p)
    if err != nil {
        panic(err)
    }

    bytes2, err2 := json.Marshal(p)
    if err2 != nil {
        panic(err2)
    }
    s := string(bytes2[:])

    fmt.Println(s)

}

func toPipe(stats [][]int) {
	bytes2, err2 := json.Marshal(stats)
	if err2 != nil {
		panic(err2)
	}
	s := string(bytes2[:])
	fmt.Println(s)
}


func main11() {
	mainMap()
}

func main() {
    /*
	bio := bufio.NewReader(os.Stdin)

    line, _, _ := bio.ReadLine()

	if line == nil {}


    var m = [][]int{}

    err := json.Unmarshal(line, &m)
    if err != nil {
        panic(err)
    }
    */
	//m[8][1] = 2
	//m[13][0] = 2


    b, err3 := ioutil.ReadFile("../src/mapfile.txt")
    if err3 != nil{
        panic(err3)
    }

    var m = [][]int{}
    err := json.Unmarshal(b, &m)
    if err != nil{
        panic(err)
    }
    //m[8][1] = 2
    //m[13][0] = 2

	testmap := TileConvert(m)
	if testmap == nil {}

    c, err4 := ioutil.ReadFile("../src/playerfile.txt")
    if err4 != nil{
        panic(err4)
    }

    var mm = [][]int{}
    err5 := json.Unmarshal(c, &mm)
    if err5 != nil{
        panic(err5)
    }

    /*
	list := [][]int{
		{1,1},
		{5,1},
        {7,1},
        {8,1},
        {9,2},
        {10,3},
        {11,50},
		{22,90},
        {25,105},
        {25,125}}
		//{3,3}}
    */

	ppl := PeopleInit(testmap, mm)
	InitPlans(&testmap)


//	stats := [][]int{}
	stats := StartStats(ppl)
//	toPipe(stats)
	//fmt.Println(stats)
	//	Run(&testmap, ppl, &stats) // startstats!
	//fmt.Println(len(stats))



	for !CheckFinish(ppl) {
		toPipe(stats)
	
		//time.Sleep(10 * time.Millisecond)
		Run(&testmap, ppl, &stats)

		/*ok := true
		for _, pers := range ppl {
			if !pers.safe && !pers.wasDiag() && !pers.IsWaiting() {ok = false}
		}*/

		/*
		for check == 0 {
			check = fromPipe()
		} */
	}


	//	go func() {
	//		SingleSimulation(m, ppl)
	//	}()
	//	if timeToSend {toPipe(stats)}
}

func fromPipe() int{
	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()

	m := 0

	err := json.Unmarshal(line, &m)
	if err != nil {
		panic(err)
	}

	return m
}
