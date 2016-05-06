package collector

import (
	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type poolCollector struct {
	metrics map[string]poolMetric
	bigip   *f5.Device
}

type poolMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBPoolStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewPoolCollector(bigip *f5.Device, namespace string) (error, *poolCollector) {
	var (
		subsystem  = "pool"
		labelNames = []string{"partition", "pool"}
	)
	return nil, &poolCollector{
		metrics: map[string]poolMetric{
			"activeMemberCnt": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "active_member_cnt"),
					"active_member_cnt",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ActiveMemberCnt.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.CurSessions.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsIn.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"serverside_bitsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bits_out"),
					"serverside_bits_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsOut.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_curConns.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsIn.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_pktsOut.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_totConns.Value)
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					if entries.Status_availabilityState.Description == "available" {
						return 1
					}
					return 0
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
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.TotRequests.Value)
				},
				valueType: prometheus.CounterValue,
			},
		},
		bigip: bigip,
	}
}

func (c *poolCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, pools := c.bigip.ShowPools()
	if err != nil {
		log.Fatal(err)
	}
	for _, pool := range pools.Items {
		err, poolStats := c.bigip.ShowPoolStats("/" + pool.Partition + "/" + pool.Name)
		if err != nil {
			log.Fatal(err)
		}
		lables := []string{pool.Partition, pool.Name}
		urlKey := "https://localhost/mgmt/tm/ltm/pool/~" + pool.Partition + "~" + pool.Name + "/~" + pool.Partition + "~" + pool.Name + "/stats"
		for _, metric := range c.metrics {
			ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, metric.extract(poolStats.Entries[urlKey].NestedStats.Entries), lables...)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Getting stats took %s", elapsed)
}

func (c *poolCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
}
