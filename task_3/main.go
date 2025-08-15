package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("data/7")
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
	// fmt.Printf("num_of_levels = %d\n", num_of_levels)
	for i := 0; i < num_of_levels; i++ {
		fmt.Fscan(reader, &y_size, &x_size)
		reader.ReadString('\n')
		// fmt.Printf("x_size = %d  y_size = %d\n", x_size, y_size)

		mtx := NewMatrix(x_size, y_size)

		// fmt.Println(str)
		// reader.ReadString('\n')
		for j := 0; j < int(y_size); j++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			// fmt.Printf("%s", line)
			mtx.field[j] = []rune(line)
			// fmt.Println(string(mtx.field[j]))
		}

		mtx.Print()
		mtx.findCellSize()
		fmt.Printf("cell_width = %d  cell_height = %d \n", mtx.cell_width, mtx.cell_height)

	}

}

type Dimension int

type Matrix struct {
	field  [][]rune  // поле
	x_size Dimension //размер поля в символах: ШИРИНА
	y_size Dimension //размер поля в символах: ВЫСОТА

	cell_width  Dimension // размер верхней и нижней граней шестиугольника
	cell_height Dimension // размер каждой из боковых граней шестиугольника

	x_anchor_pos Dimension // x-координата стартовой точки поиска
	y_anchor_pos Dimension // y-координата стартовой точки поиска
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
	// fmt.Println("m.Print() - STARTS")
	// fmt.Printf("x_size = %d\ny_size = %d\n", m.x_size, m.y_size)

	var bldr strings.Builder
	bldr.Grow(int(m.x_size) * int(m.y_size))

	for i := 0; i < int(m.y_size); i++ {
		bldr.WriteString(string(m.field[i]) + "\n")
	}

	fmt.Println(bldr.String())

}

// ищет размерность ячеек
func (m *Matrix) findCellSize() {

	cell_width := 0
	cell_height := 0
	cell_width_calc := false
	cell_height_calc := false

tag:
	for y := 0; y < int(m.y_size); y++ {
		for x := 0; x < int(m.x_size); x++ {
			if m.field[y][x] != ' ' {
				switch m.field[y][x] {
				case '_':
					_x := x
					for !cell_width_calc {

						if (_x >= int(m.x_size)) || (m.field[y][_x] != '_') {

							cell_width_calc = true
							break
						}
						cell_width++
						_x++
					}
				case '/':
					_x := x
					_y := y
					for !cell_height_calc {

						if (_x < 0 || _y >= int(m.y_size)) || (m.field[_y][_x] != '/') {
							cell_height_calc = true
							break
						}
						cell_height++
						_x--
						_y++
					}
				default:
					continue
				}
			}
			if cell_width_calc && cell_height_calc {
				m.cell_width = Dimension(cell_width)
				m.cell_height = Dimension(cell_height)
				break tag
			}
		}
	}
}

// ищет стартовую позицию для дальнейшего обхода по ячекам
func (m *Matrix) findAnchorPos() {

}

// рисует воду там где нужно по условию
func (m *Matrix) PaintWater() {

}
