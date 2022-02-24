package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "task"
)

type TaskMetric struct {
	ClassName string  `json:"class_name"`
	Host      string  `json:"host"`
	Duration  float64 `json:"duration"`
	Completed float64 `json:"completed"`
	Failed    float64 `json:"failed"`
	Retried   float64 `json:"retried"`
}

type taskCollector struct{}

var (
	labels       = []string{"class_name", "host"}
	taskDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "duration"),
		"Duration task in ms",
		labels,
		nil,
	)
	taskCompleted = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "completed"),
		"Completed task count",
		labels,
		nil,
	)
	taskFailed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "failed"),
		"Failed task count",
		labels,
		nil,
	)
	taskRetried = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "retried"),
		"Retried task count",
		labels,
		nil,
	)
)

func (c taskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- taskDuration
	ch <- taskCompleted
	ch <- taskFailed
	ch <- taskRetried
}

func (c taskCollector) Collect(ch chan<- prometheus.Metric) {
	taskMetrics := getCollectedMetrics()

	for _, metric := range taskMetrics {
		ch <- getMetric(taskDuration, metric)
		ch <- getMetric(taskCompleted, metric)
		ch <- getMetric(taskFailed, metric)
		ch <- getMetric(taskRetried, metric)
	}
}

func getMetric(desc *prometheus.Desc, metric TaskMetric) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		desc,
		prometheus.CounterValue,
		metric.Duration,
		metric.ClassName,
		metric.Host,
	)
}

func main() {
	flag.Parse()
	var listenAddress = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	var collector taskCollector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError, ErrorLog: log.Default()},
	))
	log.Print("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func getCollectedMetrics() []TaskMetric {
	var taskMetricsAll []TaskMetric

	err := filepath.Walk("metrics", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		if filepath.Ext(path) == ".json" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
				return err
			}

			var taskMetrics []TaskMetric
			_ = json.Unmarshal(file, &taskMetrics)
			//log.Print(taskMetrics)
			taskMetricsAll = append(taskMetricsAll, taskMetrics...)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return taskMetricsAll
}
