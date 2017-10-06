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

// DependentJobs represents a call graph, that is,
// a collection of Kube Jobs or CronJobs with dependencies.
type DependentJobs struct {
	wg   *sync.WaitGroup
	jobs []Job
}

// New creates a new call graph.
func New() DependentJobs {
	dj := DependentJobs{}
	dj.wg = &sync.WaitGroup{}
	return dj
}

// Run takes a call graph of jobs
// and runs it in order of its dependencies.
func (dj DependentJobs) Run() {
	dj.wg.Add(len(dj.jobs))
	go dj.jobs[0].launch(dj.wg)
	dj.wg.Wait()
}

// Add adds a job to the call graph.
func (dj *DependentJobs) Add(name string, numupstream int) {
	j := newjob(name, numupstream)
	dj.jobs = append(dj.jobs, j)
}

// AddDependents adds one or more dependent jobs to a job.
func (dj DependentJobs) AddDependents(job int, depjobs ...int) {
	j := &(dj.jobs[job])
	depj := []Job{}
	for _, r := range depjobs {
		depj = append(depj, dj.jobs[r])
	}
	j.adddep(depj...)
}

// New creates a new job.
func newjob(name string, numupstream int) Job {
	j := Job{
		ID:                fmt.Sprintf("j%v", time.Now().UnixNano()),
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
