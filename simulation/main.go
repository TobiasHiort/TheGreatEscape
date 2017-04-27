package main

/*

-- Paste Python communication stuff here

*/

func GameLoop() {
	//newMap := MapInit(foo, bar)
	//do all the Inits

	//simulates 5000 steps
	for i := 0; i < 5000; i++ {
		go Tick()
		go Discretize()
	}
}


