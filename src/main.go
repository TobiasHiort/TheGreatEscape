package main

/*

-- Paste Python communication stuff here

*/

//[][]int, peopleList []*Person, fireStartPos []int

//TODO func for input from pipe
//TODO runtime single simulation
//TODO runtime multiple simulation
//TODO firestuff


func sendToPipe(exitStatus *int, statsList[][] int, PIPEIN, a *int, b *int) {

  for *exitStatus == 0 {
   if !(a == b) {
    //TODO Copy list to pipe
   *b++
   }

  }

}



func singleSimulation(mapList [][]int, peopleList [][]int, fireStartPos[]int) {

  //TODO: create lsit for positions
  //TODO: implement spinlock in gameloop
  a := 0 //pointers?
  b := 0
  exitStatus := 0
  //TODO: create function to copy list and send to python through pipe
  //TODO: implenet sem lock + spinlock t ensure wait for all people to move
  //TODO: implement that both gameloop and copy func tries to run concurrently, spinlock continously spins

  size := len(peopleList)
  statsList := make([][]int, size)
  for i := range statsList {
    statsList[i] = make([]int, 3)
  }

  GameLoop(mapList, peopleList, fireStartPos, statsList, &a, &b, &exitStatus)

  

}
