package main

import (
	"TaskRunner/taskqueue"
	"flag"
	"github.com/golang/glog"
	"path/filepath"
)

func main() {
	flag.Parse()
	glog.Info("Starting up")
	path, _ := filepath.Abs("./tasks/tasks.yaml")
	var tq = taskqueue.New(path)

	tq.RunTasks()
}
