package main

import (
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getNode(client *kubernetes.Clientset, olderThan time.Duration) (string, error) {

	nodeGetter := client.CoreV1().Nodes()
	nodes, err := nodeGetter.List(metav1.ListOptions{})

	if err != nil {
		return "", err
	}
	return oldestNode(nodes.Items, olderThan), nil
}

func oldestNode(nodes []v1.Node, olderThan time.Duration) string {

	sortNodes(nodes)
	for _, node := range nodes {
		log.Debugf("node: %s, creation: %s", node.Name, node.CreationTimestamp)
	}

	if time.Now().Sub(nodes[0].CreationTimestamp.Time) >= olderThan {
		log.Infof("kubernetes node: %s will be terminated", nodes[0].Name)
		return nodes[0].Spec.ProviderID
	}

	return ""

}

func sortNodes(nodes []v1.Node) {
	sort.Slice(nodes, func(i int, j int) bool {
		return nodes[i].CreationTimestamp.Time.Before(nodes[j].CreationTimestamp.Time)
	})
}

func createClientConfig(kubeconfig string) (*rest.Config, error) {

	if kubeconfig != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("error creating kube config from config file: %s", err)
		}
		return config, nil
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating kube config from in cluster config: %s", err)

	}

	return config, nil

}
