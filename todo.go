package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const DEFAULT_FILE = "/default"
const DEFAULT_FOLDER = "/.todo/"

func getHomeDir() string {
	file, err := os.Getwd()
	if err == nil {
		return file
	}
	return file
}

func main() {
	add, rem, file := setFlags()
	file_type, err := getFile(file)
	if err == nil {
		if add != "\n" {
			fmt.Println("add item")
			addItem(add, file_type)
		}
		if rem != 0 {
			fmt.Println("flag")
			map_of_items := fileToMap(file_type)
			file_type.Close()
			file_type, err := os.OpenFile(getHomeDir()+DEFAULT_FOLDER+file, os.O_RDWR|os.O_TRUNC, 0755)
			if err == nil {
				fmt.Println(map_of_items)
				delete(map_of_items, rem)
				fmt.Println(map_of_items)
				mapToFile(file_type, map_of_items)
			}
			file_type.Close()
		}
	}
}

func fileToMap(file_handle *os.File) map[int]string {
	list := make(map[int]string)
	fileScanner := bufio.NewScanner(file_handle)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		tokens := strings.Fields(line)
		index, err := strconv.Atoi(tokens[0])
		if err == nil {
			size := len(tokens) - 1
			item := tokens[1 : size+1]
			list[index] = strings.Join(item, " ")
		}
	}
	return list
}

func mapToFile(file_handle *os.File, list map[int]string) {
	keys := make([]int, len(list))
	i := 0
	for key := range list {
		keys[i] = key
		i++
	}
	new_keys := sort.IntSlice(keys)
	fmt.Println(keys)
	for _, key := range new_keys {
		fmt.Println(key)
		item_to_write := strconv.Itoa(key) + " " + list[key] + "\n"
		fmt.Println(item_to_write)
		bytes, err := file_handle.Write([]byte(item_to_write))
		if err == nil && bytes > 0 {
			continue
		}
	}
}

func getFile(file string) (*os.File, error) {
	path := getHomeDir()
	path = path + DEFAULT_FOLDER
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
	default_file := flag.String("d", "no default", "Make default file")
	flag.Parse()
	fmt.Println(*add)
	return *add, *rem, *file, *default_file
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
