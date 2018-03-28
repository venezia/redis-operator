package controller

import (
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	api "gitlab.com/mvenezia/redis-operator/pkg/apis/redis/v1alpha1"

	"github.com/sirupsen/logrus"
	"gitlab.com/mvenezia/redis-operator/pkg/util/k8sutil"
	"gitlab.com/mvenezia/redis-operator/pkg/redis"

	"fmt"
)

type Controller struct {
	logger *logrus.Entry
	Config

}

type Config struct {
	Namespace 		string
	ClusterWide 	bool
	ServiceAccount 	string
	KubeCli			kubernetes.Interface
	KubeExtCli		apiextensionsclient.Interface
	CreateCRD		bool
}

func New(cfg Config) *Controller {
	return &Controller{
		logger: logrus.WithField("pkg", "controller"),
		Config: cfg,
	}
}

func (c *Controller) makeClusterConfig() redis.Config {
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