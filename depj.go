package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a Kube Job or CronJob resource.
type Job struct {
	ID                string
	Name              string
	Exectime          time.Duration
	Status            string
	Dependents        []Job
	CompletedUpstream chan bool
}

// Run takes a call graph of jobs
// and runs it in order of its dependencies.
func Run(jobs []Job) {
	var wg sync.WaitGroup
	go jobs[0].Launch()
	wg.Wait()
}

// New creates a new job.
func New(name string, numupstream int) Job {
	j := Job{
		ID:                fmt.Sprintf("j%v", time.Now().UnixNano()),
		Name:              name,
		Status:            "scheduled",
		CompletedUpstream: make(chan bool, numupstream),
	}
	return j
}

// AddDep adds one or more dependent jobs.
func (j *Job) AddDep(depj ...Job) {
	for _, d := range depj {
		j.Dependents = append(j.Dependents, d)
	}
}

// Launch launches a job, making sure it's only executed
// when all upstream jobs have completed.
func (j Job) Launch() {
	j.wait4upstream()
	fmt.Printf(j.render("Launched"))
	j.execute()
	for _, d := range j.Dependents {
		go d.Launch()
	}
}

func (j Job) execute() {
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
