package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func terminateNode(node string, session *session.Session) error {

	split := strings.SplitAfter(node, "/")
	node = split[len(split)-1]

	svc := ec2.New(session)

	log.Infof("terminating node: %s", node)

	_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{&node},
	})

	if err != nil {
		return err
	}

	log.Infof("node terminated: %s", node)

	return nil
}
