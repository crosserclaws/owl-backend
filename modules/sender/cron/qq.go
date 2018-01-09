package cron

import (
	"github.com/fwtpe/owl/modules/sender/g"
	"github.com/fwtpe/owl/modules/sender/model"
	"github.com/fwtpe/owl/modules/sender/proc"
	"github.com/fwtpe/owl/modules/sender/redis"
	//"github.com/toolkits/net/httplib"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

func ConsumeQQ() {
	queue := g.Config().Queue.QQ
	for {
		L := redis.PopAllQQ(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendQQList(L)
	}
}

func SendQQList(L []*model.QQ) {
	for _, qq := range L {
		QQWorkerChan <- 1
		go SendQQ(qq)
	}
}

func SendQQ(qq *model.QQ) {
	defer func() {
		<-QQWorkerChan
	}()

	url := g.Config().Api.QQ
	cmd := exec.Command("/bin/bash", "./qq_sms.sh", url, qq.Subject, qq.Content)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	proc.IncreQQCount()

	if g.Config().Debug {
		log.Println("==qq==>>>>", qq.Subject)
	}

}
