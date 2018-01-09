package cron

import (
	"github.com/fwtpe/owl/modules/sender/g"
	"github.com/fwtpe/owl/modules/sender/model"
	"github.com/fwtpe/owl/modules/sender/proc"
	"github.com/fwtpe/owl/modules/sender/redis"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"
	"time"
)

func ConsumeSms() {
	queue := g.Config().Queue.Sms
	for {
		L := redis.PopAllSms(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*model.Sms) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

func SendSms(sms *model.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config().Api.Sms
	r := httplib.Post(url).SetTimeout(5*time.Second, 2*time.Minute)
	r.Param("tos", sms.Tos)
	r.Param("content", sms.Content)
	resp, err := r.String()
	if err != nil {
		log.Println(err)
	}

	proc.IncreSmsCount()

	if g.Config().Debug {
		log.Println("==sms==>>>>", sms)
		log.Println("<<<<==sms==", resp)
	}

}
