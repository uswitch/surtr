package main

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"k8s.io/client-go/kubernetes"
)

type options struct {
	kubeconfig string
	olderThan  time.Duration
	debug      bool
}

func main() {
	ctx := context.Background()

	opts := &options{}
	kingpin.Flag("kubeconfig", "Path to kubeconfig.").StringVar(&opts.kubeconfig)
	kingpin.Flag("older-than", "age of nodes to terminate").Required().DurationVar(&opts.olderThan)
	kingpin.Flag("debug", "Debug mode").BoolVar(&opts.debug)
	kingpin.Parse()
	log.SetOutput(os.Stderr)

	if opts.debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Debug logging enabled")
	}

	config, err := createClientConfig(opts.kubeconfig)
	if err != nil {
		log.Fatalf("error creating kube config: %s", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error creating kube client: %s", err)
	}

	node, err := getNode(ctx, client, opts.olderThan)
	if err != nil {
		log.Fatalf("failed to get node: %s", err)
	}
	if node == "" {
		log.Info("no nodes to terminate, shutting down")
		os.Exit(0)
	}

	session, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create AWS session: %s", err)
	}

	err = terminateNode(node, session)
	if err != nil {
		log.Fatalf("error terminating node: %s", err)
	}

}
