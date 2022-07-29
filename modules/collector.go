package modules

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetNamespace(clientset *kubernetes.Clientset, namespaceSeselector string, ignorenamespaces []string) []string {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: namespaceSeselector,
	})
	if err != nil {
		panic(err.Error())
	}
	var validnamespaces []string
	for index := range namespaces.Items {
		var flag bool = true
		namespace := namespaces.Items[index].Name
		for _, ns := range ignorenamespaces {
			if ns == namespace {
				flag = false
			}
		}
		if flag == true {
			validnamespaces = append(validnamespaces, namespace)
		}
	}
	return validnamespaces
}

func GetMetric(validnamespaces []string, clientset *kubernetes.Clientset, clientmetrics *metricsv.Clientset) [250]PodMetric {
	var podmetric [250]PodMetric
	var podindex int = 0
	for _, namespace := range validnamespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		var k int = 0
		for _, v := range pods.Items {
			if v.Status.Phase == "Running" {
				podMetrics, err := clientmetrics.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}
				metrics := podMetrics.Items[k]
				podmetric[podindex].PodName = podMetrics.Items[k].Name
				podmetric[podindex].Namespace = podMetrics.Items[k].Namespace
				for j := range metrics.Containers {
					if metrics.Containers[j].Name != "envoy" {
						podmetric[podindex].CurrentCpu = int(metrics.Containers[j].Usage.Cpu().MilliValue())
						podmetric[podindex].CurrentMem = int(metrics.Containers[j].Usage.Memory().Value() / 1048576)
					}
				}
				podindex++
				k++
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	return podmetric
}
