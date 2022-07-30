package main

import (
	"fmt"
	"os"

	"github.com/clphan/k8s-resource-report/modules"
)

type podMetric struct {
	namespace  string
	podname    string
	currentmem int
	currentcpu int
}

func main() {
	ns := os.Args[1]
	// var ignorenamespaces []string = []string{"abc", "cdb"}
	clientset := modules.GetClient()
	var pods modules.PodMetricsList
	err := modules.GetMetricClientApi(clientset, &pods, ns)
	if err != nil {
		panic(err.Error())
	}
	for _, m := range pods.Items {
		fmt.Println(m.Metadata.Name, m.Metadata.Namespace, m.Timestamp.String())
	}
	// validnamespaces := modules.GetNamespace(clientset, label, ignorenamespaces)
	// podmetrics := modules.GetMetric(validnamespaces, clientset, clientmetrics, 100)
	// fmt.Println("Num object:", len(podmetrics))
	// for i := range podmetrics {
	// 	fmt.Println(podmetrics[i])
	// }
	// var csvdata [][]string
	// for i := range podmetrics {
	// 	csvdata = append(csvdata, []string{podmetrics[i].Namespace, podmetrics[i].PodName, strconv.Itoa(podmetrics[i].CurrentCpu), strconv.Itoa(podmetrics[i].CurrentMem)})
	// }
	// fmt.Println(csvdata)
}
