package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var str = `2
10 12
20 24
`

func main() {
	file, err := os.Open("data/1")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// reader := bufio.NewReader(strings.NewReader(str))

	reader := bufio.NewReader(file)

	// reader := bufio.NewReader(os.Stdin)

	var (
		num_of_levels int
		x_size        Dimension
		y_size        Dimension
	)

	fmt.Fscan(reader, &num_of_levels)
	fmt.Printf("num_of_levels = %d\n", num_of_levels)
	for i := 0; i < num_of_levels; i++ {
		fmt.Fscan(reader, &x_size, &y_size)
		fmt.Printf("x_size = %d  y_size = %d\n", x_size, y_size)

		mtx := NewMatrix(x_size, y_size)
		reader.ReadString('\n')

		for j := 0; j < int(y_size); j++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			// fmt.Printf("%s", line)
			mtx.field[j] = []rune(line)
			fmt.Println(string(mtx.field[j]))
		}
		mtx.Print()

	}

}

type Dimension int

type Matrix struct {
	field  [][]rune  // поле
	x_size Dimension //размер поля в символах: ШИРИНА
	y_size Dimension //размер поля в символах: ВЫСОТА

	cell_width  Dimension // размер верхней и нижней граней шестиугольника
	cell_height Dimension // размер каждой из боковых граней шестиугольника
}

func NewMatrix(x_size, y_size Dimension) *Matrix {
	res := &Matrix{
		x_size: x_size,
		y_size: y_size,
	}

	res.field = make([][]rune, y_size)
	for i := Dimension(0); i < y_size; i++ {
		res.field[i] = make([]rune, x_size)
	}

	return res
}

func (m *Matrix) Print() {
	fmt.Println("m.Print() - STARTS")
	fmt.Printf("x_size = %d\ny_size = %d\n", m.x_size, m.x_size)

	var bldr strings.Builder
	bldr.Grow(int(m.x_size+2) * int(m.y_size+2))
	// mtx_wdth := len(matrix[0])

	// //верхняя граница
	// bldr.WriteRune('+')
	// bldr.WriteString(strings.Repeat("-", int(m.x_size)))
	// bldr.WriteString("+\n")

	for i := 0; i < int(m.y_size); i++ {
		// bldr.WriteString("|" + string(m.field[i]) + "|\n")
		bldr.WriteString(string(m.field[i]) + "\n")
	}

	// // нижняя граница
	// bldr.WriteRune('+')
	// bldr.WriteString(strings.Repeat("-", int(m.x_size)))
	// bldr.WriteString("+\n")

	fmt.Println(bldr.String())

}
