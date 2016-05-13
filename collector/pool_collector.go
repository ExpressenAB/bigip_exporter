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
	partitions_list []string
}

type poolMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBPoolStatsInnerEntries) float64
	valueType prometheus.ValueType
}

func NewPoolCollector(bigip *f5.Device, namespace string, partitions_list []string) (error, *poolCollector) {
	var (
		subsystem  = "pool"
		labelNames = []string{"partition", "pool"}
	)
	return nil, &poolCollector{
		metrics: map[string]poolMetric{
			"connqAll_ageMax": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_age_max"),
					"connq_all_age_max",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_ageMax.Value/1000)
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
				valueType: prometheus.GaugeValue,
			},
			"minActiveMembers": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "min_active_members"),
					"min_active_members",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.MinActiveMembers.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"serverside_bytesIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_in"),
					"serverside_bytes_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsIn.Value/8)
				},
				valueType: prometheus.CounterValue,
			},
			"connqAll_serviced": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_serviced"),
					"connq_all_serviced",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_serviced.Value)
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
			"serverside_maxConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_max_conns"),
					"serverside_max_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_maxConns.Value)
				},
				valueType: prometheus.CounterValue,
			},
			"connq_depth": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_depth"),
					"connq_depth",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_depth.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"connqAll_depth": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_depth"),
					"connq_all_depth",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_depth.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"connqAll_ageHead": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_age_head"),
					"connq_all_age_head",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_ageHead.Value/1000)
				},
				valueType: prometheus.GaugeValue,
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
				valueType: prometheus.GaugeValue,
			},
			"connq_serviced": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_serviced"),
					"connq_serviced",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_serviced.Value)
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
			"connqAll_ageEdm": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_age_edm"),
					"connq_all_age_edm",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_ageEdm.Value/1000)
				},
				valueType: prometheus.GaugeValue,
			},
			"connq_ageHead": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_age_head"),
					"connq_age_head",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_ageHead.Value/1000)
				},
				valueType: prometheus.GaugeValue,
			},
			"connq_ageMax": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_age_max"),
					"connq_age_max",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_ageMax.Value/1000)
				},
				valueType: prometheus.CounterValue,
			},
			"connq_ageEdm": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_age_edm"),
					"connq_age_edm",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_ageEdm.Value)
				},
				valueType: prometheus.GaugeValue,
			},
			"serverside_bytesOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "serverside_bytes_out"),
					"serverside_bytes_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Serverside_bitsOut.Value/8)
				},
				valueType: prometheus.CounterValue,
			},
			"connq_ageEma": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_age_ema"),
					"connq_age_ema",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.Connq_ageEma.Value/1000)
				},
				valueType: prometheus.GaugeValue,
			},
			"connqAll_ageEma": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "connq_all_age_ema"),
					"connq_all_age_ema",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBPoolStatsInnerEntries) float64 {
					return float64(entries.ConnqAll_ageEma.Value/1000)
				},
				valueType: prometheus.GaugeValue,
			},
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
				valueType: prometheus.GaugeValue,
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
				valueType: prometheus.GaugeValue,
			},
		},
		bigip: bigip,
		partitions_list: partitions_list,
	}
}

func (c *poolCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, pools := c.bigip.ShowPools()
	if err != nil {
		log.Fatal(err)
	}
	for _, pool := range pools.Items {
		if c.partitions_list != nil && !stringInSlice(pool.Partition, c.partitions_list) {
			continue
		}
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
