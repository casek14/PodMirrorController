package controller

import (
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listcorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"log"
)

const (
	SOURCEPODNAMESPACE = ""
	SYNCLABEL          = ""
)

// Pod controller represents controller which watch pods in given namespace and distribute them
// to another namespaces based on label,
//
type PodController struct {
	podGetter       corev1.PodsGetter
	podLister       listcorev1.PodLister
	podListerSynced cache.InformerSynced
}

// Create new PodController
func NewPodController(client *kubernetes.Clientset,
	podInformer informercorev1.PodInformer) *PodController {
	c := PodController{
		podGetter:       client.CoreV1(),
		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,
	}

	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				c.onAdd(obj)
				log.Printf("Pod added !\n")
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				c.onUpdate(oldObj, newObj)
				log.Printf("Pod updated !\n")
			},
			DeleteFunc: func(obj interface{}) {
				c.onDelete(obj)
				log.Printf("Pod deleted !\n")
			},
		},
	)

	return &c
}

func (c *PodController) Run(stop chan struct{}) {
	log.Println("Waiting for cache to sync !!!")
	if !cache.WaitForCacheSync(stop, c.podListerSynced) {
		log.Println("Timeout waiting for cache sync")
		return
	}
	log.Println("Caches are synced !!!")
	log.Println("Waiting for stop signal !!!")
	<-stop
	log.Println("Received stop signal")
}

// Run action onAdd event
func (c *PodController) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Printf("onAdd: Unable to get key for #%v", obj, err)
	}
	log.Printf("onAdd: %v", key)
}

// Run action onUpdate event
func (c *PodController) onUpdate(oldObj, _ interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(oldObj)
	if err != nil {
		log.Printf("onUpdate: Unable to get key for #%v", oldObj, err)
	}
	log.Printf("onAdd: %v", key)
}

// Run action onDelete event
func (c *PodController) onDelete(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Printf("onDelete: Unable to get key for #%v", obj, err)
	}
	log.Printf("onDelete: %v", key)
}
