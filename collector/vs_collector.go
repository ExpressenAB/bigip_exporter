package collector

import (
	"strings"
	"time"

	"github.com/pr8kerl/f5er/f5"
	"github.com/prometheus/client_golang/prometheus"
)

// A VSCollector implements the prometheus.Collector.
type VSCollector struct {
	metrics                 map[string]vsMetric
	bigip                   *f5.Device
	partitionsList          []string
	collectorScrapeStatus   *prometheus.GaugeVec
	collectorScrapeDuration *prometheus.SummaryVec
}

type vsMetric struct {
	desc      *prometheus.Desc
	extract   func(f5.LBVirtualStatsInnerEntries) float64
	valueType prometheus.ValueType
}

// NewVSCollector returns a collector that collecting virtual server statistics.
func NewVSCollector(bigip *f5.Device, namespace string, partitionsList []string) (*VSCollector, error) {
	var (
		subsystem  = "vs"
		labelNames = []string{"partition", "vs"}
	)
	return &VSCollector{
		metrics: map[string]vsMetric{
			"syncookie_accepts": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_accepts"),
					"syncookie_accepts",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_accepts.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_bytesOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_bytes_out"),
					"ephemeral_bytes_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_bitsOut.Value / 8
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_bytesOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_bytes_out"),
					"clientside_bytes_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Clientside_bitsOut.Value / 8
				},
				valueType: prometheus.CounterValue,
			},
			"fiveMinAvgUsageRatio": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "five_min_avg_usage_ratio"),
					"five_min_avg_usage_ratio",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.FiveMinAvgUsageRatio.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"fiveSecAvgUsageRatio": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "five_sec_avg_usage_ratio"),
					"five_sec_avg_usage_ratio",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.FiveSecAvgUsageRatio.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"syncookie_syncookies": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_syncookies"),
					"syncookie_syncookies",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_syncookies.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_slowKilled": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_slow_killed"),
					"ephemeral_slow_killed",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_slowKilled.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_pktsOut": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_pkts_out"),
					"ephemeral_pkts_out",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_pktsOut.Value
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
					return entries.Syncookie_rejects.Value
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_syncacheCurr": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_syncache_curr"),
					"syncookie_syncache_curr",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_syncacheCurr.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"csMinConnDur": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "cs_min_conn_dur"),
					"cs_min_conn_dur",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.CsMinConnDur.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"csMeanConnDur": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "cs_mean_conn_dur"),
					"cs_mean_conn_dur",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.CsMeanConnDur.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"syncookie_swsyncookieInstance": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_swsyncookie_instance"),
					"syncookie_swsyncookie_instance",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_swsyncookieInstance.Value
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_syncacheOver": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_syncache_over"),
					"syncookie_syncache_over",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_syncacheOver.Value
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_hwAccepts": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_hw_accepts"),
					"syncookie_hw_accepts",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_hwAccepts.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_pktsIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_pkts_in"),
					"ephemeral_pkts_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_pktsIn.Value
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
					return entries.Clientside_totConns.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_curConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_cur_conns"),
					"ephemeral_cur_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_curConns.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"clientside_evictedConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_evicted_conns"),
					"clientside_evicted_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Clientside_evictedConns.Value
				},
				valueType: prometheus.CounterValue,
			},
			"oneMinAvgUsageRatio": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "one_min_avg_usage_ratio"),
					"one_min_avg_usage_ratio",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.OneMinAvgUsageRatio.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"ephemeral_evictedConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_evicted_conns"),
					"ephemeral_evicted_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_evictedConns.Value
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_slowKilled": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_slow_killed"),
					"clientside_slow_killed",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Clientside_slowKilled.Value
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_bytesIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_bytes_in"),
					"clientside_bytes_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Clientside_bitsIn.Value / 8
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_maxConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_max_conns"),
					"ephemeral_max_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_maxConns.Value
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_hwsyncookieInstance": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_hwsyncookie_instance"),
					"syncookie_hwsyncookie_instance",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_hwsyncookieInstance.Value
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
					return entries.Clientside_pktsOut.Value
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
					return entries.Clientside_curConns.Value
				},
				valueType: prometheus.GaugeValue,
			},
			"ephemeral_bytesIn": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_bytes_in"),
					"ephemeral_bytes_in",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_bitsIn.Value / 8
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
					return entries.Clientside_pktsIn.Value
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
					return entries.TotRequests.Value
				},
				valueType: prometheus.CounterValue,
			},
			"csMaxConnDur": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "cs_max_conn_dur"),
					"cs_max_conn_dur",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.CsMaxConnDur.Value
				},
				valueType: prometheus.CounterValue,
			},
			"syncookie_hwSyncookies": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "syncookie_hw_syncookies"),
					"syncookie_hw_syncookies",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Syncookie_hwSyncookies.Value
				},
				valueType: prometheus.CounterValue,
			},
			"clientside_maxConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "clientside_max_conns"),
					"clientside_max_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Clientside_maxConns.Value
				},
				valueType: prometheus.CounterValue,
			},
			"ephemeral_totConns": {
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "ephemeral_tot_conns"),
					"ephemeral_tot_conns",
					labelNames,
					nil,
				),
				extract: func(entries f5.LBVirtualStatsInnerEntries) float64 {
					return entries.Ephemeral_totConns.Value
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
				valueType: prometheus.GaugeValue,
			},
		},
		collectorScrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
		collectorScrapeDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "collector_scrape_duration",
				Help:      "collector_scrape_duration",
			},
			[]string{"collector"},
		),
		bigip:          bigip,
		partitionsList: partitionsList,
	}, nil
}

// Collect collects metrics for BIG-IP virtual servers.
func (c *VSCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	err, allVirtualServerStats := c.bigip.ShowAllVirtualStats()
	if err != nil {
		c.collectorScrapeStatus.WithLabelValues("vs").Set(0)
		logger.Warningf("Failed to get statistics for virtual servers")
	} else {
		for key, virtualStats := range allVirtualServerStats.Entries {
			keyParts := strings.Split(key, "/")
			path := keyParts[len(keyParts)-2]
			pathParts := strings.Split(path, "~")
			partition := pathParts[1]
			vsName := pathParts[len(pathParts)-1]

			if c.partitionsList != nil && !stringInSlice(partition, c.partitionsList) {
				continue
			}

			labels := []string{partition, vsName}
			for _, metric := range c.metrics {
				ch <- prometheus.MustNewConstMetric(
					metric.desc,
					metric.valueType,
					metric.extract(virtualStats.NestedStats.Entries),
					labels...,
				)
			}
		}
		c.collectorScrapeStatus.WithLabelValues("vs").Set(1)
		logger.Debugf("Successfully fetched statistics for virtual servers")
	}

	elapsed := time.Since(start)
	c.collectorScrapeDuration.WithLabelValues("vs").Observe(elapsed.Seconds())
	c.collectorScrapeStatus.Collect(ch)
	c.collectorScrapeDuration.Collect(ch)
	logger.Debugf("Getting virtual server statistics took %s", elapsed)
}

// Describe describes the metrics exported from this collector.
func (c *VSCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
	c.collectorScrapeStatus.Describe(ch)
	c.collectorScrapeDuration.Describe(ch)
}
