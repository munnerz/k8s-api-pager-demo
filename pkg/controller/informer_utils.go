package controller

import (
	"time"

	factory_interfaces "github.com/srossross/k8s-test-controller/pkg/informers/externalversions/internalinterfaces"

	client "github.com/srossross/k8s-test-controller/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgRuntime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

func (ctrl *TestController) newPodInformerFactory() factory_interfaces.NewInformerFunc {
	return func(cl client.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		sharedIndexInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (pkgRuntime.Object, error) {
					return ctrl.CoreV1().Pods(metav1.NamespaceAll).List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return ctrl.CoreV1().Pods(metav1.NamespaceAll).Watch(options)
				},
			},
			&corev1.Pod{},
			resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
		)
		return sharedIndexInformer
	}
}
