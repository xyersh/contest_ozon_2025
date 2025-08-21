package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {

		// n - кол-во строк поля
		// m - кол-во столбцов поля
		var n, m int

		var (
			// координаты начального пункта
			startR int // Y
			startC int // X

			// координаты конечного пункта
			endR int // Y
			endC int // X
		)

		fmt.Fscan(reader, &n, &m)

		// сетка с гексагонами высотой n строк
		grid := make([][]rune, n)
		for j := 0; j < n; j++ {

			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			grid[j] = []rune(line)
		}

		fmt.Fscan(reader, &startR, &startC) //
		fmt.Fscan(reader, &endR, &endC)     //

		// Adjust to 0-indexed
		startR--
		startC--
		endR--
		endC--

		solver := NewSolver(n, m, grid)

		// Определяем HexWidth и HexHeight один раз по стартовой точке
		solver.HexWidth, solver.HexHeight = solver.findHexagonDimensions(startR, startC)

		result := solver.solve(startR, startC, endR, endC)
		fmt.Fprintln(writer, result)
	}
}

type Point struct {
	R, C int
}

type HexagonID struct {
	R, C int // Координаты верхнего левого символа '/'
}

type Solver struct {
	N, M       int
	Grid       [][]rune
	HexWidth   int // W
	HexHeight  int // H
	VisitedHex map[HexagonID]bool
	HexMap     map[Point]HexagonID // Кэш: Point -> HexagonID
	HexType    map[HexagonID]bool  // true = land, false = sea
}

// NewSolver создает новый экземпляр Solver
func NewSolver(n, m int, grid [][]rune) *Solver {
	return &Solver{
		N:          n,
		M:          m,
		Grid:       grid,
		VisitedHex: make(map[HexagonID]bool),
		HexMap:     make(map[Point]HexagonID),
		HexType:    make(map[HexagonID]bool),
	}
}

// isValid проверяет, находятся ли координаты в пределах поля
func (s *Solver) isValid(r, c int) bool {
	return r >= 0 && r < s.N && c >= 0 && c < s.M
}

// findHexagonDimensions определяет ширину (W) и высоту (H) шестиугольника
// по заданной стартовой точке (r, c), которая гарантированно является пробелом внутри шестиугольника.
func (s *Solver) findHexagonDimensions(startR, startC int) (width, height int) {
	// Используем BFS для поиска всех пробелов, принадлежащих одному шестиугольнику
	q := []Point{{startR, startC}}
	visited := make(map[Point]bool)
	visited[Point{startR, startC}] = true

	minR, maxR := startR, startR
	// minC, maxC := startC, startC // Эти переменные не используются для определения W и H

	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.R < minR {
			minR = curr.R
		}
		if curr.R > maxR {
			maxR = curr.R
		}
		// if curr.C < minC { minC = curr.C } // Не нужны для определения W и H
		// if curr.C > maxC { maxC = curr.C } // Не нужны для определения W и H

		for i := 0; i < 4; i++ {
			nR, nC := curr.R+dr[i], curr.C+dc[i]
			if s.isValid(nR, nC) && s.Grid[nR][nC] == ' ' && !visited[Point{nR, nC}] {
				visited[Point{nR, nC}] = true
				q = append(q, Point{nR, nC})
			}
		}
	}

	// Теперь minR, maxR охватывают все пробелы шестиугольника по вертикали.
	// Верхняя граница '_' будет на строке minR - 1.
	// Нижняя граница '_' будет на строке maxR + 1.

	// Высота шестиугольника в строках: (maxR + 1) - (minR - 1) + 1 = maxR - minR + 3
	// H = (total_rows - 1) / 2
	height = (maxR - minR + 3 - 1) / 2 // (maxR - minR + 2) / 2

	// Ширина шестиугольника (количество '_')
	// На строке minR - 1, ищем количество '_'
	// Найдем первую и последнюю колонку с '_' на строке minR - 1
	firstUnderscoreC := -1
	lastUnderscoreC := -1
	for c := 0; c < s.M; c++ {
		if s.isValid(minR-1, c) && s.Grid[minR-1][c] == '_' {
			if firstUnderscoreC == -1 {
				firstUnderscoreC = c
			}
			lastUnderscoreC = c
		}
	}
	if firstUnderscoreC != -1 {
		width = lastUnderscoreC - firstUnderscoreC + 1
	} else {
		width = 0 // Если '_' не найден, что не должно произойти для суши
	}

	return width, height
}

// getHexagonID определяет уникальный идентификатор шестиугольника (HexagonID)
// для заданной точки (r, c).
// HexagonID - это координаты верхнего левого символа '/' шестиугольника.
func (s *Solver) getHexagonID(r, c int) HexagonID {
	p := Point{r, c}
	if hexID, ok := s.HexMap[p]; ok {
		return hexID
	}

	// Используем BFS для поиска всех пробелов, принадлежащих одному шестиугольнику
	q := []Point{p}
	visited := make(map[Point]bool)
	visited[p] = true

	hexPoints := []Point{}
	hexPoints = append(hexPoints, p)

	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		for i := 0; i < 4; i++ {
			nR, nC := curr.R+dr[i], curr.C+dc[i]
			if s.isValid(nR, nC) && s.Grid[nR][nC] == ' ' && !visited[Point{nR, nC}] {
				visited[Point{nR, nC}] = true
				q = append(q, Point{nR, nC})
				hexPoints = append(hexPoints, Point{nR, nC})
			}
		}
	}

	// Теперь, когда у нас есть все пробелы шестиугольника,
	// найдем его верхний левый символ '/'.
	// Для этого, найдем самую верхнюю строку с пробелами (minR_space)
	minR_space := r // Инициализируем minR_space и minC_space
	minC_space := c

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.R < minR_space {
			minR_space = curr.R
		}
		if curr.C < minC_space {
			minC_space = curr.C
		} // Обновляем minC_space

		for i := 0; i < 4; i++ {
			nR, nC := curr.R+dr[i], curr.C+dc[i]
			if s.isValid(nR, nC) && s.Grid[nR][nC] == ' ' && !visited[Point{nR, nC}] {
				visited[Point{nR, nC}] = true
				q = append(q, Point{nR, nC})
				hexPoints = append(hexPoints, Point{nR, nC})
			}
		}
	}

	// Строка с верхней границей '_' будет minR_space - 1.
	rTopHex := minR_space - 1

	// Найдем c_left_top_underscore на строке rTopHex.
	// Ищем '_' влево от minC_space, затем вправо.
	cLeftTopUnderscoreHex := minC_space
	for cLeftTopUnderscoreHex >= 0 && s.isValid(rTopHex, cLeftTopUnderscoreHex) && s.Grid[rTopHex][cLeftTopUnderscoreHex] == '_' {
		cLeftTopUnderscoreHex--
	}
	cLeftTopUnderscoreHex++ // Возвращаемся к первому '_'

	hexID := HexagonID{R: rTopHex + s.HexHeight, C: cLeftTopUnderscoreHex - s.HexHeight}

	// Кэшируем HexagonID для всех точек, принадлежащих этому шестиугольнику
	for _, pt := range hexPoints {
		s.HexMap[pt] = hexID
	}

	return hexID
}

// isLand определяет, является ли шестиугольник сушей.
// HexagonID - это координаты верхнего левого символа '/' шестиугольника.
func (s *Solver) isLand(hexID HexagonID) bool {
	if isLand, ok := s.HexType[hexID]; ok {
		return isLand
	}

	// Проверяем 6 сторон шестиугольника
	// HexagonID.R, HexagonID.C - это координаты верхнего левого символа '/'
	// W = s.HexWidth, H = s.HexHeight

	// 1. Верхняя сторона '_'
	// Находится на строке hexID.R - H
	// Начинается в hexID.C + H
	// Длина W
	for i := 0; i < s.HexWidth; i++ {
		if !s.isValid(hexID.R-s.HexHeight, hexID.C+s.HexHeight+i) || s.Grid[hexID.R-s.HexHeight][hexID.C+s.HexHeight+i] != '_' {
			s.HexType[hexID] = false
			return false
		}
	}

	// 2. Левая верхняя сторона '/'
	// Начинается в hexID.R, hexID.C
	// Длина H
	for i := 0; i < s.HexHeight; i++ {
		if !s.isValid(hexID.R+i, hexID.C+i) || s.Grid[hexID.R+i][hexID.C+i] != '/' {
			s.HexType[hexID] = false
			return false
		}
	}

	// 3. Правая верхняя сторона '\'
	// Начинается в hexID.R, hexID.C + W + H + 1
	// Длина H
	for i := 0; i < s.HexHeight; i++ {
		if !s.isValid(hexID.R+i, hexID.C+s.HexWidth+s.HexHeight+1-i) || s.Grid[hexID.R+i][hexID.C+s.HexWidth+s.HexHeight+1-i] != '\\' {
			s.HexType[hexID] = false
			return false
		}
	}

	// 4. Нижняя сторона '_'
	// Находится на строке hexID.R + H + H
	// Начинается в hexID.C + H
	// Длина W
	for i := 0; i < s.HexWidth; i++ {
		if !s.isValid(hexID.R+s.HexHeight, hexID.C+s.HexHeight+i) || s.Grid[hexID.R+s.HexHeight][hexID.C+s.HexHeight+i] != '_' {
			s.HexType[hexID] = false
			return false
		}
	}

	// 5. Левая нижняя сторона '\'
	// Начинается в hexID.R + H, hexID.C + H - 1
	// Длина H
	for i := 0; i < s.HexHeight; i++ {
		if !s.isValid(hexID.R+s.HexHeight+i, hexID.C+s.HexHeight-1-i) || s.Grid[hexID.R+s.HexHeight+i][hexID.C+s.HexHeight-1-i] != '\\' {
			s.HexType[hexID] = false
			return false
		}
	}

	// 6. Правая нижняя сторона '/'
	// Начинается в hexID.R + H, hexID.C + W + H + 1
	// Длина H
	for i := 0; i < s.HexHeight; i++ {
		if !s.isValid(hexID.R+s.HexHeight+i, hexID.C+s.HexWidth+s.HexHeight+1+i) || s.Grid[hexID.R+s.HexHeight+i][hexID.C+s.HexWidth+s.HexHeight+1+i] != '/' {
			s.HexType[hexID] = false
			return false
		}
	}

	s.HexType[hexID] = true
	return true
}

// getNeighbors возвращает список соседних шестиугольников.
// Используем 6 "проверочных" точек, которые находятся в соседних шестиугольниках.
func (s *Solver) getNeighbors(hexID HexagonID) []HexagonID {
	neighbors := []HexagonID{}

	// Смещения для 6 "проверочных" точек относительно HexagonID (верхний левый '/')
	// Эти точки должны быть пробелами, принадлежащими соседним шестиугольникам.
	// (dr, dc)
	testPoints := []Point{
		// Сосед справа: точка внутри него, справа от правой границы текущего
		{hexID.R + s.HexHeight, hexID.C + s.HexWidth + 2*s.HexHeight + 1},
		// Сосед слева: точка внутри него, слева от левой границы текущего
		{hexID.R + s.HexHeight, hexID.C - 1},
		// Сосед сверху-справа: точка внутри него, над верхней правой границей текущего
		{hexID.R - 1, hexID.C + s.HexWidth + s.HexHeight},
		// Сосед сверху-слева: точка внутри него, над верхней левой границей текущего
		{hexID.R - 1, hexID.C + s.HexHeight - 1},
		// Сосед снизу-справа: точка внутри него, под нижней правой границей текущего
		{hexID.R + 2*s.HexHeight + 1, hexID.C + s.HexWidth + s.HexHeight},
		// Сосед снизу-слева: точка внутри него, под нижней левой границей текущего
		{hexID.R + 2*s.HexHeight + 1, hexID.C + s.HexHeight - 1},
	}

	for _, p := range testPoints {
		if s.isValid(p.R, p.C) && s.Grid[p.R][p.C] == ' ' {
			neighborHexID := s.getHexagonID(p.R, p.C)
			// Проверяем, что это не текущий шестиугольник
			if neighborHexID != hexID {
				neighbors = append(neighbors, neighborHexID)
			}
		}
	}

	return neighbors
}

// solve выполняет BFS для поиска пути между двумя шестиугольниками суши.
func (s *Solver) solve(startR, startC, endR, endC int) string {
	startHex := s.getHexagonID(startR, startC)
	endHex := s.getHexagonID(endR, endC)

	if !s.isLand(startHex) || !s.isLand(endHex) {
		return "NO"
	}

	q := []HexagonID{startHex}
	s.VisitedHex[startHex] = true

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr == endHex {
			return "YES"
		}

		for _, neighbor := range s.getNeighbors(curr) {
			if !s.VisitedHex[neighbor] && s.isLand(neighbor) {
				s.VisitedHex[neighbor] = true
				q = append(q, neighbor)
			}
		}
	}

	return "NO"
}
