package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"time"
)

type bigipCollector struct {
	collectors            map[string]prometheus.Collector
	total_scrape_duration prometheus.Summary
}

func NewBigIpCollector(bigip *f5.Device, namespace string, partitions_list []string) (error, *bigipCollector) {
	_, vsCollector := NewVSCollector(bigip, namespace, partitions_list)
	_, poolCollector := NewPoolCollector(bigip, namespace, partitions_list)
	_, nodeCollector := NewNodeCollector(bigip, namespace, partitions_list)
	_, ruleCollector := NewRuleCollector(bigip, namespace, partitions_list)
	return nil, &bigipCollector{
		collectors: map[string]prometheus.Collector{
			"node": nodeCollector,
			"pool": poolCollector,
			"rule": ruleCollector,
			"vs":   vsCollector,
		},
		total_scrape_duration: prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "total_scrape_duration",
				Help:      "total_scrape_duration",
			},
		),
	}
}

func (c *bigipCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))
	start := time.Now()
	for _, collector := range c.collectors {
		go func(coll prometheus.Collector) {
			coll.Collect(ch)
			wg.Done()
		}(collector)
	}
	wg.Wait()
	elapsed := time.Since(start)
	c.total_scrape_duration.Observe(float64(elapsed.Seconds()))
	ch <- c.total_scrape_duration
	log.Println(elapsed.Seconds())
}

func (c *bigipCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
	ch <- c.total_scrape_duration.Desc()
}
