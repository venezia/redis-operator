package controller

import (
	api "gitlab.com/mvenezia/redis-operator/pkg/apis/redis/v1alpha1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	"github.com/sirupsen/logrus"
	"gitlab.com/mvenezia/redis-operator/pkg/redis"
	"gitlab.com/mvenezia/redis-operator/pkg/util/k8sutil"
	kwatch "k8s.io/apimachinery/pkg/watch"

	"fmt"
	"time"

	"context"
	"gitlab.com/mvenezia/redis-operator/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

var initRetryWaitTime = 30 * time.Second

var pt *panicTimer

func init() {
	pt = newPanicTimer(time.Minute, "unexpected long blocking (> 1 Minute) when handling cluster event")
}

type Event struct {
	Type   kwatch.EventType
	Object *api.Redis
}

type Controller struct {
	logger *logrus.Entry
	Config

	redii map[string]*redis.Redis
}

type Config struct {
	Namespace      string
	ClusterWide    bool
	ServiceAccount string
	KubeCli        kubernetes.Interface
	KubeExtCli     apiextensionsclient.Interface
	RedisCRCli     versioned.Interface
	CreateCRD      bool
}

func New(cfg Config) *Controller {
	return &Controller{
		logger: logrus.WithField("pkg", "controller"),
		Config: cfg,
		redii:  make(map[string]*redis.Redis),
	}
}

func (c *Controller) makeRedisConfig() redis.Config {
	return redis.Config{
		ServiceAccount: c.Config.ServiceAccount,
		KubeCli:        c.Config.KubeCli,
	}
}

func (c *Controller) InitCRD() error {
	err := k8sutil.CreateCRD(c.KubeExtCli, api.RedisCRDName, api.RedisResourceKind, api.RedisResourcePlural, "redis")
	if err != nil {
		return fmt.Errorf("failed to create CRD: %v", err)
	}
	return k8sutil.WaitCRDReady(c.KubeExtCli, api.RedisCRDName)
}

// handleRedisEvent returns true if redis is ignored (not managed) by this instance.
func (c *Controller) handleRedisEvent(event *Event) (bool, error) {
	clus := event.Object

	switch event.Type {
	case kwatch.Added:
		if _, ok := c.redii[clus.Name]; ok {
			return false, fmt.Errorf("unsafe state. cluster (%s) was created before but we received event (%s)", clus.Name, event.Type)
		}

		nc := redis.New(c.makeRedisConfig(), clus)

		c.redii[clus.Name] = nc

	case kwatch.Modified:
		if _, ok := c.redii[clus.Name]; !ok {
			return false, fmt.Errorf("unsafe state. cluster (%s) was never created but we received event (%s)", clus.Name, event.Type)
		}
		c.redii[clus.Name].Update(clus)

	case kwatch.Deleted:
		if _, ok := c.redii[clus.Name]; !ok {
			return false, fmt.Errorf("unsafe state. cluster (%s) was never created but we received event (%s)", clus.Name, event.Type)
		}
		c.redii[clus.Name].Delete(clus)
		delete(c.redii, clus.Name)
	}
	return false, nil
}

func (c *Controller) Start() error {
	// TODO: get rid of this init code. CRD and storage class will be managed outside of operator.
	for {
		err := c.initResource()
		if err == nil {
			break
		}
		c.logger.Errorf("initialization failed: %v", err)
		c.logger.Infof("retry in %v...", initRetryWaitTime)
		time.Sleep(initRetryWaitTime)
	}

	c.run()
	panic("unreachable")
}

func (c *Controller) run() {
	var ns string
	if c.Config.ClusterWide {
		ns = metav1.NamespaceAll
	} else {
		ns = c.Config.Namespace
	}

	source := cache.NewListWatchFromClient(
		c.Config.RedisCRCli.RedisV1alpha1().RESTClient(),
		api.RedisResourcePlural,
		ns,
		fields.Everything())

	_, informer := cache.NewIndexerInformer(source, &api.Redis{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAddRedis,
		UpdateFunc: c.onUpdateRedis,
		DeleteFunc: c.onDeleteRedis,
	}, cache.Indexers{})

	ctx := context.TODO()
	// TODO: use workqueue to avoid blocking
	informer.Run(ctx.Done())
}

func (c *Controller) initResource() error {
	if c.Config.CreateCRD {
		err := c.InitCRD()
		if err != nil {
			return fmt.Errorf("fail to init CRD: %v", err)
		}
	}
	return nil
}

func (c *Controller) onAddRedis(obj interface{}) {
	c.syncRedis(obj.(*api.Redis))
}

func (c *Controller) onUpdateRedis(oldObj, newObj interface{}) {
	c.syncRedis(newObj.(*api.Redis))
}

func (c *Controller) onDeleteRedis(obj interface{}) {
	clus, ok := obj.(*api.Redis)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			panic(fmt.Sprintf("unknown object from Redis delete event: %#v", obj))
		}
		clus, ok = tombstone.Obj.(*api.Redis)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not an Redis: %#v", obj))
		}
	}
	ev := &Event{
		Type:   kwatch.Deleted,
		Object: clus,
	}

	pt.start()
	_, err := c.handleRedisEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}

func (c *Controller) syncRedis(clus *api.Redis) {
	ev := &Event{
		Type:   kwatch.Added,
		Object: clus,
	}
	// re-watch or restart could give ADD event.
	// If for an ADD event the cluster spec is invalid then it is not added to the local cache
	// so modifying that cluster will result in another ADD event
	if _, ok := c.redii[clus.Name]; ok {
		ev.Type = kwatch.Modified
	}

	pt.start()
	_, err := c.handleRedisEvent(ev)
	if err != nil {
		c.logger.Warningf("fail to handle event: %v", err)
	}
	pt.stop()
}
