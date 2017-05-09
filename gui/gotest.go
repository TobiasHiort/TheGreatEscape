package main

import (
    //"fmt"
    "os"
    "bufio"
    //"fmt"
    "encoding/json"
    //"log"
)

func assembleMap(pipeAmount int) /*[][]int*/ {
	bio2 := bufio.NewReader(os.Stdin)
	var tmp []int
	for i := 0; i < pipeAmount; i ++ {
		//send hello to py
		//handshake()
		//recievee stuff from py
		if bio2.Peek(1) != !nil {
			//yay
			//shame overwhelms my senses



			piece, _, _ := bio2.ReadLine()
			//append to blah
			err3 := json.Unmarshal(piece, &tmp)

			if err3 != nil {
				panic(err3)
			}
		}

	}
	//return aMap
}

/*
func handshake() {
	sayhi := true
	bytes3, _ := json.Marshal(sayhi)
	histr := string(bytes3[:])
	fmt.Println(histr)
}
*/

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
		//handshake()

		//bio = bufio.NewReader(os.Stdin)
		//line, _, _ := bio.ReadLine()
		//		arrsiz := bio.ReadLine()


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
