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
}

type PublishQueueProcessor struct {
	Workers    chan chan PublishJob
	maxWorkers int
	Sem        chan bool
}

func NewPublishWorker(workers chan chan PublishJob) PublishWorker {
	return PublishWorker{
		Workers:    workers,
		JobChannel: make(chan PublishJob),
	}
}

func (pw PublishWorker) StartPublishWorker(pubnub *Pubnub) {
	go func() {
		for {
			pw.Workers <- pw.JobChannel
			pubnub.infoLogger.Println("INFO: StartPublishWorker: Got job")
			select {
			case publishJob := <-pw.JobChannel:
				pubnub.infoLogger.Println("INFO: StartPublishWorker processing job FOR CHANNEL %s: Got job %d", publishJob.Channel, publishJob.PublishURL)
				value, responseCode, err := pubnub.publishHTTPRequest(publishJob.PublishURL)
				pubnub.readPublishResponseAndCallSendResponse(publishJob.Channel, value, responseCode, err, publishJob.CallbackChannel, publishJob.ErrorChannel)
			}
		}
	}()
}

func (pubnub *Pubnub) newPublishQueueProcessor(maxWorkers int) {
	//func (pubnub *Pubnub) newPublishQueueProcessor(maxWorkers int) *PublishQueueProcessor {
	workers := make(chan chan PublishJob, maxWorkers)
	sem := make(chan bool, maxWorkers)
	pubnub.infoLogger.Println("INFO: Init PublishQueueProcessor: workers %d", maxWorkers)

	p := &PublishQueueProcessor{
		Workers:    workers,
		maxWorkers: maxWorkers,
		Sem:        sem,
	}
	p.Run(pubnub)
	//go p.process(pubnub)
	//return p
}

func (p *PublishQueueProcessor) Run(pubnub *Pubnub) {
	//func (p *PublishQueueProcessor) Run(pubnub *Pubnub, publishJob PublishJob) {
	pubnub.infoLogger.Println("INFO: PublishQueueProcessor: Running with workers %d", p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		pubnub.infoLogger.Println("INFO: PublishQueueProcessor: StartPublishWorker %d", i)
		publishWorker := NewPublishWorker(p.Workers)
		publishWorker.StartPublishWorker(pubnub)
	}
	go p.process(pubnub)
	/*p.Sem <- true
	go func(publishJob PublishJob) {
		pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: Got job %d", publishJob.PublishURL)
		defer func() { <-p.Sem }()
		// get the url
		pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: Running job %d", publishJob.PublishURL)
		value, responseCode, err := pubnub.publishHTTPRequest(publishJob.PublishURL)
		pubnub.readPublishResponseAndCallSendResponse(publishJob.Channel, value, responseCode, err, publishJob.CallbackChannel, publishJob.ErrorChannel)

	}(publishJob)
	for i := 0; i < cap(p.Sem); i++ {
		p.Sem <- true
	}*/
}

func (p *PublishQueueProcessor) process(pubnub *Pubnub) {
	for {
		select {
		case publishJob := <-pubnub.publishJobQueue:
			pubnub.infoLogger.Println("INFO: PublishQueueProcessor process: Got job %d", publishJob.PublishURL)
			go func(publishJob PublishJob) {
				jobChannel := <-p.Workers
				jobChannel <- publishJob
			}(publishJob)
			/*pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: Got job, check sem %d len:%d", publishJob.PublishURL, len(pubnub.publishJobQueue))
			p.Sem <- true
			go func(publishJob PublishJob) {
				defer func() {
					pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: Defer job %d", publishJob.PublishURL)
					b := <-p.Sem
					pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: After Defer job %d", b)
				}()
				// get the url
				pubnub.infoLogger.Println("INFO: StartPublishWorker processing job: Running job %d", publishJob.PublishURL)
				value, responseCode, err := pubnub.publishHTTPRequest(publishJob.PublishURL)
				pubnub.readPublishResponseAndCallSendResponse(publishJob.Channel, value, responseCode, err, publishJob.CallbackChannel, publishJob.ErrorChannel)

			}(publishJob)
			/*for i := 0; i < cap(p.Sem); i++ {
				p.Sem <- true
			}*/
		}
	}
}

func (p *PublishQueueProcessor) Close(pubnub *Pubnub) {
	close(p.Workers)
}
