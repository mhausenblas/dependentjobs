package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a Kube Job or CronJob resource.
type Job struct {
	ID                string        `yaml:"id" json:"id"`
	Name              string        `yaml:"name" json:"name"`
	Exectime          time.Duration `yaml:"-" json:"exectime"`
	Status            string        `yaml:"-" json:"status"`
	Dependents        []Job         `yaml:"deps" json:"deps"`
	CompletedUpstream chan bool     `yaml:"-" json:"-"`
}

// New creates a new job.
func newjob(id, name string, numupstream int) Job {
	j := Job{
		ID:                id,
		Name:              name,
		Status:            "scheduled",
		CompletedUpstream: make(chan bool, numupstream),
	}
	return j
}

// AddDep adds one or more dependent jobs.
func (j *Job) adddep(depj ...Job) {
	for _, d := range depj {
		j.Dependents = append(j.Dependents, d)
	}
}

// Launch launches a job, making sure it's only executed
// when all upstream jobs have completed.
func (j Job) launch(wg *sync.WaitGroup) {
	j.wait4upstream()
	fmt.Printf(j.render("Launched"))
	j.execute(wg)
	for _, d := range j.Dependents {
		go d.launch(wg)
	}
}

func (j Job) execute(wg *sync.WaitGroup) {
	defer wg.Done()
	et := time.Duration(500 + 1000000*rand.Intn(2000))
	j.Exectime = et
	time.Sleep(et)
	j.Status = "completed"
	fmt.Printf(j.render("Executed"))
	for _, d := range j.Dependents {
		d.CompletedUpstream <- true
	}
}

func (j Job) wait4upstream() {
	upstreamcount := cap(j.CompletedUpstream)
	if upstreamcount == 0 {
		return
	}
	i := 0
	for {
		select {
		case <-j.CompletedUpstream:
			i++
		}
		if upstreamcount == i {
			break
		}
	}
}

func (j Job) render(msg string) string {
	now := time.Now().Unix()
	return fmt.Sprintf("%v| %s: <%v,%v,%v,%v>\n", now, msg, j.ID, j.Name, j.Exectime, j.Status)
}
