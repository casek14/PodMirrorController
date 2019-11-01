package main

import (
	"fmt"
	"github.com/casek14/PodMirrorController/controller"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {

	//	kubeconfig := ""

	config, err := clientcmd.BuildConfigFromFlags("", "/home/casek/.kube/config")
	if err != nil {
		fmt.Printf("Unable to load config file %v", err)
	}
	client := kubernetes.NewForConfigOrDie(config)
	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)
	podController := controller.NewPodController(client, sharedInformers.Core().V1().Pods())
	sharedInformers.Start(nil)
	podController.Run(nil)
}
