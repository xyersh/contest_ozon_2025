package main

import (
	"fmt"
	"strings"
)

func main() {
	matrix := CreateMatrix(40, 10)
	// PaintHex(matrix, 5, 1, 4, 3)
	// PaintHex(matrix, 10, 1, 4, 3)

	MakeGrid(matrix, 2, 1, 15)

	PrintMatrix(matrix)

}

func CreateMatrix(mtx_width, mtx_height int) [][]rune {

	// mtx_width - горизонтальный размер поля (X)
	// mtx_height - вертикалоьный размер поля (Y)

	matrix := make([][]rune, mtx_height)
	for j := 0; j < mtx_height; j++ {
		matrix[j] = make([]rune, mtx_width)
		for k := 0; k < mtx_width; k++ {
			matrix[j][k] = ' '
		}
	}

	return matrix
}

func PrintMatrix(matrix [][]rune) {
	var bldr strings.Builder
	mtx_wdth := len(matrix[0])

	//верхняя граница
	bldr.WriteString("+" + strings.Repeat("-", mtx_wdth) + "+\n")

	for i := 0; i < len(matrix); i++ {

		bldr.WriteString("|" + string(matrix[i]) + "|")
		bldr.WriteString("\n")

	}

	// нижняя граница
	bldr.WriteString("+" + strings.Repeat("-", mtx_wdth) + "+\n")

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
			if (i - y) <= h {
				if (i - y) > 0 {
					if (i-y)+(j-x) == h {
						mtx[i][j] = '/'
					}
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
			if i-y <= 2*h {
				if (i - y) > h {
					if (j - x) < h*2+w {
						if (i-y)-(j-x) == h+1 {
							mtx[i][j] = '\\'
						}
					}
				}
			}

			// рисуем нижнюю-правую грань
			if (i - y) <= 2*h {
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

}

func CheckBorders(mtx [][]rune, x_coord int, y_coord int, cell_width int, cell_height int) (res bool) {
	x_grid_size := len(mtx[0])
	y_grid_size := len(mtx)

	// если фигура выходит за правый край матрицы
	if x_coord+2*cell_height+cell_width > x_grid_size {
		return false
	}

	// если фигура выходит за нижний край матрицы
	if y_coord+2*cell_height > y_grid_size {
		return false
	}

	return true
}

func MakeGrid(mtx [][]rune, cell_width int, cell_height int, num_cells int) {
	curr_cell := 0 //номер текущей добавленной ячейки

	var curr_x int // x-координата фигуры
	var curr_y int // y-координата фигуры
	var curr_row int

	for curr_cell < num_cells {

		// если фигура не выходит за границы сетки
		if CheckBorders(mtx, curr_x, curr_y, cell_width, cell_height) {
			PaintHex(mtx, curr_x, curr_y, cell_width, cell_height)

			curr_x += cell_height + cell_width
			curr_cell++
			if curr_cell%2 == 1 {
				curr_y = (curr_row * 2 * cell_height) + cell_height
			} else {
				curr_y = (curr_row * 2 * cell_height)
			}

		} else {
			curr_x = 0
			curr_row++
		}
	}
}
