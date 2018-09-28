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
			pubnub.Config.Log.Printf("Process: Worker started: %d", pw.id)
			job := <-pw.JobChannel
			pubnub.Config.Log.Printf("Process: Worker %d processing job %s", pw.id, job.Req)
			res, err := job.Client.Do(job.Req)
			pubnub.Config.Log.Println(res, err)
			jqr := &JobQResponse{
				Error: err,
				Resp:  res,
			}
			pubnub.Config.Log.Printf("Process: send on channel")
			job.JobResponse <- jqr
		}
	}()
}

func (p *RequestWorkers) Start(pubnub *PubNub) {
	pubnub.Config.Log.Printf("Start: Running with workers %d", p.MaxWorkers)
	for i := 0; i < p.MaxWorkers; i++ {
		pubnub.Config.Log.Printf("Start: StartNonSubWorker %d", i)
		worker := NewRequestWorkers(p.Workers, i)
		worker.Process(pubnub)
	}
	go p.ReadQueue(pubnub)
}

func (p *RequestWorkers) ReadQueue(pubnub *PubNub) {
	for job := range pubnub.jobQueue {
		pubnub.Config.Log.Printf("ReadQueue: Got job for channel %s ", job.Req)
		go func(job *JobQItem) {
			jobChannel := <-p.Workers
			jobChannel <- job
		}(job)
	}
	pubnub.Config.Log.Printf("ReadQueue: Exit")
}

func (p *RequestWorkers) Close() {
	close(p.Workers)
}
