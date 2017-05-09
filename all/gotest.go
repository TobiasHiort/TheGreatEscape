package main

import (
    //"fmt"
    "os"
    "bufio"
    "fmt"
	"encoding/json"
	"time"
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


func main() {
	//fmt.Println("GO STARTED2")

	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()

	if line == nil {}
		

    var m = [][]int{}

    err := json.Unmarshal(line, &m)
    if err != nil {
        panic(err)
    }

//	m[8][1] = 2
//	m[13][0] = 2

  
	testmap := TileConvert(m)
	if testmap == nil {}

	list := [][]int{
		{1,1},
		{1,3},
    {1,8},
    {3,3}}

	ppl := PeopleInit(testmap, list)


	stats := [][]int{}
  Run(&testmap, ppl, &stats) // startstats!
	//fmt.Println(len(stats))

	//	check := 0
	
	for !CheckFinish(ppl) {
		toPipe(stats)
		time.Sleep(10 * time.Millisecond)
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
