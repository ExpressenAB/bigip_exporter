package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type nodeCollector struct {
	metrics map[string]nodeMetric
	bigip   *f5.Device
	partitions_list []string
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
			"serverside_bitsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bits_out"),
					"serverside_bits_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsOut.Value)
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
				valueType: prometheus.CounterValue,
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
			"serverside_bitsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bits_in"),
					"serverside_bits_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBNodeStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsIn.Value)
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
				valueType: prometheus.CounterValue,
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
				valueType: prometheus.CounterValue,
			},
		},
		bigip: bigip,
		partitions_list: partitions_list,
	}
}

func (c *nodeCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, nodes := c.bigip.ShowNodes()
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes.Items {
		if c.partitions_list != nil && !stringInSlice(node.Partition, c.partitions_list) {
			continue
		}
		err, nodeStats := c.bigip.ShowNodeStats("/" + node.Partition + "/" + node.Name)
		if err != nil {
			log.Fatal(err)
		}
		lables := []string{node.Partition, node.Name}
		urlKey := "https://localhost/mgmt/tm/ltm/node/~" + node.Partition + "~" + node.Name + "/~" + node.Partition + "~" + node.Name + "/stats"
		for _, metric := range c.metrics {
			ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(nodeStats.Entries[urlKey].NestedStats.Entries), lables...)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Getting stats took %s", elapsed)
}

func (c *nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
}
