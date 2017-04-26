package main

//toppen av stupröret, yttersta lagret av löken

//time := int
//time = 0

type event struct {
	tileMap [][]tile
	timeStamp int
}

func TimeInit(){
	time := *int
	*time = 0
}

//func get ppl movement from map

func timeStamp(inMap [][]tile, currentTime int) event{
	ev := event
	ev -> tileMap = inMap
	ev -> timeStamp = currentTime
	return ev
}

//this function will in the future preferably grab just the changes on the map, rather than it as a whole
	//outEvent := event
	//outEvent = timeStamp(tileMap, time)
	//grabs data from map
	//returns an event
func Tick(){
	*time++
}

func MapInit(peopleMap [][]int, newMap [][]int){
	//gets map data from GM and inits the map
	tileMap := *[][]tile
	*tileMap = TileConvert(newMap)
	//Populate(newMap)
	//gets ppl data from GM and calls the initpppl function in map 
}

/*
func fetch(inMap [][]tile){
}
*/


func ClockCycle() {
	// should be a goroutine

	// save timeStamp((fetch(*tileMap)), *time) to mem
}
