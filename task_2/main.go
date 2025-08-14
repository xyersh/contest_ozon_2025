package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// var str = `30 10 4 2 6`

var str = `3 3 1 1 1`

func main() {
	file, err := os.Open("data/9")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// reader := bufio.NewReader(strings.NewReader(str))

	reader := bufio.NewReader(file)

	// reader := bufio.NewReader(os.Stdin)

	var (
		x_grid_size int
		y_grid_size int
		hex_width   int
		hex_height  int
		num_of_hex  int
	)

	line, _ := reader.ReadString('\n')
	fmt.Sscan(line, &x_grid_size, &y_grid_size, &hex_width, &hex_height, &num_of_hex)

	// fmt.Printf("x_grid_size=%d y_grid_size=%d\nhex_width=%d hex_height=%d\nnum_of_hex=%d\n", x_grid_size, y_grid_size, hex_width, hex_height, num_of_hex)

	matrix := CreateMatrix(x_grid_size, y_grid_size)
	// PaintHex(matrix, 5, 1, 4, 3)
	// PaintHex(matrix, 10, 1, 4, 3)

	MakeGrid(matrix, hex_width, hex_height, num_of_hex)

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

func PaintHex(mtx [][]rune, x_coord, y_coord, w, h int) {

	// mtx - инициализированный двухмерный массив рун. В него рисуем фигуру.
	// x_coord - x-координата фигуры на сетке "mtx"
	// y_coord - y-координата фигуры на сетке "mtx"
	// w - width фигуры
	// h - hgeight фигуры

	// Расчитываем нижнюю-правую границу фигуры
	x_end := x_coord + 2*h + w
	y_end := y_coord + 2*h

	// Проверка на вхождение за границы сетки
	if x_end > len(mtx[0]) {
		x_end = len(mtx[0]) - 1
	}
	if y_end >= len(mtx) {
		y_end = len(mtx) - 1
	}

	// i - ряды
	for i := y_coord; i <= y_end; i++ {

		//j - колонки
		for j := x_coord; j < x_end; j++ {
			i_rel := i - y_coord
			j_rel := j - x_coord

			//рисуем верхнюю и нижнюю грани фигуры
			if (i_rel == 0) || (i_rel == 2*h) {
				if (j_rel < (h + w)) && (j_rel >= h) {
					mtx[i][j] = '_'
					// fmt.Println("test")
				}
			}

			// рисуем верхнюю-левую грань
			if i_rel <= h {
				if i_rel > 0 {
					if i_rel+j_rel == h {
						mtx[i][j] = '/'
					}
				}
			}
			// рисуем верхнюю-правую грань
			if i_rel > 0 {
				if j_rel < h*2+w {
					if j_rel-i_rel == w+h-1 {
						mtx[i][j] = '\\'
					}
				}
			}

			// рисуем нижнюю-левую грань
			if i_rel <= 2*h {
				if i_rel > h {
					if j_rel < h*2+w {
						if i_rel-j_rel == h+1 {
							// fmt.Printf("i=%d j=%d", i, j)
							mtx[i][j] = '\\'
						}
					}
				}
			}

			// рисуем нижнюю-правую грань
			if i_rel <= 2*h {
				if i_rel > 0 {
					if j_rel < h*2+w {
						if i_rel+j_rel == h*3+w {
							mtx[i][j] = '/'
						}
					}
				}
			}
		}
	}

}

func CheckBorders(mtx [][]rune, x_coord int, y_coord int, cell_width int, cell_height int) (res string) {
	// fmt.Printf("CheckBorders: \n	x_coord=%d y_coord=%d\n	cell_width=%d cell_height=%d\n", x_coord, y_coord, cell_width, cell_height)
	x_grid_size := len(mtx[0])
	y_grid_size := len(mtx)

	// если фигура выходит за правый край матрицы
	if x_coord+2*cell_height+cell_width > x_grid_size {
		return "RIGHT_BORDER"
	}

	// если фигура выходит за нижний край матрицы
	if y_coord+2*cell_height >= y_grid_size {
		return "BOTTOM_BORDER"
	}

	return "OK"
}

func MakeGrid(mtx [][]rune, cell_width int, cell_height int, num_cells int) {

	var curr_x int // x-координата фиsгуры
	var curr_y int // y-координата фигуры

	var curr_cell int //номер текущей добавленной ячейки
	var curr_row int
	var cell_in_row int // номер яцейки в ряду

	for curr_cell < num_cells {

		if cell_in_row%2 == 1 {
			// fmt.Printf("odd\n")
			curr_y = (curr_row * 2 * cell_height) + cell_height
		} else {
			// fmt.Printf("even\n")
			curr_y = (curr_row * 2 * cell_height)
		}

		// fmt.Printf("curr_row=%d  curr_x=%d  curr_y=%d cell_in_row=%d\n", curr_row, curr_x, curr_y, cell_in_row)

		//  если фигура не выходит за границы сетки
		switch CheckBorders(mtx, curr_x, curr_y, cell_width, cell_height) {
		case "OK":
			// fmt.Printf("CheckBorders = OK\n")
			PaintHex(mtx, curr_x, curr_y, cell_width, cell_height)
			curr_x += cell_height + cell_width
			curr_cell++

			cell_in_row += 1

		case "BOTTOM_BORDER":
			// fmt.Printf("CheckBorders = BOTTOM_BORDER\n")
			if curr_x == 0 {
				break // если выходим за границы на  первом же элементе в ряду - выходим нафиг
			}
			curr_x += cell_height + cell_width
			cell_in_row++
		case "RIGHT_BORDER":
			// fmt.Printf("CheckBorders = RIGHT_BORDER\n")

			curr_x = 0
			curr_row++
			cell_in_row = 0

		}
	}
}
