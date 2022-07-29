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
	podmetrics := modules.GetMetric(validnamespaces, clientset, clientmetrics)
	fmt.Println("Num object:", len(podmetrics))
	var csvdata [][]string
	for i := range podmetrics {
		csvdata[i][0] = podmetrics[i].namespace
		csvdata[i][1] = podmetrics[i].podname
		csvdata[i][2] = strconv.Itoa(podmetrics[i].currentcpu)
		csvdata[i][3] = strconv.Itoa(podmetrics[i].currentmem)
	}
}
