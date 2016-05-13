package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type ruleCollector struct {
	metrics map[string]ruleMetric
	bigip   *f5.Device
	partitions_list []string
}

type ruleMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBRuleStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewRuleCollector(bigip *f5.Device, namespace string, partitions_list []string) (error, *ruleCollector) {
	var (
		subsystem  = "rule"
		labelNames = []string{"partition", "rule"}
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
		bigip: bigip,
		partitions_list: partitions_list,
	}
}

func (c *ruleCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, rules := c.bigip.ShowRules()
	if err != nil {
		log.Fatal(err)
	}
	for _, rule := range rules.Items {
		if c.partitions_list != nil && !stringInSlice(rule.Partition, c.partitions_list) {
			continue
		}
		err, ruleStats := c.bigip.ShowRuleStats("/" + rule.Partition + "/" + rule.Name)
		if err != nil {
			log.Fatal(err)
		}
		lables := []string{rule.Partition, rule.Name}
		urlKey := "https://localhost/mgmt/tm/ltm/rule/~" + rule.Partition + "~" + rule.Name + "/~" + rule.Partition + "~" + rule.Name + ":HTTP_REQUEST/stats"
		for _, metric := range c.metrics {
			ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(ruleStats.Entries[urlKey].NestedStats.Entries), lables...)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Getting stats took %s", elapsed)
}

func (c *ruleCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
}
