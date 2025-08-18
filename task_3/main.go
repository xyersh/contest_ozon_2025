package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var str = `1
10 11
       _   
      / \  
     /   \ 
   _/     \
  / \     /
 /   \   / 
/     \_/  
\     /    
 \   /     
  \_/      
`

func main() {
	file, err := os.Open("data/1")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// reader := bufio.NewReader(strings.NewReader(str))

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

		// mtx.Print()

		mtx.findCellSize()
		// fmt.Printf("cell_width = %d   cell_height=%d\n", mtx.cell_width, mtx.cell_height)
		mtx.PaintWater()
		// mtx._printCopy()
		mtx.Print()

		// fmt.Printf("cell_width = %d  cell_height = %d \n", mtx.cell_width, mtx.cell_height)

	}

}

type Dimension int

type Matrix struct {
	field      [][]rune // поле
	field_copy [][]rune

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

	res.field_copy = make([][]rune, y_size)
	for i := Dimension(0); i < y_size; i++ {
		res.field_copy[i] = make([]rune, x_size)
		for j := 0; j < int(res.x_size); j++ {
			res.field_copy[i][j] = ' '
		}
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

func (m *Matrix) _printCopy() {
	// fmt.Println("m.Print() - STARTS")
	// fmt.Printf("x_size = %d\ny_size = %d\n", m.x_size, m.y_size)

	var bldr strings.Builder
	bldr.Grow(int(m.x_size) * int(m.y_size))

	for i := 0; i < int(m.y_size); i++ {
		bldr.WriteString(string(m.field_copy[i]) + "\n")
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

// рисует воду там где нужно по условию
func (m *Matrix) PaintWater() {
	// fmt.Printf("PaintWater() - START\n")
	for y := 0; y < (int(m.y_size) - 2*int(m.cell_height)); y += int(m.cell_height) {
		// fmt.Printf("y=%d\n", y)
		cell_in_row := 0
		var _y = y

		start_x_idx := strings.IndexRune(string(m.field[y]), '_') - int(m.cell_height)

		for x := start_x_idx; x <= int(m.x_size)-(2*int(m.cell_height)+int(m.cell_width))+1; x += int(m.cell_height) + int(m.cell_width) {
			// fmt.Printf("x=%d\n", x)
			// if cell_in_row%2 == 1 {
			// 	_y = y + int(m.cell_height)
			// } else {
			// 	_y = y
			// }
			// fmt.Printf("_y=%d\n", _y)
			if m._isCell(x, _y) {
				m._markGroundCell(x, _y)
			}

			// fmt.Printf("y = %d   x = %d\n", _y, x)
			cell_in_row++
		}
	}
	m._markWater()
}

// проверяет, является ли содержиимое блока ячейкой
func (m *Matrix) _isCell(x, y int) (res bool) {
	// fmt.Printf("_isCell(x=%d, y=%d)\n", x, y)
	// defer func() {
	// 	fmt.Printf("_isCell(x=%d, y=%d) = %t\n", x, y, res)
	// }()

	res = true
	if x < 0 || y < 0 {
		res = false
		return
	}

	// проверка верхней стороны
	if m.field[y][x+int(m.cell_height)] != '_' {
		// fmt.Printf("		test 1 sym=%c\n", m.field[y][x+int(m.cell_height)])
		res = false
		return
	}
	// проверка нижней стороны
	if m.field[y+2*int(m.cell_height)][x+int(m.cell_height)] != '_' {
		// fmt.Printf("		test 2 sym= %c\n", m.field[y][x+int(m.cell_height)])
		res = false
		return
	}

	// проверка  левой верхней стороны
	if m.field[y+1][x+int(m.cell_height)-1] != '/' {
		// fmt.Println("		test 3")
		res = false
		return
	}

	// проверка левой нижней  стороны
	if m.field[y+2*int(m.cell_height)][x+int(m.cell_height)-1] != '\\' {
		// fmt.Println("		test 4")
		res = false
		return
	}

	// проверка правой верхней стороны
	if m.field[y+1][x+int(m.cell_height+m.cell_width)] != '\\' {
		// fmt.Println("		test 5")
		res = false
		return
	}

	// проверка  стороны
	if m.field[y+2*int(m.cell_height)][x+int(m.cell_height+m.cell_width)] != '/' {
		// fmt.Println("		test 6")
		res = false
		return
	}

	return
}

func (m *Matrix) _getBorder(x_pos, y_pos int) [][]rune {
	res := make([][]rune, 2*m.cell_height+1)
	for i := 0; i < len(res); i++ {
		res[i] = m.field[i+y_pos][x_pos : x_pos+2*int(m.cell_height)+int(m.cell_width)]
	}

	return res
}

func (m *Matrix) _markGroundCell(x_pos, y_pos int) {
	// fmt.Printf("_markGroundCell(x_pos=%d, y_pos=%d)\n", x_pos, y_pos)
	//маркируем ячейки - как землю ('*')

	border := m._getBorder(x_pos, y_pos)

	// _printBorder(border)
	// fmt.Printf(" border len = %d\n", len(border))
	for y := 0; y < len(border); y++ {
		// fmt.Printf("y=%d test %s\n", y, string(border[y]))

		// проверка верхней и нижней сторон ячейки
		if y == 0 || y == len(border)-1 {

			for x := 0; x < len(border[0]); x++ {

				if border[y][x] == '_' {
					m.field_copy[y_pos+y][x_pos+x] = '*'
				}
			}
		}
		// fmt.Println("weqwqwee")
		// m._printCopy()

		if y > 0 && y < len(border) {
			// fmt.Printf(" test2 %s\n", string(border[y]))
			//верхние боковые стороны (проверка символов '/'  и '\')
			if y <= int(m.cell_height) {
				first_idx := strings.IndexRune(string(border[y]), '/')
				last_idx := strings.IndexRune(string(border[y]), '\\')

				// fmt.Printf("upper side  %d'nth row first_idx=%d  last_idx=%d \n ", y, first_idx, last_idx)

				for x := first_idx; x <= last_idx; x++ {
					m.field_copy[y_pos+y][x_pos+x] = '*'
				}
			}

			//нижние боковые стороны (проверка символов '/'  и '\')
			if y > int(m.cell_height) {

				first_idx := strings.IndexRune(string(border[y]), '\\')
				last_idx := strings.IndexRune(string(border[y]), '/')
				// fmt.Printf("lower side  %d'nth row first_idx=%d last_idx=%d\n", y, first_idx, last_idx)
				for x := first_idx; x <= last_idx; x++ {

					m.field_copy[y_pos+y][x_pos+x] = '*'
				}
			}

		}
		// fmt.Println("!!!!!")
		// m._printCopy()
	}

}

func (m *Matrix) _markWater() {
	// маркируем все остальное как воду ('~')
	for y := 0; y < int(m.y_size); y++ {
		for x := 0; x < int(m.x_size); x++ {
			if m.field_copy[y][x] != '*' {
				m.field[y][x] = '~'
			}
		}
	}
}

func _printBorder(b [][]rune) {
	fmt.Println("_printBorder() - START")
	for i := 0; i < len(b); i++ {
		fmt.Println(string(b[i]))
	}
}
