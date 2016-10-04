package f5

import (
	"fmt"
	"github.com/fatih/structs"
	"os"
	"strings"
	"time"
)

type LBStatsValue struct {
	// use a wide float - some values can be BIG
	Value float64 `json:"value"`
}

type LBStatsDescription struct {
	Description string `json:"description"`
}

type LBObjectStatsMap map[string]LBStatsValue

type LBObjectStats struct {
	Entries LBObjectStatsMap `json:"entries"`
}

// Metric is a struct that defines the relevant properties of a graphite metric
type GraphiteDataPoint struct {
	Key       string
	Value     float64
	Timestamp int64
}

func NewGraphiteDataPoint(key string, value float64, timestamp int64) GraphiteDataPoint {
	return GraphiteDataPoint{
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
	}
}

func (dp GraphiteDataPoint) String() string {
	return fmt.Sprintf(
		"%s %.f %d",
		dp.Key,
		dp.Value,
		dp.Timestamp,
	)
}

func (f *Device) Stats() (error, []GraphiteDataPoint) {

	data := make([]GraphiteDataPoint, 0, 16384)
	err, pools := f.StatsPools()
	if err != nil {
		return err, nil
	}
	data = append(data, pools...)

	err, virtuals := f.StatsVirtuals()
	if err != nil {
		return err, nil
	}
	data = append(data, virtuals...)

	err, nodes := f.StatsNodes()
	if err != nil {
		return err, nil
	}
	data = append(data, nodes...)

	err, rules := f.StatsRules()
	if err != nil {
		return err, nil
	}
	data = append(data, rules...)
	return nil, data
}

func (f *Device) StatsPool(pname string) (error, []GraphiteDataPoint) {

	pool := strings.Replace(pname, "/", "~", -1)
	splitter := func(c rune) bool {
		// split if "/" or "~"
		return c == '\u002f' || c == '\u007e'

	}
	fields := strings.FieldsFunc(pool, splitter)
	partition := fields[0]
	poolname := fields[1]
	if len(partition) < 1 || len(poolname) < 1 {
		return fmt.Errorf("error: cannot parse partition and poolname for given pool: %s", pname), nil
	}
	prefix := f.StatsPathPrefix + partition + ".pool."

	data := make([]GraphiteDataPoint, 0, 1024)
	start := time.Now()
	timestamp := start.Unix()

	err, res := f.ShowPoolStats(pool)
	if err != nil {
		return err, nil
	} else {

		for key := range res.Entries {

			skey := prefix + poolname + "." + key
			var v interface{}
			v = res.Entries[key]

			if value, ok := v.(LBStatsValue); ok {
				if value.Value > 0 {
					// only print value if it is greater than zero
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				} else if f.StatsShowZeroes {
					// except if I really want it
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				}
			}

		}

		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data

	}

}

func (f *Device) StatsPools() (error, []GraphiteDataPoint) {

	data := make([]GraphiteDataPoint, 0, 4096)
	start := time.Now()
	timestamp := start.Unix()
	err, res := f.ShowAllPoolStats()
	if err != nil {
		return err, nil
	} else {

		splitter := func(c rune) bool {
			// split if "/" or "~"
			return c == '\u002f' || c == '\u007e'

		}

		for surl, stats := range res.Entries {

			fields := strings.FieldsFunc(surl, splitter)
			partition := fields[6]
			poolname := fields[7]
			if len(partition) < 1 || len(poolname) < 1 {
				fmt.Fprintf(os.Stderr, "warn: cannot parse partition and poolname for pool given url: %s", surl)
				continue
			}

			entries := structs.New(stats.NestedStats.Entries)
			fnames := entries.Names()

			prefix := f.StatsPathPrefix + partition + ".pool."

			for key := range fnames {

				skey := prefix + poolname + "." + fnames[key]
				v := entries.Field(fnames[key])

				if value, ok := v.Value().(LBStatsValue); ok {
					if value.Value > 0 {
						// only print value if it is greater than zero
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					} else if f.StatsShowZeroes {
						// except if I really want it
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					}
				}

			}

		}
		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data
	}

}

func (f *Device) StatsNode(nname string) (error, []GraphiteDataPoint) {

	node := strings.Replace(nname, "/", "~", -1)
	splitter := func(c rune) bool {
		// split if "/" or "~"
		return c == '\u002f' || c == '\u007e'

	}
	fields := strings.FieldsFunc(node, splitter)
	partition := fields[0]
	nodename := fields[1]
	if len(partition) < 1 || len(nodename) < 1 {
		return fmt.Errorf("error: cannot parse partition and nodename for given node: %s", nname), nil
	}
	prefix := f.StatsPathPrefix + partition + ".node."

	data := make([]GraphiteDataPoint, 0, 1024)
	start := time.Now()
	timestamp := start.Unix()

	err, res := f.ShowNodeStats(node)
	if err != nil {
		return err, nil
	} else {

		for key := range res.Entries {

			skey := prefix + nodename + "." + key
			var v interface{}
			v = res.Entries[key]

			// only print if it is an object of LBStatsValue
			if value, ok := v.(LBStatsValue); ok {
				if value.Value > 0 {
					// only print value if it is greater than zero
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				} else if f.StatsShowZeroes {
					// except if I really want it
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				}
			}

		}

		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data

	}

}

func (f *Device) StatsNodes() (error, []GraphiteDataPoint) {

	data := make([]GraphiteDataPoint, 0, 4096)
	start := time.Now()
	timestamp := start.Unix()
	err, res := f.ShowAllNodeStats()
	if err != nil {
		return err, nil
	} else {

		splitter := func(c rune) bool {
			// split if "/" or "~"
			return c == '\u002f' || c == '\u007e'

		}

		for surl, stats := range res.Entries {

			fields := strings.FieldsFunc(surl, splitter)
			partition := fields[6]
			nodename := fields[7]
			if len(partition) < 1 || len(nodename) < 1 {
				fmt.Fprintf(os.Stderr, "warn: cannot parse partition and nodename for node given url: %s", surl)
				continue
			}

			entries := structs.New(stats.NestedStats.Entries)
			fnames := entries.Names()

			prefix := f.StatsPathPrefix + partition + ".node."

			for key := range fnames {

				skey := prefix + nodename + "." + fnames[key]
				v := entries.Field(fnames[key])

				if value, ok := v.Value().(LBStatsValue); ok {
					if value.Value > 0 {
						// only print value if it is greater than zero
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					} else if f.StatsShowZeroes {
						// except if I really want it
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					}
				}

			}

		}
		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data
	}

}

func (f *Device) StatsVirtual(vname string) (error, []GraphiteDataPoint) {

	virtual := strings.Replace(vname, "/", "~", -1)
	splitter := func(c rune) bool {
		// split if "/" or "~"
		return c == '\u002f' || c == '\u007e'

	}
	fields := strings.FieldsFunc(virtual, splitter)
	partition := fields[0]
	virtualname := fields[1]
	if len(partition) < 1 || len(virtualname) < 1 {
		return fmt.Errorf("error: cannot parse partition and virtualname for given pool: %s", vname), nil
	}
	prefix := f.StatsPathPrefix + partition + ".virtual."

	data := make([]GraphiteDataPoint, 0, 1024)
	start := time.Now()
	timestamp := start.Unix()

	err, res := f.ShowVirtualStats(virtual)
	if err != nil {
		return err, nil
	} else {

		for key := range res.Entries {

			skey := prefix + virtualname + "." + key
			var v interface{}
			v = res.Entries[key]

			if value, ok := v.(LBStatsValue); ok {
				if value.Value > 0 {
					// only print value if it is greater than zero
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				} else if f.StatsShowZeroes {
					// except if I really want it
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				}
			}

		}

		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data

	}

}

func (f *Device) StatsVirtuals() (error, []GraphiteDataPoint) {

	data := make([]GraphiteDataPoint, 0, 4096)
	start := time.Now()
	timestamp := start.Unix()
	err, res := f.ShowAllVirtualStats()
	if err != nil {
		return err, nil
	} else {

		splitter := func(c rune) bool {
			// split if "/" or "~"
			return c == '\u002f' || c == '\u007e'

		}

		for surl, stats := range res.Entries {

			fields := strings.FieldsFunc(surl, splitter)
			partition := fields[6]
			virtualname := fields[7]
			if len(partition) < 1 || len(virtualname) < 1 {
				fmt.Fprintf(os.Stderr, "warn: cannot parse partition and virtualname for virtual given url: %s", surl)
				continue
			}

			entries := structs.New(stats.NestedStats.Entries)
			fnames := entries.Names()

			prefix := f.StatsPathPrefix + partition + ".virtual."

			for key := range fnames {

				skey := prefix + virtualname + "." + fnames[key]
				v := entries.Field(fnames[key])

				if value, ok := v.Value().(LBStatsValue); ok {
					if value.Value > 0 {
						// only print value if it is greater than zero
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					} else if f.StatsShowZeroes {
						// except if I really want it
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					}
				}

			}

		}
		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data
	}

}

func (f *Device) StatsRule(rname string) (error, []GraphiteDataPoint) {

	rule := strings.Replace(rname, "/", "~", -1)
	splitter := func(c rune) bool {
		// split if "/" or "~"
		return c == '\u002f' || c == '\u007e'

	}
	fields := strings.FieldsFunc(rule, splitter)
	partition := fields[0]
	rulename := fields[1]
	if len(partition) < 1 || len(rulename) < 1 {
		return fmt.Errorf("error: cannot parse partition and rulename for given rule: %s", rname), nil
	}
	prefix := f.StatsPathPrefix + partition + ".rule."

	data := make([]GraphiteDataPoint, 0, 1024)
	start := time.Now()
	timestamp := start.Unix()

	err, res := f.ShowRuleStats(rule)
	if err != nil {
		return err, nil
	} else {

		for key := range res.Entries {

			skey := prefix + rulename + "." + key
			var v interface{}
			v = res.Entries[key]

			if value, ok := v.(LBStatsValue); ok {
				if value.Value > 0 {
					// only print value if it is greater than zero
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				} else if f.StatsShowZeroes {
					// except if I really want it
					dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
					data = append(data, dp)
				}
			}

		}

		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data

	}

}

func (f *Device) StatsRules() (error, []GraphiteDataPoint) {

	data := make([]GraphiteDataPoint, 0, 4096)
	start := time.Now()
	timestamp := start.Unix()
	err, res := f.ShowAllRuleStats()
	if err != nil {
		return err, nil
	} else {

		splitter := func(c rune) bool {
			// split if "/" or "~"
			return c == '\u002f' || c == '\u007e'

		}

		for surl, stats := range res.Entries {

			fields := strings.FieldsFunc(surl, splitter)
			partition := fields[6]
			rulename := fields[7]
			if len(partition) < 1 || len(rulename) < 1 {
				fmt.Fprintf(os.Stderr, "warn: cannot parse partition and rulename for rule given url: %s", surl)
				continue
			}

			entries := structs.New(stats.NestedStats.Entries)
			fnames := entries.Names()

			prefix := f.StatsPathPrefix + partition + ".rule."

			for key := range fnames {

				skey := prefix + rulename + "." + fnames[key]
				v := entries.Field(fnames[key])

				if value, ok := v.Value().(LBStatsValue); ok {
					if value.Value > 0 {
						// only print value if it is greater than zero
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					} else if f.StatsShowZeroes {
						// except if I really want it
						dp := NewGraphiteDataPoint(skey, value.Value, timestamp)
						data = append(data, dp)
					}
				}

			}

		}
		//		elapsed := time.Since(start)
		//		fmt.Printf("elapsed: %v\n", elapsed.Seconds())
		return nil, data
	}

}
