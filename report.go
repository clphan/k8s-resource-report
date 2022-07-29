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

var ignorenamespaces []string = []string{
	"aurora-automations-apps",
	"card-batchjob-service",
	"card-sftp-service",
	"card-mock-service",
}

func main() {
	label := os.Args[1]
	clientset, clientmetrics := getclient()
	validnamespaces := getnamespace(clientset, label)
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
