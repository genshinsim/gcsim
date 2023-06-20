package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/containrrr/shoutrrr"
	"github.com/genshinsim/gcsim/backend/pkg/notify"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
)

type service struct {
	c           *notify.Client
	infoUrl     string
	criticalUrl string
}

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

func run() error {
	var err error
	s := &service{}
	s.c, err = notify.New("notification-service")
	if err != nil {
		return err
	}
	s.infoUrl = fmt.Sprintf("discord://%v@%v", os.Getenv("NOTIFY_INFO_TOKEN"), os.Getenv("NOTIFY_INFO_ID"))
	s.criticalUrl = fmt.Sprintf("discord://%v@%v", os.Getenv("NOTIFY_CRITICAL_TOKEN"), os.Getenv("NOTIFY_CRITICAL_ID"))
	err = s.sub()
	if err != nil {
		return err
	}
	log.Println(s.infoUrl)
	log.Println(s.criticalUrl)

	return nil
}

func (s *service) sub() error {
	err := s.c.ActivateListener()
	if err != nil {
		return err
	}
	err = s.c.Subscribe(db.TopicReplace, s.onDBReplace)
	if err != nil {
		return err
	}
	err = s.c.Subscribe(db.TopicSubmissionDelete, s.onSubmissionDeleted)
	if err != nil {
		return err
	}
	err = s.c.Subscribe(db.TopicComputeCompleted, s.onComputeCompleted)
	if err != nil {
		return err
	}
	err = s.c.Subscribe(db.TopicSubmissionComputeFailed, s.onSubmissionComputeFailed)
	if err != nil {
		return err
	}
	err = s.c.Subscribe(db.TopicDBComputeFailed, s.onDBComputeFailed)
	if err != nil {
		return err
	}

	s.info("notification service now online")
	return nil
}

func (s *service) onDBReplace(topic string, payload []byte) {
	m := &model.EntryReplaceEvent{}
	err := protojson.Unmarshal(payload, m)
	if err != nil {
		log.Println("error marshalling event:", err)
		return
	}
	s.info(fmt.Sprintf("DB entry %v config has been replaced", m.DbId))
}

func (s *service) onSubmissionDeleted(topic string, payload []byte) {
	m := &model.SubmissionDeleteEvent{}
	err := protojson.Unmarshal(payload, m)
	if err != nil {
		log.Println("error marshalling event:", err)
		return
	}
	s.info(fmt.Sprintf("Submission %v has been deleted by original submitter", m.DbId))
}

func (s *service) onComputeCompleted(topic string, payload []byte) {
	m := &model.ComputeCompletedEvent{}
	err := protojson.Unmarshal(payload, m)
	if err != nil {
		log.Println("error marshalling event:", err)
		return
	}
	s.info(fmt.Sprintf("Compute for DB %v has been completed with share key %v. Link: https://gcsim.app/db/%v", m.DbId, m.ShareId, m.DbId))
}

func (s *service) onSubmissionComputeFailed(topic string, payload []byte) {
	m := &model.ComputeFailedEvent{}
	err := protojson.Unmarshal(payload, m)
	if err != nil {
		log.Println("error marshalling event:", err)
		return
	}
	s.info(fmt.Sprintf("Compute for submission with id %v has failed (%v); entry has been deleted", m.DbId, m.Reason))
}

func (s *service) onDBComputeFailed(topic string, payload []byte) {
	m := &model.ComputeFailedEvent{}
	err := protojson.Unmarshal(payload, m)
	if err != nil {
		log.Println("error marshalling event:", err)
		return
	}
	s.critical(fmt.Sprintf("Compute for db entry with id %v has failed (%v). Old link: https://gcsim.app/db/%v ", m.DbId, m.Reason, m.DbId))
}

func (s *service) info(msg string) {
	sender, err := shoutrrr.CreateSender(s.infoUrl)
	if err != nil {
		log.Println("creating info url sender failed:", err)
		return
	}
	sender.Send(msg, nil)
}

func (s *service) critical(msg string) {
	sender, err := shoutrrr.CreateSender(s.criticalUrl)
	if err != nil {
		log.Println("creating info url sender failed:", err)
		return
	}
	sender.Send(msg, nil)
}
