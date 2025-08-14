package main

import (
	"bufio"
	"os"
)

func main() {
	// file, err := os.Open("data/9")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer file.Close()

	// reader := bufio.NewReader(strings.NewReader(str))

	// reader := bufio.NewReader(file)

	reader := bufio.NewReader(os.Stdin)

}
