package taskqueue

import (
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type TaskQueue struct {
	Tasks []Task `yaml:"tasks"`
}

type TaskResult struct {
	task   Task
	result bool
	err    error
}

type Job struct {
	task  Task
	index int
}

/*
	Generate a new task queue using tasks defined in `filename`
*/
func New(filename string) *TaskQueue {
	glog.Info("Reading file", filename)
	yamlFile, err := ioutil.ReadFile(filename)

	tq := new(TaskQueue)
	err = yaml.Unmarshal(yamlFile, &tq)

	if err != nil {
		fmt.Println(err)
	}

	glog.Infof("Loaded %d tasks", len(tq.Tasks))

	/* Post-process yaml tasks */
	for idx := range tq.Tasks {
		t := &tq.Tasks[idx]
		if t.From == "now" {
			t.StartFrom = time.Now()
		}
		t.NextRunTime = t.StartFrom
	}

	return tq
}

func (tq *TaskQueue) AddTask(task Task) {
	tq.Tasks = append(tq.Tasks, task)
}

func taskWorker(job Job, jobs <-chan Job, results chan<- Task) {
	// glog.Infof("Task worker for %s ", job.task.Name)
	go job.task.Run(results)
}

/* Main worker function */
func (tq TaskQueue) RunTasks() {
	tasks := make(chan Job, 100)
	results := make(chan Task, 100)
	glog.Info("Running tasks ... forever")

	go func() {
		for job := range tasks {
			// glog.Infof("Incoming task %s", job.task.Name)
			go taskWorker(job, tasks, results)
		}
	}()

	go func() {
		for result := range results {
			glog.Infof("Task %s completed", result.Name)

			if result.Requeue {
				glog.Info(result.NextRunTime)
				tq.AddTask(result)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			for idx := len(tq.Tasks) - 1; idx >= 0; idx-- {
				t := tq.Tasks[idx]
				var currentTime = time.Now()

				if t.NextRunTime.Before(currentTime) && t.StartFrom.Before(currentTime) {
					tq.Tasks = append(tq.Tasks[:idx], tq.Tasks[idx+1:]...)
					tasks <- Job{task: t, index: idx}
				}
			}
		}
	}()

	for {
	}
}
