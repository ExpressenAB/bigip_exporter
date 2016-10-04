package f5

import (
	"encoding/json"
	"strings"
)

type LBNodeFQDN struct {
	AddressFamily string `json:"addressFamily"`
	AutoPopulate  string `json:"autopopulate"`
	DownInterval  int    `json:"downInterval"`
	// hack - ref issue https://github.com/pr8kerl/f5er/issues/9
	// BIG-IP v12.0 returns a string, whereas v11 returns an int
	// if you use this field, you'll have to convert it explicitly before use :(
	Interval interface{} `json:"interval"`
}

type LBNode struct {
	Name            string     `json:"name"`
	Partition       string     `json:"partition"`
	FullPath        string     `json:"fullPath"`
	Generation      int        `json:"generation"`
	Address         string     `json:"address,omitEmpty"`
	ConnectionLimit int        `json:"connectionLimit"`
	Fqdn            LBNodeFQDN `json:"fqdn"`
	Logging         string     `json:"logging"`
	Monitor         string     `json:"monitor"`
	RateLimit       string     `json:"rateLimit"`
	Session         string     `json:"session,omitEmpty"`
	State           string     `json:"state,omitEmpty"`
}

type LBNodeRef struct {
	Link  string   `json:"selfLink"`
	Items []LBNode `json":items"`
}

type LBNodes struct {
	Items []LBNode `json:"items"`
}

type LBNodeFQDNUpdate struct {
	DownInterval int `json:"downInterval"`
	Interval     int `json:"interval"`
}

type LBNodeUpdate struct {
	Name            string           `json:"name"`
	Partition       string           `json:"partition"`
	FullPath        string           `json:"fullPath"`
	Generation      int              `json:"generation"`
	ConnectionLimit int              `json:"connectionLimit"`
	Fqdn            LBNodeFQDNUpdate `json:"fqdn"`
	Logging         string           `json:"logging"`
	Monitor         string           `json:"monitor"`
	RateLimit       string           `json:"rateLimit"`
}

type LBNodeStatsDescription struct {
	Description string `json:"description"`
}

type LBNodeStatsValue struct {
	Value int `json:"value"`
}

type LBNodeStatsInnerEntries struct {
	Addr                     LBNodeStatsDescription `json:"addr"`
	CurSessions              LBStatsValue           `json:"curSessions"`
	MonitorRule              LBNodeStatsDescription `json:"monitorRule"`
	MonitorStatus            LBNodeStatsDescription `json:"monitorStatus"`
	TmName                   LBNodeStatsDescription `json:"tmName"`
	Serverside_bitsIn        LBStatsValue           `json:"serverside.bitsIn"`
	Serverside_bitsOut       LBStatsValue           `json:"serverside.bitsOut"`
	Serverside_curConns      LBStatsValue           `json:"serverside.curConns"`
	Serverside_maxConns      LBStatsValue           `json:"serverside.maxConns"`
	Serverside_pktsIn        LBStatsValue           `json:"serverside.pktsIn"`
	Serverside_pktsOut       LBStatsValue           `json:"serverside.pktsOut"`
	Serverside_totConns      LBStatsValue           `json:"serverside.totConns"`
	SessionStatus            LBNodeStatsDescription `json:"sessionStatus"`
	Status_availabilityState LBNodeStatsDescription `json:"status.availabilityState"`
	Status_enabledState      LBNodeStatsDescription `json:"status.enabledState"`
	Status_statusReason      LBNodeStatsDescription `json:"status.statusReason"`
	TotRequests              LBStatsValue           `json:"totRequests"`
}

type LBNodeStatsNestedStats struct {
	Kind     string                  `json:"kind"`
	SelfLink string                  `json:"selfLink"`
	Entries  LBNodeStatsInnerEntries `json:"entries"`
}

type LBNodeURLKey struct {
	NestedStats LBNodeStatsNestedStats `json:"nestedStats"`
}
type LBNodeStatsOuterEntries map[string]LBNodeURLKey

type LBNodeStats struct {
	Kind       string                  `json:"kind"`
	Generation int                     `json:"generation"`
	SelfLink   string                  `json:"selfLink"`
	Entries    LBNodeStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowNodes() (error, *LBNodes) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node"
	res := LBNodes{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowNode(nname string) (error, *LBNode) {

	//u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/~" + partition + "~" + pname + "?expandSubcollections=true"
	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := LBNode{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowNodeStats(nname string) (error, *LBObjectStats) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowAllNodeStats() (error, *LBNodeStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/stats"
	res := LBNodeStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddNode(body *json.RawMessage) (error, *LBNode) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node"
	res := LBNode{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateNode(nname string, body *json.RawMessage) (error, *LBNode) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := LBNode{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteNode(nname string) (error, *Response) {

	node := strings.Replace(nname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/node/" + node
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
