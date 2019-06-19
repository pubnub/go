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
	ctx        Context
	id         int
}

var Workers []Worker

func NewRequestWorkers(workers chan chan *JobQItem, id int, ctx Context) Worker {
	return Worker{
		Workers:    workers,
		JobChannel: make(chan *JobQItem),
		ctx:        ctx,
		id:         id,
	}
}

func (pw Worker) Process(pubnub *PubNub) {
	go func() {
	B:
		for {
			select {
			case pw.Workers <- pw.JobChannel:
				job := <-pw.JobChannel

				if job != nil {

					res, err := job.Client.Do(job.Req)
					jqr := &JobQResponse{
						Error: err,
						Resp:  res,
					}
					job.JobResponse <- jqr
					pubnub.Config.Log.Println("Request sent using worker id ", pw.id)
				}
			case <-pw.ctx.Done():
				pubnub.Config.Log.Println("Exiting Worker Process id ", pw.id)
				break B
			case <-pubnub.ctx.Done():
				pubnub.Config.Log.Println("2 Exiting Worker Process id ", pw.id)
				break B
			}
		}
	}()
}

func (p *RequestWorkers) Start(pubnub *PubNub, ctx Context) {
	pubnub.Config.Log.Println("Start: Running with workers ", p.MaxWorkers)
	Workers = make([]Worker, p.MaxWorkers)
	for i := 0; i < p.MaxWorkers; i++ {
		pubnub.Config.Log.Println("Start: StartNonSubWorker ", i)
		worker := NewRequestWorkers(p.Workers, i, ctx)
		worker.Process(pubnub)
		Workers[i] = worker
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

	for _, w := range Workers {
		close(w.JobChannel)
		w.ctx.Done()
	}
}
