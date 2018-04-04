package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/mvenezia/redis-operator/pkg/client/clientset/versioned"
	"gitlab.com/mvenezia/redis-operator/pkg/controller"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/juju/loggo"
	"github.com/venezia/redis-operator/pkg/util"
)

var (
	logger loggo.Logger
	config *rest.Config
)

const (
	kubeconfigDir  = ".kube"
	kubeconfigFile = "config"
)

func main() {
	var err error
	logger := util.GetModuleLogger("cmd.redis-operator", loggo.INFO)

	// setup viper
	viperInit()

	// get flags
	portNumber := viper.GetInt("port")
	kubeconfigLocation := viper.GetString("kubeconfig")

	// Debug for now
	logger.Infof("Parsed Variables: \n  Port: %d \n  Kubeconfig: %s", portNumber, kubeconfigLocation)

	if kubeconfigLocation != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigLocation)
		if err != nil {
			logErrorAndExit(err)
		}
	} else {
		configPath := filepath.Join(homeDir(), kubeconfigDir, kubeconfigFile)
		if _, err := os.Stat(configPath); err == nil {
			config, err = clientcmd.BuildConfigFromFlags("", configPath)
		} else {
			config, err = rest.InClusterConfig()
		}
	}

	// create the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logErrorAndExit(err)
	}

	operatorController := controller.New(controller.Config{
		KubeCli:    clientSet,
		KubeExtCli: apiextensionsclient.NewForConfigOrDie(config),
		RedisCRCli: versioned.NewForConfigOrDie(config),
	})

	operatorController.InitCRD()
	operatorController.Start()

	monitorPods(clientSet, "default", "example-xxxxx")
}

func viperInit() {
	viper.SetEnvPrefix("redisoperator")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// using standard library "flag" package
	flag.Int("port", 8081, "Port to listen on")
	flag.String("kubeconfig", "", "Location of kubeconfig file")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()
}

func monitorPods(clientSet *kubernetes.Clientset, namespace string, pod string) {
	for {
		pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			logErrorAndExit(err)
		}

		logger.Infof("There are %d pods in the cluster", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientSet.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})

		if errors.IsNotFound(err) {
			logger.Warningf("Pod %s in namespace %s not found", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			logger.Warningf("Error getting pod %s in namespace %s: %v", pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			logErrorAndExit(err)
		} else {
			logger.Infof("Found pod %s in namespace %s", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}

func logErrorAndExit(err error) {
	logger.Criticalf("error: %s", err)
	os.Exit(1)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
