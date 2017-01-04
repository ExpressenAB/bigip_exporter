package collector

import (
	"strings"
	"time"

	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

type ruleCollector struct {
	metrics                   map[string]ruleMetric
	bigip                     *f5.Device
	partitions_list           []string
	collector_scrape_status   *prometheus.GaugeVec
	collector_scrape_duration *prometheus.SummaryVec
}

type ruleMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBRuleStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewRuleCollector(bigip *f5.Device, namespace string, partitions_list []string) (error, *ruleCollector) {
	var (
		subsystem  = "rule"
		labelNames = []string{"partition", "rule", "event"}
	)
	return nil, &ruleCollector{
		metrics: map[string]ruleMetric{
			"priority": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "priority"),
					"priority",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.Priority.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"failures": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "failures"),
					"failures",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.Failures.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"totalExecutions": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "total_executions"),
					"total_executions",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.TotalExecutions.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"aborts": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "aborts"),
					"aborts",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.Aborts.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"minCycles": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "min_cycles"),
					"min_cycles",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.MinCycles.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"maxCycles": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "max_cycles"),
					"max_cycles",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.MaxCycles.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"avgCycles": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "avg_cycles"),
					"avg_cycles",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBRuleStatsInnerEntries) float64 {
					return float64(entries.AvgCycles.Value)
				},
				valueType: prometheus.GaugeValue,
			},
		},
		collector_scrape_status: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
		collector_scrape_duration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "collector_scrape_duration",
				Help:      "collector_scrape_duration",
			},
			[]string{"collector"},
		),
		bigip:           bigip,
		partitions_list: partitions_list,
	}
}

func (c *ruleCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, allRuleStats := c.bigip.ShowAllRuleStats()
	if err != nil {
		c.collector_scrape_status.WithLabelValues("rule").Set(float64(0))
		logger.Warningf("Failed to get statistics for rules")
	} else {
		for key, ruleStats := range allRuleStats.Entries {
			keyParts := strings.Split(key, "/")
			path := keyParts[len(keyParts)-2]
			pathParts := strings.Split(path, "~")
			partition := pathParts[1]
			eventParts := strings.Split(pathParts[len(pathParts)-1], ":")
			ruleName := eventParts[0]
			event := eventParts[1]

			if c.partitions_list != nil && !stringInSlice(partition, c.partitions_list) {
				continue
			}

			labels := []string{partition, ruleName, event}
			for _, metric := range c.metrics {
				ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(ruleStats.NestedStats.Entries), labels...)
			}
		}
		c.collector_scrape_status.WithLabelValues("rule").Set(float64(1))
		logger.Debugf("Successfully fetched statistics for rules")
	}

	elapsed := time.Since(start)
	c.collector_scrape_duration.WithLabelValues("rule").Observe(float64(elapsed.Seconds()))
	c.collector_scrape_status.Collect(ch)
	c.collector_scrape_duration.Collect(ch)
	logger.Debugf("Getting rule stats took %s", elapsed)
}

func (c *ruleCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
	c.collector_scrape_status.Describe(ch)
	c.collector_scrape_duration.Describe(ch)
}
