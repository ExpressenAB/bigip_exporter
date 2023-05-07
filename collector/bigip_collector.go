package collector

import (
	"sync"
	"time"

	"github.com/juju/loggo"
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

// A BigipCollector implements the prometheus.Collector.
type BigipCollector struct {
	collectors          map[string]prometheus.Collector
	totalScrapeDuration prometheus.Summary
}

var logger = loggo.GetLogger("")

// NewBigipCollector returns a collector that wraps all the collectors.
func NewBigipCollector(bigip *f5.Device, namespace string, partitionsList []string) (*BigipCollector, error) {
	vsCollector, _ := NewVSCollector(bigip, namespace, partitionsList)
	poolCollector, _ := NewPoolCollector(bigip, namespace, partitionsList)
	nodeCollector, _ := NewNodeCollector(bigip, namespace, partitionsList)
	ruleCollector, _ := NewRuleCollector(bigip, namespace, partitionsList)

	return &BigipCollector{
		collectors: map[string]prometheus.Collector{
			"node": nodeCollector,
			"pool": poolCollector,
			"rule": ruleCollector,
			"vs":   vsCollector,
		},
		totalScrapeDuration: prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "total_scrape_duration",
				Help:      "total_scrape_duration",
			},
		),
	}, nil
}

// Collect collects all metrics exported by this exporter by delegating
// to the different collectors.
func (c *BigipCollector) Collect(ch chan<- prometheus.Metric) {
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
	c.totalScrapeDuration.Observe(elapsed.Seconds())
	ch <- c.totalScrapeDuration
	logger.Debugf("Total collection time was: %s", elapsed)
}

// Describe describes all metrics exported by this exporter by delegating
// to the different collectors.
func (c *BigipCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
	ch <- c.totalScrapeDuration.Desc()
}
