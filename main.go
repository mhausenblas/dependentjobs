package main

func main() {
	jobs := []Job{
		New("job 1", 0),
		New("job 2", 1),
		New("job 3", 1),
		New("job 4", 2),
	}
	jobs[1].AddDep(jobs[3])
	jobs[2].AddDep(jobs[3])
	jobs[0].AddDep(jobs[1], jobs[2])
	Run(jobs)
}
