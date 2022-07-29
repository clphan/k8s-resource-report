package main

import (
	"fmt"
	"os"
	"strconv"
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
	clientset, clientmetrics := modules.getclient()
	validnamespaces := modules.getnamespace(clientset, label, ignorenamespaces)
	podmetrics := getmetric(validnamespaces, clientset, clientmetrics)
	fmt.Println("Num object:", len(podmetrics))
	var csvdata [][]string
	for i := range podmetrics {
		csvdata[i][0] = podmetrics[i].namespace
		csvdata[i][1] = podmetrics[i].podname
		csvdata[i][2] = strconv.Itoa(podmetrics[i].currentcpu)
		csvdata[i][3] = strconv.Itoa(podmetrics[i].currentmem)
	}
}
