package main

import (
	"bufio"
	"fmt"

	// "log"
	// "os"
	"strings"
)

var empty_elem = ' '

var in string = `3
1 1
2 1
5 3`

func main() {
	// file, err := os.Open("data/1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// input := bufio.NewReader(file)

	input := bufio.NewReader(strings.NewReader(in))

	var num_of_hex int

	fmt.Fscan(input, &num_of_hex)
	// fmt.Printf("num_of_hex = %d\n", num_of_hex)

	var height, width int
	for i := 0; i < num_of_hex; i++ {
		fmt.Fscan(input, &width, &height)

		matrix := CreateMatrix(width, height)

		PaintHex(matrix, 0, 0, width, height)
		PrintMatrix(matrix)
	}
	fmt.Println()
}

func CreateMatrix(width, height int) [][]rune {
	matrix := make([][]rune, 2*height+1)
	for j := 0; j < len(matrix); j++ {
		matrix[j] = make([]rune, width+2*height)
		for k := 0; k < len(matrix[j]); k++ {
			matrix[j][k] = empty_elem
		}
	}

	return matrix
}

func PrintMatrix(matrix [][]rune) {
	var bldr strings.Builder

	for i := 0; i < len(matrix); i++ {

		bldr.WriteString(strings.TrimRight(string(matrix[i]), " "))

		if i != len(matrix)-1 {
			bldr.WriteString("\n")
		}

	}

	fmt.Println(bldr.String())

}

func PaintHex(mtx [][]rune, x, y, w, h int) {

	// mtx - инициализированный двухмерный массив рун. В него рисуем фигуру.
	// x - x-координата фигуры на сетке "mtx"
	// y - y-координата фигуры на сетке "mtx"
	// w - width фигуры
	// h - hgeight фигуры

	for i := 0; i < len(mtx); i++ {
		for j := 0; j < len(mtx[i]); j++ {

			//рисуем верхнюю и нижнюю грани фигуры
			if (i == y) || (i == y+2*h) {
				if (j < x+(h+w)) && (j >= x+h) {
					mtx[i][j] = '_'
				}
			}

			// рисуем верхнюю-левую грань
			if (i - y) > 0 {
				if (i-y)+(j-x) == h {
					mtx[i][j] = '/'
				}
			}
			// рисуем верхнюю-правую грань
			if (i - y) > 0 {
				if (j - x) < h*2+w {
					if (j-x)-(i-y) == w+h-1 {
						mtx[i][j] = '\\'
					}
				}
			}

			// рисуем нижнюю-левую грань
			if (i - y) > 0 {
				if (j - x) < h*2+w {
					if (i-y)-(j-x) == h+1 {
						mtx[i][j] = '\\'
					}
				}
			}

			// рисуем нижнюю-правую грань
			if (i - y) > 0 {
				if (j - x) < h*2+w {
					if (i-y)+(j-x) == h*3+w {
						mtx[i][j] = '/'
					}
				}
			}
		}
	}

}
