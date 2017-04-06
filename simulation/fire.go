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
	printMatrix(matrix, 0)

	for n := 0; n < 8; n++ { // n = number of time units
		matrix = fireChecker(matrix)
		printMatrix(matrix, n+1)	
	}
}

// check adjacent tiles in matrix for fire status
func fireChecker(matrix [][]int) [][]int {
// return = (...) _[][]int_ {}
	for row := 0; row < len(matrix); row++ {
		for column:= 0 ; column < len(matrix[0]); column++ {
			if matrix[row][column] == 3 {
				if row+1 < len(matrix) && row+1 >= 0 {
					if matrix[row+1][column] != 3 {
						matrix[row+1][column] += 1
					}
				}
				if column+1 < len(matrix[0]) && column+1 >= 0 {
					if matrix[row][column+1] != 3 {
						matrix[row][column+1] += 1
					}
				}
				if row-1 < len(matrix) && row-1 >= 0 {
					if matrix[row-1][column] != 3 {
						matrix[row-1][column] += 1
					}
				}
				if column-1 < len(matrix[0]) && column-1 >= 0 {
					if matrix[row][column-1] != 3 {
						matrix[row][column-1] += 1
					}
				}
			}
		}
	}
	return matrix
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
	fire(startMatrix, 2, 2) // [row, column] = start pos

	//fmt.Print(len(startMatrix))
	fmt.Print(len(startMatrix[0])) // column
	fmt.Print(len(startMatrix)) // rows
    
}