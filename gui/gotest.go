package main

import (
    //"fmt"
    "os"
    "bufio"
    "fmt"
    "encoding/json"
    //"log"
)

func assembleMap(pipeAmount int) [][]int {
	bio2 := bufio.NewReader(os.Stdin)
	tmp := []int
	for i := 0; i < pipeAmount; i ++ {
		//recievee stuff from py
		piece, _, _ := bio.ReadLine()
		//append to blah
		err3 := json.Unmarshal(piece, &tmp1)
		//send hello to py
	}
}

func main() {
	//fmt.Println("GO STARTED2")

	//arrsiz := bufio.NewReader(os.Stdin)
	bio := bufio.NewReader(os.Stdin)

	arrsiz, _, _ := bio.ReadLine()
	var first int
	lolerr := json.Unmarshal(arrsiz, &first)
	//arrsiz2 := int(first)
	if lolerr != nil {
		panic(lolerr)
	}


	if first > 0 {

		//bio = bufio.NewReader(os.Stdin)
		//line, _, _ := bio.ReadLine()
		//		arrsiz := bio.ReadLine()

		sayhi := true
		bytes3, _ := json.Marshal(sayhi)
		histr := string(bytes3[:])
		fmt.Println(histr)

		//var p = [][]float32{}

		/*
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
		*/
	}

}
