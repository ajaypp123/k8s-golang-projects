package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var currentObserver *ConfigObserver

func StartController(clientset *kubernetes.Clientset, dynClient dynamic.Interface, stopCh chan struct{}) {
	fmt.Println("Starting controller to watch over deployment")

	factory := informers.NewSharedInformerFactory(clientset, 0)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	gvr := schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1",
		Resource: "configobservers",
	}

	observerInformer := cache.NewSharedInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (pkgruntime.Object, error) {
				return dynClient.Resource(gvr).Namespace("testing").List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return dynClient.Resource(gvr).Namespace("testing").Watch(context.TODO(), options)
			},
		},
		&unstructured.Unstructured{},
		0,
	)

	observerInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// If having multiple CR for CRD it should be handled here, by storing each CR copy by AddFunc
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			observer, err := ParseUnstructuredToObserver(u)
			if err != nil {
				fmt.Println("Failed to parse observer:", err)
				return
			}
			currentObserver = observer
			fmt.Println("Observer CR created:", observer.Spec)
			_ = Reconcile(clientset, observer)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			u := newObj.(*unstructured.Unstructured)
			observer, err := ParseUnstructuredToObserver(u)
			if err != nil {
				fmt.Println("Failed to parse observer:", err)
				return
			}
			currentObserver = observer
			fmt.Println("Observer CR updated:", observer.Spec)
			_ = Reconcile(clientset, observer)
		},
	})

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			if currentObserver == nil {
				return
			}
			deploy := newObj.(*appsv1.Deployment)

			if deploy.Name == currentObserver.Spec.DeploymentName {
				_ = Reconcile(clientset, currentObserver)
			}
		},
	})

	go observerInformer.Run(stopCh)
	go deploymentInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, observerInformer.HasSynced, deploymentInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	<-stopCh
}

func updateConfigMap(clientset *kubernetes.Clientset, namespace, name, key, value string) error {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get configmap error: %v", err)
	}
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	cm.Data[key] = value
	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return err
}
