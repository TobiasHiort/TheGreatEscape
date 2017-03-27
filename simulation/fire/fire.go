package main

import "fmt"

// print readable matrix
func printMatrix(matrix [][]int, num int) {
	fmt.Printf("\n Matrix %v :", num)
	fmt.Printf("\n %v \n %v \n %v \n %v \n %v \n ", matrix[0], matrix[1], matrix[2], matrix[3], matrix[4])
}

// fire spread algorithm
func fire(matrix [][]int, x int, y int) {

	matrix[x][y] += 3 // [row, column], 3 = full fire from start
	printMatrix(matrix, 0) // print first fire

	for i := 0; i < 3; i++ {
		fireAddAdjacent(matrix, x, y)
		x += 1
		y += 1
		fireAddAdjacent(matrix, x, y)

		printMatrix(matrix, i+1)
	}

}

func fireAddAdjacent(matrix [][]int, x int, y int) {
	
	if (x+1) < len(matrix) && (y+1) < len(matrix[0]) {
		matrix[x+1][y] += 1
		matrix[x][y+1] += 1
		matrix[x-1][y] += 1
		matrix[x][y-1] += 1	
	}
}

func main() {

	// predefined size of matrix
	startMatrix := [][] int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0}}
	
	// run fire algorithm
	fire(startMatrix, 2, 3) // [row, column]

	//fmt.Print(len(startMatrix))
	fmt.Print(len(startMatrix[0])) // column
	fmt.Print(len(startMatrix)) // rows
    
}
