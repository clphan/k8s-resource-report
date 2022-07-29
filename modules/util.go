package main

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

type podMetric struct {
	namespace  string
	podname    string
	currentmem int
	currentcpu int
}

type clientConfig struct {
	clientset     *kubernetes.Clientset
	clientmetrics *metricsv.Clientset
}

func getclient() (*kubernetes.Clientset, *metricsv.Clientset) {
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

func getnamespace(clientset *kubernetes.Clientset, namespaceSeselector string, ignorenamespaces []string) []string {
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

func getmetric(validnamespaces []string, clientset *kubernetes.Clientset, clientmetrics *metricsv.Clientset) []podMetric {
	var podmetric []podMetric
	var podindex int = 0
	for _, namespace := range validnamespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for i, v := range pods.Items {
			if v.Status.Phase == "Running" {
				podmetric[podindex].namespace = namespace
				podmetric[podindex].podname = pods.Items[i].Name
				podMetrics, err := clientmetrics.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}
				metrics := podMetrics.Items[i]
				for j := range metrics.Containers {
					if metrics.Containers[j].Name != "envoy" {
						podmetric[podindex].currentcpu = int(metrics.Containers[j].Usage.Cpu().MilliValue())
						podmetric[podindex].currentmem = int(metrics.Containers[j].Usage.Memory().Value() / 1048576)
					}
				}
				podindex++
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	return podmetric
}