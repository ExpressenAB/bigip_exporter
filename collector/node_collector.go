package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"time"
)

type nodeCollector struct {
	metrics                   map[string]nodeMetric
	bigip                     *f5.Device
	partitions_list           []string
	collector_scrape_status   *prometheus.GaugeVec
	collector_scrape_duration *prometheus.SummaryVec
}

type nodeMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBNodeStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewNodeCollector(bigip *f5.Device, namespace string, partitions_list []string) (error, *nodeCollector) {
	var (
		subsystem  = "node"
		labelNames = []string{"partition", "node"}
	)
	return nil, &nodeCollector{
		metrics: map[string]nodeMetric{
			"serverside_bytesOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_out"),
					"serverside_bytes_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsOut.Value / 8)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_maxConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_max_conns"),
					"serverside_max_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_maxConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_curConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_cur_conns"),
					"serverside_cur_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_curConns.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"serverside_pktsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_pkts_out"),
					"serverside_pkts_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsOut.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"totRequests": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "tot_requests"),
					"tot_requests",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.TotRequests.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_pktsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_pkts_in"),
					"serverside_pkts_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsIn.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_totConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_tot_conns"),
					"serverside_tot_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_totConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_bytesIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_in"),
					"serverside_bytes_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsIn.Value / 8)
				},
				valueType: prometheus.CounterValue,
			},
			"curSessions": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "cur_sessions"),
					"cur_sessions",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.CurSessions.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"status_availabilityState": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "status_availability_state"),
					"status_availability_state",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					if entries.Status_availabilityState.Description == "available" {
						return 1
					}
					return 0
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

func (c *nodeCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, allNodeStats := c.bigip.ShowAllNodeStats()
	success := true
	if err != nil {
		success = false
		log.Println(err)
	} else {
		for key, nodeStats := range allNodeStats.Entries {
			keyParts := strings.Split(key, "/")
			path := keyParts[len(keyParts)-2]
			pathParts := strings.Split(path, "~")
			partition := pathParts[1]
			nodeName := pathParts[len(pathParts)-1]

			if c.partitions_list != nil && !stringInSlice(partition, c.partitions_list) {
				continue
			}

			lables := []string{partition, nodeName}
			for _, metric := range c.metrics {
				ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(nodeStats.NestedStats.Entries), lables...)
			}
		}
	}
	elapsed := time.Since(start)
	if success {
		c.collector_scrape_status.WithLabelValues("node").Set(float64(1))
	} else {
		c.collector_scrape_status.WithLabelValues("node").Set(float64(0))
	}
	c.collector_scrape_duration.WithLabelValues("node").Observe(float64(elapsed.Seconds()))
	c.collector_scrape_status.Collect(ch)
	c.collector_scrape_duration.Collect(ch)
	log.Printf("Node was succes: %t", success)
	log.Printf("Getting node stats took %s", elapsed)
}

func (c *nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
	c.collector_scrape_status.Describe(ch)
	c.collector_scrape_duration.Describe(ch)
}
