package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// DependentJobs represents a call graph, that is,
// a collection of Kube Jobs or CronJobs with dependencies.
type DependentJobs struct {
	wg      *sync.WaitGroup `yaml:"-"`
	jobs    map[string]Job  `yaml:"jobs"`
	callseq chan Job        `yaml:"-"`
	result  []string        `yaml:"-"`
}

var jticks map[string]int

// New creates a new call graph.
func New() DependentJobs {
	dj := DependentJobs{}
	dj.jobs = make(map[string]Job)
	dj.wg = &sync.WaitGroup{}
	rand.Seed(time.Now().UTC().UnixNano())
	return dj
}

// FromFile reads a call graph from a YAML file.
func (dj *DependentJobs) FromFile(cgfile string) error {
	yamlFile, err := ioutil.ReadFile(cgfile)
	if err != nil {
		return err
	}
	spec := make(map[string]Job)
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		return err
	}
	// generate DJ out of spec
	for id, job := range spec {
		dj.Add(id, job.Name, countupstream(spec, id))
		dj.AddDependents(id, spec[id].Dependents...)
		if spec[id].Every > 0 {
			dj.AddPeriodic(id, spec[id].Every)
		}
	}
	return nil
}

// Store stores a call graph into a YAML file.
func (dj DependentJobs) Store(cgfile string) error {
	bytes, err := yaml.Marshal(dj.jobs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cgfile, bytes, 0644)
}

// Run takes a call graph of jobs and runs it in order of its dependencies.
func (dj *DependentJobs) Run() {
	dj.callseq = make(chan Job, len(dj.jobs))
	dj.wg.Add(len(dj.jobs))
	r := dj.Lookup("root")
	go r.launch(*dj, dj.wg)
	dj.wg.Wait()
	// need to close the call sequence channel
	// in order to be able to drain it in CallSeq():
	close(dj.callseq)
}

// Add adds a job to the call graph.
func (dj *DependentJobs) Add(id, name string, numupstream int) {
	j := newjob(id, name, numupstream)
	dj.jobs[id] = j
}

// AddDependents adds one or more dependent jobs to a job.
func (dj DependentJobs) AddDependents(id string, depjobs ...string) {
	j := dj.jobs[id]
	j.adddep(depjobs...)
	dj.jobs[id] = j
}

// AddPeriodic adds a periodic schedule to a job.
func (dj DependentJobs) AddPeriodic(id string, every int) {
	j := dj.jobs[id]
	j.periodic(every - 1)
	dj.jobs[id] = j
}

// Lookup retrieves a job by ID.
func (dj DependentJobs) Lookup(id string) Job {
	return dj.jobs[id]
}

// CallSeq returns the sequence in which the jobs have been called.
func (dj DependentJobs) CallSeq() []string {
	return dj.result
}

// Complete waits until the cycle is complete using the
// call sequence for synchronization.
func (dj *DependentJobs) Complete() {
	for j := range dj.callseq {
		p := fmt.Sprintf("%s %v %v", j.ID, j.Starttime, j.Endtime)
		dj.result = append(dj.result, p)
	}
}

// TimeToRun checks if a periodic job can run.
func (dj DependentJobs) TimeToRun(id string) bool {
	d := dj.Lookup(id)
	if d.Every > 0 {
		if jticks[id] <= d.Every-1 {
			jticks[id]++
			return false
		}
		jticks[id] = 0
	}
	return true
}

// GoString return a canonical string represenation of a dependent job
func (dj DependentJobs) GoString() string {
	res := ""
	for _, j := range dj.jobs {
		res = fmt.Sprintf("%s%#v\n", res, j)
	}
	return res
}

func countupstream(jobs map[string]Job, jobid string) int {
	numupstream := 0
	switch jobid {
	case "root":
		numupstream = 0
	default:
		for _, j := range jobs {
			if contains(j.Dependents, jobid) {
				numupstream++
			}
		}
	}
	return numupstream
}

func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}
