package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type vsCollector struct {
	metrics map[string]vsMetric
	bigip   *f5.Device
}

type vsMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBVirtualStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewVSCollector(bigip *f5.Device, namespace string) (error, *vsCollector) {
	var (
		subsystem  = "vs"
		labelNames = []string{"partition", "vs"}
	)
	return nil, &vsCollector{
		metrics: map[string]vsMetric{
			"clientside_bitsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_bits_in"),
					"clientside_bits_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_bitsIn.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_curConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_cur_conns"),
					"clientside_cur_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_curConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_bitsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_bits_out"),
					"clientside_bits_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_bitsOut.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_pktsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_pkts_in"),
					"clientside_pkts_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_pktsIn.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_pktsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_pkts_out"),
					"clientside_pkts_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_pktsOut.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_totConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_tot_conns"),
					"clientside_tot_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_totConns.Value)
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
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					if entries.Status_availabilityState.Description == "available" {
						return 1
					}
					return 0
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_rejects": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_rejects"),
					"syncookie_rejects",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Syncookie_rejects.Value)
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
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return float64(entries.Clientside_totConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
		},
		bigip: bigip,
	}
}

func (c *vsCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, virtualServers := c.bigip.ShowVirtuals()
	if err != nil {
		log.Fatal(err)
	}
	for _, virtualServer := range virtualServers.Items {
		err, virtualStats := c.bigip.ShowVirtualStats("/" + virtualServer.Partition + "/" + virtualServer.Name)
		if err != nil {
			log.Fatal(err)
		}
		lables := []string{virtualServer.Partition, virtualServer.Name}
		urlKey := "https://localhost/mgmt/tm/ltm/virtual/~" + virtualServer.Partition + "~" + virtualServer.Name + "/~" + virtualServer.Partition + "~" + virtualServer.Name + "/stats"
		for _, metric := range c.metrics {
			ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(virtualStats.Entries[urlKey].NestedStats.Entries), lables...)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Getting stats took %s", elapsed)
}

func (c *vsCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
}
