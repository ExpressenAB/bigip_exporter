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
	bigipHost           string
	collectors          map[string]prometheus.Collector
	totalScrapeDuration *prometheus.SummaryVec
}

var (
	logger = loggo.GetLogger("")
)

// NewBigipCollector returns a collector that wraps all the collectors
func NewBigipCollector(bigip *f5.Device, namespace string, partitionsList []string, bigipHost string) (*BigipCollector, error) {
	vsCollector, _ := NewVSCollector(bigip, namespace, partitionsList, bigipHost)
	poolCollector, _ := NewPoolCollector(bigip, namespace, partitionsList, bigipHost)
	nodeCollector, _ := NewNodeCollector(bigip, namespace, partitionsList, bigipHost)
	ruleCollector, _ := NewRuleCollector(bigip, namespace, partitionsList, bigipHost)
	return &BigipCollector{
		collectors: map[string]prometheus.Collector{
			"node": nodeCollector,
			"pool": poolCollector,
			"rule": ruleCollector,
			"vs":   vsCollector,
		},
		totalScrapeDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "total_scrape_duration",
				Help:      "total_scrape_duration",
			},
			[]string{"host"},
		),
		bigipHost: bigipHost,
	}, nil
}

// Collect collects all metrics exported by this exporter by delegating
// to the different collectors
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
	c.totalScrapeDuration.WithLabelValues(c.bigipHost).Observe(float64(elapsed.Seconds()))
	c.totalScrapeDuration.Collect(ch)
	logger.Debugf("Total collection time was: %s", elapsed)
}

// Describe describes all metrics exported by this exporter by delegating
// to the different collectors
func (c *BigipCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
	c.totalScrapeDuration.Describe(ch)
}
