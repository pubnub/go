package pubnub

import "net/http"

type nonSubMsgType int

const (
	messageTypePublish nonSubMsgType = 1 << iota
	messageTypePAM
)

// JobQResponse is the type to store the resposne and error of the requests in the queue.
type JobQResponse struct {
	Resp  *http.Response
	Error error
}

// JobQItem is the type to store the request, client and its resposne.
type JobQItem struct {
	Req         *http.Request
	Client      *http.Client
	JobResponse chan *JobQResponse
}

// RequestWorkers is the type to store the workers info
type RequestWorkers struct {
	Workers        []Worker
	WorkersChannel chan chan *JobQItem
	MaxWorkers     int
	Sem            chan bool
}

// Worker is the type to store the worker info
type Worker struct {
	WorkersChannel chan chan *JobQItem
	JobChannel     chan *JobQItem
	ctx            Context
	id             int
}

func newRequestWorkers(workers chan chan *JobQItem, id int, ctx Context) Worker {
	return Worker{
		WorkersChannel: workers,
		JobChannel:     make(chan *JobQItem),
		ctx:            ctx,
		id:             id,
	}
}

// Process runs a goroutine for the worker
func (pw Worker) Process(pubnub *PubNub) {
	go func() {
	ProcessLabel:
		for {
			select {
			case pw.WorkersChannel <- pw.JobChannel:
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
				pubnub.Config.Log.Println("Exiting Worker Process by worker ctx, id ", pw.id)
				break ProcessLabel
			case <-pubnub.ctx.Done():
				pubnub.Config.Log.Println("Exiting Worker Process by PN ctx, id ", pw.id)
				break ProcessLabel
			}
		}
	}()
}

// Start starts the workers
func (p *RequestWorkers) Start(pubnub *PubNub, ctx Context) {
	pubnub.Config.Log.Println("Start: Running with workers ", p.MaxWorkers)
	p.Workers = make([]Worker, p.MaxWorkers)
	for i := 0; i < p.MaxWorkers; i++ {
		pubnub.Config.Log.Println("Start: StartNonSubWorker ", i)
		worker := newRequestWorkers(p.WorkersChannel, i, ctx)
		worker.Process(pubnub)
		p.Workers[i] = worker
	}
	go p.ReadQueue(pubnub)
}

// ReadQueue reads the queue and passes on the job to the workers
func (p *RequestWorkers) ReadQueue(pubnub *PubNub) {
	for job := range pubnub.jobQueue {
		pubnub.Config.Log.Println("ReadQueue: Got job for channel ", job.Req)
		go func(job *JobQItem) {
			jobChannel := <-p.WorkersChannel
			jobChannel <- job
		}(job)
	}
	pubnub.Config.Log.Println("ReadQueue: Exit")
}

// Close closes the workers
func (p *RequestWorkers) Close() {

	for _, w := range p.Workers {
		close(w.JobChannel)
		w.ctx.Done()
	}
}
