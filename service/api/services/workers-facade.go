package services

import "sync"

type Work struct {
	idx int
	res interface{}
	err error
}

func NewWork(id int, res interface{}, err error) *Work {
	return &Work{
		idx: id,
		res: res,
		err: err,
	}
}

type SendFunc func(res interface{}, err error)
type Worker func(newWorkRes SendFunc)

type Job struct {
	worker Worker
}

func NewJob(worker Worker) *Job {
	return &Job{
		worker: worker,
	}
}

func NewSendWorkRes(id int, out chan<- *Work) SendFunc {
	return func(res interface{}, err error) {
		work := NewWork(id, res, err)
		out <- work
	}
}

func NewWorkersFacade(jobs ...*Job) <-chan *Work {
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	out := make(chan *Work)

	for i, job := range jobs {
		sendWorkRes := NewSendWorkRes(i, out)

		go func(i int, job *Job) {
			defer wg.Done()
			job.worker(sendWorkRes)
		}(i, job)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

	// result := make([]interface{}, len(jobs))
	// for work := range out {
	// 	if work.err != nil {
	// 		return nil, work.err
	// 	}
	// 	result[work.idx] = work.res
	// }
	// return result
}
