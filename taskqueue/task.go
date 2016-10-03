package taskqueue

import "time"
import "github.com/golang/glog"

type Task struct {
	Name        string `yaml:"name"`
	Interval    string `yaml:"interval"`
	Requeue     bool   `yaml:"requeue"`
	From        string
	StartFrom   time.Time
	NextRunTime time.Time
}

func (this *Task) Run(results chan<- Task) {
	glog.Info("In task coroutine")
	go func() {
		nextDuration, _ := time.ParseDuration(this.Interval)
		glog.Infof("Task %s running again in %s", this.Name, this.Interval)
		this.NextRunTime = this.NextRunTime.Add(nextDuration)
		results <- *this
	}()
}
