package main

import "testing"

func makeTestQueue(size int) *queue {
	testQ := queue{}

	for i := 0; i < size; i++ {
		tile := makeNewTile(0, i, i)
		testQ.Add(&tile, float32(i))
	}
	return &testQ
}

func TestAdd(t *testing.T) {
	size := 10
	testQ := queue{}
	i := 0
	for ; i < size; i++ {
		tile := makeNewTile(0, i, i)
		testQ.Add(&tile, float32(i))
	}
	if testQ.Len() != i {
		t.Errorf("Expected len: %d, but got len: %d", i, testQ.Len())
	}

	j := 0
	for ; j < size*2; j++ {
		tile := makeNewTile(0, j, j)
		testQ.Add(&tile, float32(j))
	}
	if testQ.Len() != j {
		t.Errorf("Expected len: %d, but got len: %d", j, testQ.Len())
	}

}

func TestRemove(t *testing.T) {
	size := 10

	testQ := *makeTestQueue(size)
	checkQ := *makeTestQueue(size * 2)

	// remove all elements, and then some
	for i := 0; i < size*2; i++ {
		testQ.Remove(checkQ[i].tile)
	}
	if testQ.Len() != 0 {
		t.Errorf("Expected len: %d, but got len: %d", 0, testQ.Len())
	}

	testQ = *makeTestQueue(size)

	// remove elements that are not in the queue
	for i := size; i < size*2; i++ {
		testQ.Remove(checkQ[i].tile)
	}
	if testQ.Len() != size {
		t.Errorf("Expected len: %d, but got len: %d", size, testQ.Len())
	}
}

func TestInQueue(t *testing.T) {
	size := 10 //... time to make this shit global or what?
	testQ := *makeTestQueue(size)

	for i := 0; i < size*2; i++ {
		//	if i < size && !testQ.inQueue(testQ[i].tile) {
		// stooopid golang räknar saker i onödan... I MISS HASKELL!
		if i < size {
			if !testQ.inQueue(testQ[i].tile) {
				t.Errorf("Expected tile: %d to be in queue", testQ[i].tile)
			} else if i >= size && testQ.inQueue(testQ[i].tile) {
				t.Errorf("Expected tile: %d to not  be in queue", testQ[i].tile)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	size := 10
	newVal := float32(1)

	testQ := *makeTestQueue(size)
	checkQ := *makeTestQueue(size * 2)

	// update entire queue tocost= newVal
	for i := 0; i < size; i++ {
		testQ.Update(checkQ[i].tile, newVal)
	}

	for _, tc := range testQ {
		if *tc.cost != newVal {
			t.Errorf("Expected cost: %f, but got cost: %f", newVal, *tc.cost)
		}
		//jenny kan inte organisera paranteser korrekt.... time wasted

	}
	//updating tiles not in queue+updateíng tiles to their current cost again
	for i := 0; i > size*2; i++ {
		testQ.Update(checkQ[i].tile, newVal)
	}
	if testQ.Len() != size {
		t.Errorf("Expected len: %d, but got len: %d, update changed length of queue", size, testQ.Len())
	}
	for _, tc := range testQ {
		if *tc.cost != newVal {
			t.Errorf("Expected cost: %f, but got cost: %f", newVal, *tc.cost)
		}
		//ghaaaa writing tests are boring....
	}
}

func TestPop(t *testing.T) {
	size := 10

	testQ := *makeTestQueue(size)
	// pops all elements and then some
	for i := 0; i < size*2; i++ {
		current := testQ.Pop()	
		if i < size {
			if *current.cost != float32(i) { // checks that popped tile has correct value
				t.Errorf("Expected cost: %d but got cost: %d", i, current.cost)
			}
		} else if current.cost != nil && *current.cost != 0 { // checks that popped tile is 'nil', zero cost
			t.Errorf("Expected cost: %d but got cost: %d", 0, current.cost)
		}
	}
	// queue should be empty after all elements have been popped
	if testQ.Len() != 0 {
		t.Errorf("Expected empty queue but got len: %d", testQ.Len())
	}

}
