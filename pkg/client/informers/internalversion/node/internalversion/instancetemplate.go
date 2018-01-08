// This file was automatically generated by informer-gen

package internalversion

import (
	node "github.com/gardener/node-controller-manager/pkg/apis/node"
	internalinterfaces "github.com/gardener/node-controller-manager/pkg/client/informers/internalversion/internalinterfaces"
	internalclientset "github.com/gardener/node-controller-manager/pkg/client/internalclientset"
	internalversion "github.com/gardener/node-controller-manager/pkg/client/listers/node/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// InstanceTemplateInformer provides access to a shared informer and lister for
// InstanceTemplates.
type InstanceTemplateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.InstanceTemplateLister
}

type instanceTemplateInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

// NewInstanceTemplateInformer constructs a new informer for InstanceTemplate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewInstanceTemplateInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.Node().InstanceTemplates(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.Node().InstanceTemplates(namespace).Watch(options)
			},
		},
		&node.InstanceTemplate{},
		resyncPeriod,
		indexers,
	)
}

func defaultInstanceTemplateInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewInstanceTemplateInformer(client, v1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *instanceTemplateInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&node.InstanceTemplate{}, defaultInstanceTemplateInformer)
}

func (f *instanceTemplateInformer) Lister() internalversion.InstanceTemplateLister {
	return internalversion.NewInstanceTemplateLister(f.Informer().GetIndexer())
}