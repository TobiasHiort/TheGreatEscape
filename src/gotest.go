package main

import (
	//"fmt"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	//"time"
	"io/ioutil"
	//"log"
)

func toPipe(stats [][]int) {
	bytes2, err2 := json.Marshal(stats)
	if err2 != nil {
		panic(err2)
	}
	s := string(bytes2[:])
	fmt.Println(s)
}


func mains() {
	mainMap()
}

func main() {

    b, err3 := ioutil.ReadFile("../src/mapfile.txt")
    if err3 != nil{
        panic(err3)
    }

    var m = [][]int{}
    err := json.Unmarshal(b, &m)
    if err != nil{
        panic(err)
    }
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

	ppl := PeopleInit(testmap, mm)
	InitPlans(&testmap)


	//fire := GetTile(testmap, 20, 20)
	
//	stats := [][]int{}
	stats := StartStats(ppl)
	SetFire(GetTile(testmap, 20, 20))
	fireStats := FireStats(&testmap) 

	//	SetFire(GetTile(testmap, 2, 2))
	

	for !CheckFinish(ppl) {
		toPipe(stats)
		toPipe(fireStats)//FireStats(&testmap))
		fireStats = FireStats(&testmap)//2(fire) //fire.getFS()
		//time.Sleep(10 * time.Millisecond)
		
		Run(&testmap, ppl, &stats)
		/*
			for check == 0 {
				check = fromPipe()
			} */
	}

	//	go func() {
	//		SingleSimulation(m, ppl)
	//	}()
	//	if timeToSend {toPipe(stats)}
	//	FireSpread2(fireStats)
	}

}

/**
func fromPipe() int{
	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()

	m := 0

	err := json.Unmarshal(line, &m)
	if err != nil {
		panic(err)
	}

	return m
}*/
