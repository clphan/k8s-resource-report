package modules

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"
	"time"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
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
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
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
	clientset     *kubernetes.Clientset
	clientmetrics *metricsv.Clientset
}

func GetClient() (*kubernetes.Clientset, *metricsv.Clientset) {
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
	metricset, err := metricsv.NewForConfig(config)
	return clientset, metricset
}

// func GetNamespace(clientset *kubernetes.Clientset, namespaceSeselector string, ignorenamespaces []string) []string {
// 	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
// 		LabelSelector: namespaceSeselector,
// 	})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var validnamespaces []string
// 	for index := range namespaces.Items {
// 		var flag bool = true
// 		namespace := namespaces.Items[index].Name
// 		for _, ns := range ignorenamespaces {
// 			if ns == namespace {
// 				flag = false
// 			}
// 		}
// 		if flag == true {
// 			validnamespaces = append(validnamespaces, namespace)
// 		}
// 	}
// 	return validnamespaces
// }

func GetMetricClientApi(namespace string, podname string, clientset *kubernetes.Clientset) {
	var pods *PodMetricsList
	// var apipath string
	// apipath = "apis/metrics.k8s.io/v1beta1/" + namespace + "pods" + podname
	data, err := clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/pods").DoRaw(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(data, &pods)
	// rs := json.Unmarshal()
}

// func GetMetric(validnamespaces []string, clientset *kubernetes.Clientset, clientmetrics *metricsv.Clientset, timeInterval time.Duration) []PodMetric {
// 	var podmetric []PodMetric
// 	for _, namespace := range validnamespaces {
// 		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		var k int = 0
// 		for _, v := range pods.Items {
// 			if v.Status.Phase == "Running" {
// 				podMetrics, err := clientmetrics.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
// 				if err != nil {
// 					panic(err.Error())
// 				}
// 				metrics := podMetrics.Items[k]
// 				var currentcpu, currentmem int
// 				for j := range metrics.Containers {
// 					if metrics.Containers[j].Name != "envoy" {
// 						currentcpu = int(metrics.Containers[j].Usage.Cpu().MilliValue())
// 						currentmem = int(metrics.Containers[j].Usage.Memory().Value() / 1048576)
// 					}
// 				}
// 				podmetric = append(podmetric, PodMetric{podMetrics.Items[k].Namespace, podMetrics.Items[k].Name, currentcpu, currentmem})
// 				k++
// 			}
// 			time.Sleep(timeInterval * time.Millisecond)
// 		}
// 	}
// 	return podmetric
// }
