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
	"math"
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
	//InitPlans(&testmap)
	//plans := InitPlans2(&testmap)
	InitPlans(&testmap)


	//fire := GetTile(testmap, 20, 20)
	
//	stats := [][]int{}
	stats := StartStats(ppl)
	SetFire(GetTile(testmap, 31, 31))
	fireStats := FireStats(&testmap)
	smokeStats := SmokeStats(&testmap)

	//	SetFire(GetTile(testmap, 2, 2))
	

	for !CheckFinish(ppl) {
		toPipe(stats)
		toPipe(fireStats)//FireStats(&testmap))
		toPipe(smokeStats)

		if math.Mod(float64(step), 2) == 0 {
			fireStats = FireStats(&testmap)//2(fire) //fire.getFS()
			smokeStats = SmokeStats(&testmap)}
		//time.Sleep(10 * time.Millisecond)
		
		//UpdateParentOf(&testmap, plans, fireStats)//[]*tile{&(testmap)[20][20]})
		
		Run(&testmap, ppl, &stats)
	//	FireSpread2(fireStats)
	}
	toPipe(stats)
	toPipe(fireStats)//FireStats(&testmap))

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
