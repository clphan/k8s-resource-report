package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/clphan/k8s-resource-report/modules"
)

type Containers struct {
	name       string
	currentmem string
	currentcpu string
}
type podMetric struct {
	namespace  string
	podname    string
	containers []Containers
}

func main() {
	label := os.Args[1]
	var ignorenamespaces []string = []string{""}
	clientset := modules.GetClient()
	var pods modules.PodMetricsList
	validnamespaces := modules.GetNamespace(clientset, label, ignorenamespaces)
	var podobject []podMetric
	for _, v := range validnamespaces {
		err := modules.GetMetricClientApi(clientset, &pods, v)
		if err != nil {
			panic(err.Error())
		}

		for _, m := range pods.Items {
			var containers []Containers
			for _, n := range m.Containers {
				containers = append(containers, Containers{n.Name, n.Usage.CPU, n.Usage.Memory})
			}
			podobject = append(podobject, podMetric{v, m.Metadata.Name, containers})
		}
	}
	var csvdata [][]string
	for i := range podobject {
		csvdata = append(csvdata, []string{podobject[i].Namespace, podobject[i].PodName, strconv.Itoa(podobject[i].CurrentCpu), podobject[i].CurrentMem)})
	}
	fmt.Println(csvdata)
}
