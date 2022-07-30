package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/clphan/k8s-resource-report/modules"
)

type podMetric struct {
	namespace  string
	podname    string
	currentmem int
	currentcpu int
}

func main() {
	label := os.Args[1]
	var ignorenamespaces []string = []string{"abc", "cdb"}
	clientset, clientmetrics := modules.GetClient()
	validnamespaces := modules.GetNamespace(clientset, label, ignorenamespaces)
	podmetrics := modules.GetMetric(validnamespaces, clientset, clientmetrics, 100)
	fmt.Println("Num object:", len(podmetrics))
	for i := range podmetrics {
		fmt.Println(podmetrics[i])
	}
	var csvdata [][]string
	for i := range podmetrics {
		csvdata = append(csvdata, []string{podmetrics[i].Namespace, podmetrics[i].PodName, strconv.Itoa(podmetrics[i].CurrentCpu), strconv.Itoa(podmetrics[i].CurrentMem)})
	}
	fmt.Println(csvdata)
}
