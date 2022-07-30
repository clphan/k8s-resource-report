package modules

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PodMetric struct {
	Namespace  string
	PodName    string
	CurrentMem int
	CurrentCpu int
}

type PodMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp  time.Time `json:"timestamp"`
		Window     string    `json:"window"`
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	} `json:"items"`
}

type ClientConfig struct {
	clientset *kubernetes.Clientset
}

func GetClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset
}

func GetMetricClientApi(clientset *kubernetes.Clientset, pods *PodMetricsList, namespace) error {
	var path string = "apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "pods"
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/payment-proxy-management/pods").DoRaw(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	byteArr, _ := json.Unmarshal(data)
	fmt.Println(string(byteArr))
	err = json.Unmarshal(data, &pods)
	return err
}
