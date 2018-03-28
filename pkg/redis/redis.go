package redis

import (
	api "gitlab.com/mvenezia/redis-operator/pkg/apis/redis/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"gitlab.com/mvenezia/redis-operator/pkg/client/clientset/versioned"

	"time"
)

var (
	reconcileInterval         = 8 * time.Second
	podTerminationGracePeriod = int64(5)
)

type clusterEventType string

const (
	eventModifyCluster clusterEventType = "Modify"
)

type clusterEvent struct {
	typ     clusterEventType
	cluster *api.Redis
}

type Config struct {
	ServiceAccount string

	KubeCli   kubernetes.Interface
	RedisCRCli versioned.Interface
}

