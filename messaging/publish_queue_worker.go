package messaging

type nonSubMsgType int

const (
	messageTypePublish nonSubMsgType = 1 << iota
	messageTypePAM
)

type NonSubJob struct {
	Channel         string
	NonSubURL       string
	NonSubMsgType   nonSubMsgType
	CallbackChannel chan []byte
	ErrorChannel    chan []byte
}

type NonSubWorker struct {
	Workers    chan chan NonSubJob
	JobChannel chan NonSubJob
	id         int
}

type NonSubQueueProcessor struct {
	Workers    chan chan NonSubJob
	maxWorkers int
	Sem        chan bool
}

func NewNonSubWorker(workers chan chan NonSubJob, id int) NonSubWorker {
	return NonSubWorker{
		Workers:    workers,
		JobChannel: make(chan NonSubJob),
		id:         id,
	}
}

func (pw NonSubWorker) StartNonSubWorker(pubnub *Pubnub) {
	go func() {
		for {
			pw.Workers <- pw.JobChannel
			pubnub.infoLogger.Printf("INFO: StartNonSubWorker: Worker started: %d", pw.id)
			nonSubJob := <-pw.JobChannel
			pubnub.infoLogger.Printf("INFO: StartNonSubWorker processing job FOR CHANNEL %s: Got job %s, id:%d", nonSubJob.Channel, nonSubJob.NonSubURL, pw.id)
			pn := pubnub
			value, responseCode, err := pn.nonSubHTTPRequest(nonSubJob.NonSubURL)
			if nonSubJob.NonSubMsgType == messageTypePublish {
				pubnub.readPublishResponseAndCallSendResponse(nonSubJob.Channel, value, responseCode, err, nonSubJob.CallbackChannel, nonSubJob.ErrorChannel)
			} else if nonSubJob.NonSubMsgType == messageTypePAM {
				pubnub.handlePAMResponse(nonSubJob.Channel, value, responseCode, err, nonSubJob.CallbackChannel, nonSubJob.ErrorChannel)
			}
		}
	}()
}

func (p *NonSubQueueProcessor) Run(pubnub *Pubnub) {
	pubnub.infoLogger.Printf("INFO: NonSubQueueProcessor: Running with workers %d", p.maxWorkers)
	//logic 1
	for i := 0; i < p.maxWorkers; i++ {
		pubnub.infoLogger.Printf("INFO: NonSubQueueProcessor: StartNonSubWorker %d", i)
		nonSubWorker := NewNonSubWorker(p.Workers, i)
		nonSubWorker.StartNonSubWorker(pubnub)
	}
	go p.process(pubnub)
}

func (p *NonSubQueueProcessor) process(pubnub *Pubnub) {
	for nonSubJob := range pubnub.nonSubJobQueue {
		pubnub.infoLogger.Printf("INFO: NonSubQueueProcessor process: Got job for channel %s %s", nonSubJob.Channel, nonSubJob.NonSubURL)
		go func(nonSubJob NonSubJob) {
			jobChannel := <-p.Workers
			jobChannel <- nonSubJob
		}(nonSubJob)
	}
}

func (p *NonSubQueueProcessor) Close() {
	close(p.Workers)
}
