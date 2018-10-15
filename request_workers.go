package pubnub

import "net/http"

type nonSubMsgType int

const (
	messageTypePublish nonSubMsgType = 1 << iota
	messageTypePAM
)

type JobQResponse struct {
	Resp  *http.Response
	Error error
}

type JobQItem struct {
	Req         *http.Request
	Client      *http.Client
	JobResponse chan *JobQResponse
}

type RequestWorkers struct {
	Workers    chan chan *JobQItem
	MaxWorkers int
	Sem        chan bool
}

type Worker struct {
	Workers    chan chan *JobQItem
	JobChannel chan *JobQItem
	id         int
}

func NewRequestWorkers(workers chan chan *JobQItem, id int) Worker {
	return Worker{
		Workers:    workers,
		JobChannel: make(chan *JobQItem),
		id:         id,
	}
}

func (pw Worker) Process(pubnub *PubNub) {
	go func() {
		for {
			pw.Workers <- pw.JobChannel
			job := <-pw.JobChannel
			res, err := job.Client.Do(job.Req)
			jqr := &JobQResponse{
				Error: err,
				Resp:  res,
			}
			job.JobResponse <- jqr
		}
	}()
}

func (p *RequestWorkers) Start(pubnub *PubNub) {
	pubnub.Config.Log.Println("Start: Running with workers ", p.MaxWorkers)
	for i := 0; i < p.MaxWorkers; i++ {
		pubnub.Config.Log.Println("Start: StartNonSubWorker ", i)
		worker := NewRequestWorkers(p.Workers, i)
		worker.Process(pubnub)
	}
	go p.ReadQueue(pubnub)
}

func (p *RequestWorkers) ReadQueue(pubnub *PubNub) {
	for job := range pubnub.jobQueue {
		pubnub.Config.Log.Println("ReadQueue: Got job for channel ", job.Req)
		go func(job *JobQItem) {
			jobChannel := <-p.Workers
			jobChannel <- job
		}(job)
	}
	pubnub.Config.Log.Println("ReadQueue: Exit")
}

func (p *RequestWorkers) Close() {
	close(p.Workers)
}
