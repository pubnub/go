package messaging

import (
//"time"
)

type PublishJob struct {
	Channel         string
	PublishURL      string
	CallbackChannel chan []byte
	ErrorChannel    chan []byte
}

type PublishWorker struct {
	Workers    chan chan PublishJob
	JobChannel chan PublishJob
	quit       chan bool
	id         int
}

type PublishQueueProcessor struct {
	Workers    chan chan PublishJob
	maxWorkers int
	Sem        chan bool
}

func NewPublishWorker(workers chan chan PublishJob, id int) PublishWorker {
	return PublishWorker{
		Workers:    workers,
		JobChannel: make(chan PublishJob),
		id:         id,
	}
}

func (pw PublishWorker) StartPublishWorker(pubnub *Pubnub) {
	go func() {
		for {
			pw.Workers <- pw.JobChannel
			pubnub.infoLogger.Printf("INFO: StartPublishWorker: Got job", pw.id)
			select {
			case publishJob := <-pw.JobChannel:

				pubnub.infoLogger.Printf("INFO: StartPublishWorker processing job FOR CHANNEL %s: Got job %s, id:%d", publishJob.Channel, publishJob.PublishURL, pw.id)
				pn := pubnub
				value, responseCode, err := pn.publishHTTPRequest(publishJob.PublishURL)
				pubnub.readPublishResponseAndCallSendResponse(publishJob.Channel, value, responseCode, err, publishJob.CallbackChannel, publishJob.ErrorChannel)
			}
		}
	}()
}

func (pubnub *Pubnub) newPublishQueueProcessor(maxWorkers int) {
	//logic 1
	workers := make(chan chan PublishJob, maxWorkers)
	//end logic 1

	//logic 2
	//sem := make(chan bool, maxWorkers)
	//end logic 2

	pubnub.infoLogger.Printf("INFO: Init PublishQueueProcessor: workers %d", maxWorkers)

	p := &PublishQueueProcessor{
		//logic 1
		Workers: workers,
		//end logic 1
		maxWorkers: maxWorkers,
		//logic 2
		//Sem: sem,
		//end logic 2
	}
	p.Run(pubnub)
}

func (p *PublishQueueProcessor) Run(pubnub *Pubnub) {
	pubnub.infoLogger.Printf("INFO: PublishQueueProcessor: Running with workers %d", p.maxWorkers)
	//logic 1
	for i := 0; i < p.maxWorkers; i++ {
		pubnub.infoLogger.Printf("INFO: PublishQueueProcessor: StartPublishWorker %d", i)
		publishWorker := NewPublishWorker(p.Workers, i)
		publishWorker.StartPublishWorker(pubnub)
	}
	//end logic 1
	go p.process(pubnub)
}

func (p *PublishQueueProcessor) process(pubnub *Pubnub) {
	for {
		select {
		case publishJob := <-pubnub.publishJobQueue:
			pubnub.infoLogger.Printf("INFO: PublishQueueProcessor process: Got job for channel %s %s", publishJob.Channel, publishJob.PublishURL)
			//logic 2
			//p.Sem <- true
			//end logic 2
			go func(publishJob PublishJob) {
				//logic 1
				jobChannel := <-p.Workers

				jobChannel <- publishJob
				//end logic 1

				//logic 2
				/*defer func() {
					pubnub.infoLogger.Printf("INFO: StartPublishWorker processing job: Defer job %d", publishJob.PublishURL)
					b := <-p.Sem
					pubnub.infoLogger.Printf("INFO: StartPublishWorker processing job: After Defer job %d", b)
				}()

				pubnub.infoLogger.Printf("INFO: StartPublishWorker processing job FOR CHANNEL %s: Got job %d", publishJob.Channel, publishJob.PublishURL)
				pn := pubnub
				value, responseCode, err := pn.publishHTTPRequest(publishJob.PublishURL)
				pubnub.readPublishResponseAndCallSendResponse(publishJob.Channel, value, responseCode, err, publishJob.CallbackChannel, publishJob.ErrorChannel)
				*/
				//end logic 2

			}(publishJob)
		}
	}
}

func (p *PublishQueueProcessor) Close(pubnub *Pubnub) {
	close(p.Workers)
}
