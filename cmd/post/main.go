package main

import (
	"2019_2_Next_Level/internal/logger"
	"2019_2_Next_Level/internal/post"
	mailsender "2019_2_Next_Level/internal/post/MailSender"
	"2019_2_Next_Level/internal/post/messagequeue"
	"2019_2_Next_Level/internal/post/smtpd"
	"2019_2_Next_Level/pkg/config"
	"flag"
	"sync"
)

const (
	configFilenameDefault = "post_service.config.json"
)

type daemon interface {
	Init(post.ChanPair, post.ChanPair) error
	Run(*sync.WaitGroup)
}

var log logger.Log

func main() {
	// Должны быть компоненты:
	// 	* Очередь исходящих
	// 	* Очередь Входящих
	// 	* Отправщик
	// 	* SMTP-сервер

	// 	outcomingQueue <------> mailsender <--------> smtpd <--------> incomingQueue
	log.SetPrefix("PostServerMain")
	err := initializeConfig()
	if err != nil {
		log.Println(err)
		return
	}

	daemonList := []daemon{
		&messagequeue.QueueDemon{Name: "outcoming"},
		&mailsender.MailSender{},
		&smtpd.Server{},
		&messagequeue.QueueDemon{Name: "incoming"},
	}
	Execute(daemonList...)
}

// Execute : starts daemon chain
func Execute(daemons ...daemon) {
	wg := &sync.WaitGroup{}

	previous := post.ChanPair{}.Init(post.Conf.ChannelCapasity)

	for _, daemon := range daemons {
		next := post.ChanPair{}.Init(post.Conf.ChannelCapasity)
		err := daemon.Init(previous, next)
		if err != nil {
			log.Println("Error during initializing a daemon: ", err)
			return
		}
		previous = next
		wg.Add(1)
		go daemon.Run(wg)
	}
	wg.Wait()
}

func initializeConfig() error {
	configFilename := flag.String("config", configFilenameDefault, "Path to config file")
	flag.Parse()

	return config.Configurator.Inflate(*configFilename, &post.Conf)
}
