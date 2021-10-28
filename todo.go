package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const DEFAULT_FILE = "/default"

func main() {
	add, rem, file := setFlags()
	file_type, err := getFile(file)
	if err == nil {
		if add != "\n" {
			fmt.Println("add item")
			addItem(add, file_type)
		}
		if rem != 0 {
			removeItem(rem, file_type)
		}
	}
}

/*func fileToMap(file_handle *os.File) (map[int]string, error) {
	var list map[int]string
	fileScanner := bufio.NewScanner(file_handle)
	for fileScanner.Scan() {
		line := fileScanner.Text()
	}
}

func listToFile(bytes []byte, file *os.File) (*os.File, error) {

} */

func getFile(file string) (*os.File, error) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = path + "/.todo/"
	fmt.Println(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if file == DEFAULT_FILE {
		var file_env = os.Getenv("DEFAULT_FILE")
		if file_env == "" {
			os.Setenv("DEFAULT_FILE", DEFAULT_FILE)
			file_env = DEFAULT_FILE
		}
		file = file_env
	}
	return os.OpenFile(path+file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
}

func setFlags() (string, int, string) {
	add := flag.String("a", "\n", "Add item to list")
	rem := flag.Int("r", 0, "Remove item from list")
	file := flag.String("f", DEFAULT_FILE, "File to write to")
	flag.Parse()
	fmt.Println(*add)
	return *add, *rem, *file
}

func getLineNumber(file *os.File) int {
	fileInfo, err := file.Stat()
	line := 1
	if err == nil {
		if fileInfo.Size() == 0 {
			return line
		} else {
			fileScanner := bufio.NewScanner(file)
			for fileScanner.Scan() {
				line += 1
			}
		}
	} else {
		return 0
	}
	return line
}

func addItem(item string, file *os.File) {
	line_number := getLineNumber(file)
	if line_number == 0 {
		return
	}
	item = strconv.Itoa(line_number) + " " + item
	bytes_to_write := []byte(item + "\n")
	bytes, err := file.Write(bytes_to_write)
	if err == nil {
		fmt.Printf("Bytes written: %d\n", bytes)
	}
	defer file.Close()
}

func removeItem(index int, file *os.File) {

}
