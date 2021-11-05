package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const DEFAULT_FILE = "default.json"
const DEFAULT_FOLDER = "/.todo/"

func getHomeDir() string {
	path, err := os.UserHomeDir()
	if err == nil {
		return path
	}
	return path
}

type Task struct {
	ID       int
	Priority int
	Item     string
	Error    string
}

type Tasks []Task

func getTask() Task {
	add := flag.String("a", " ", "Add task to list")
	remove := flag.Int("r", 0, "Remove task from list")
	priority := flag.Int("p", 0, "Specify priority for task")
	flag.Parse()
	if *add == " " && *remove == 0 && *priority == 0 {
		return Task{Error: "No args"}
	}
	return Task{ID: *remove, Priority: *priority, Item: *add, Error: ""}
}

func getJsonFile() (*os.File, error) {
	path := getHomeDir() + DEFAULT_FOLDER
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	return os.OpenFile(path+DEFAULT_FILE, os.O_RDWR|os.O_CREATE, 0755)

}

func writeJson(json []byte, filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0755)
	bytes, err := file.Write(json)
	if err != nil {
		fmt.Println(err)
	}
	if bytes > 0 {
		return
	}
}

func unmarshalJsonfile(jsonFile *os.File) Tasks {
	var tasks []Task
	bytes, _ := ioutil.ReadAll(jsonFile)
	if len(bytes) == 0 {
		return tasks
	}
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		panic(err)
	}
	jsonFile.Close()
	return tasks
}

func marshalJson(ts Tasks) []byte {
	bytes, _ := json.Marshal(ts)
	return bytes
}

func displayTodo(tasks []Task) {
	toWrite := ""
	for _, task := range tasks {
		toWrite = toWrite + strconv.Itoa(task.ID) + ". " + task.Item + "\n"
	}
	os.Stdout.Write([]byte(toWrite))
}

type Actions interface {
	remove(task Task)
	add(task Task)
}

func (ts *Tasks) add(t Task) {
	size := len(*ts)
	t.ID = size + 1
	*ts = append(*ts, t)
}
func (ts *Tasks) remove(id int) {
	id = id - 1
	var tsNew Tasks
	for index, task := range *ts {
		if index < id {
			tsNew = append(tsNew, task)
		} else if index == id {
			continue
		} else {
			tsNew = append(tsNew, Task{ID: (task.ID - 1), Priority: task.Priority, Item: task.Item, Error: task.Error})
		}
	}
	*ts = tsNew
}

func main() {
	var task = getTask()
	file, err := getJsonFile()
	if err != nil {
		return
	}
	tasks := unmarshalJsonfile(file)
	if task.Error == "No args" {
		displayTodo(tasks)
		return
	}
	if task.ID == 0 {
		tasks.add(task)
	} else {
		tasks.remove(task.ID)
	}
	displayTodo(tasks)
	json_format := marshalJson(tasks)
	writeJson(json_format, getHomeDir()+DEFAULT_FOLDER+DEFAULT_FILE)

}
