package f5

import (
	"encoding/json"
	"strings"
)

// a pool member
type LBPoolMember struct {
	Name            string `json:"name"`
	Partition       string `json:"partition"`
	FullPath        string `json:"fullPath"`
	Address         string `json:"address"`
	ConnectionLimit int    `json:"connectionLimit"`
	DynamicRatio    int    `json:"dynamicRatio"`
	Ephemeral       string `json:"ephemeral"`
	InheritProfile  string `json:"inheritProfile"`
	Logging         string `json:"logging"`
	Monitor         string `json:"monitor"`
	PriorityGroup   int    `json:"priorityGroup"`
	RateLimit       string `json:"rateLimit"`
	Ratio           int    `json:"ratio"`
	Session         string `json:"session"`
	State           string `json:"state"`
}

// a pool member reference - just a link and an array of pool members
type LBPoolMemberRef struct {
	Link  string         `json:"link"`
	Items []LBPoolMember `json":items"`
}

type LBPoolMembers struct {
	Link  string         `json:"selfLink"`
	Items []LBPoolMember `json":items"`
}

// used by online/offline
type LBPoolMemberState struct {
	State   string `json:"state"`
	Session string `json:"session"`
}

type LBPool struct {
	Name                   string          `json:"name"`
	Partition              string          `json:"partition"`
	FullPath               string          `json:"fullPath"`
	Generation             int             `json:"generation"`
	AllowNat               string          `json:"allowNat"`
	AllowSnat              string          `json:"allowSnat"`
	IgnorePersistedWeight  string          `json:"ignorePersistedWeight"`
	IpTosToClient          string          `json:"ipTosToClient"`
	IpTosToServer          string          `json:"ipTosToServer"`
	LinkQosToClient        string          `json:"linkQosToClient"`
	LinkQosToServer        string          `json:"linkQosToServer"`
	LoadBalancingMode      string          `json:"loadBalancingMode"`
	MinActiveMembers       int             `json:"minActiveMembers"`
	MinUpMembers           int             `json:"minUpMembers"`
	MinUpMembersAction     string          `json:"minUpMembersAction"`
	MinUpMembersChecking   string          `json:"minUpMembersChecking"`
	Monitor                string          `json:"monitor"`
	QueueDepthLimit        int             `json:"queueDepthLimit"`
	QueueOnConnectionLimit string          `json:"queueOnConnectionLimit"`
	QueueTimeLimit         int             `json:"queueTimeLimit"`
	ReselectTries          int             `json:"reselectTries"`
	ServiceDownAction      string          `json:"serviceDownAction"`
	SlowRampTime           int             `json:"slowRampTime"`
	MemberRef              LBPoolMemberRef `json:"membersReference"`
}

type LBPools struct {
	Items []LBPool `json:"items"`
}

type LBPoolStatsDescription struct {
	Description string `json:"description"`
}

type LBPoolStatsValue struct {
	Value int `json:"value"`
}

type LBPoolStatsInnerEntries struct {
	ActiveMemberCnt          LBStatsValue           `json:"activeMemberCnt"`
	ConnqAll_ageEdm          LBStatsValue           `json:"connqAll.ageEdm"`
	ConnqAll_ageEma          LBStatsValue           `json:"connqAll.ageEma"`
	ConnqAll_ageHead         LBStatsValue           `json:"connqAll.ageHead"`
	ConnqAll_ageMax          LBStatsValue           `json:"connqAll.ageMax"`
	ConnqAll_depth           LBStatsValue           `json:"connqAll.depth"`
	ConnqAll_serviced        LBStatsValue           `json:"connqAll.serviced"`
	Connq_ageEdm             LBStatsValue           `json:"connq.ageEdm"`
	Connq_ageEma             LBStatsValue           `json:"connq.ageEma"`
	Connq_ageHead            LBStatsValue           `json:"connq.ageHead"`
	Connq_ageMax             LBStatsValue           `json:"connq.ageMax"`
	Connq_depth              LBStatsValue           `json:"connq.depth"`
	Connq_serviced           LBStatsValue           `json:"connq.serviced"`
	CurSessions              LBStatsValue           `json:"curSessions"`
	MinActiveMembers         LBStatsValue           `json:"minActiveMembers"`
	MonitorRule              LBPoolStatsDescription `json:"monitorRule"`
	TmName                   LBPoolStatsDescription `json:"tmName"`
	Serverside_bitsIn        LBStatsValue           `json:"serverside.bitsIn"`
	Serverside_bitsOut       LBStatsValue           `json:"serverside.bitsOut"`
	Serverside_curConns      LBStatsValue           `json:"serverside.curConns"`
	Serverside_maxConns      LBStatsValue           `json:"serverside.maxConns"`
	Serverside_pktsIn        LBStatsValue           `json:"serverside.pktsIn"`
	Serverside_pktsOut       LBStatsValue           `json:"serverside.pktsOut"`
	Serverside_totConns      LBStatsValue           `json:"serverside.totConns"`
	Status_availabilityState LBPoolStatsDescription `json:"status.availabilityState"`
	Status_enabledState      LBPoolStatsDescription `json:"status.enabledState"`
	Status_statusReason      LBPoolStatsDescription `json:"status.statusReason"`
	TotRequests              LBStatsValue           `json:"totRequests"`
}

type LBPoolStatsNestedStats struct {
	Kind     string                  `json:"kind"`
	SelfLink string                  `json:"selfLink"`
	Entries  LBPoolStatsInnerEntries `json:"entries"`
}

type LBPoolURLKey struct {
	NestedStats LBPoolStatsNestedStats `json:"nestedStats"`
}
type LBPoolStatsOuterEntries map[string]LBPoolURLKey

type LBPoolStats struct {
	Kind       string                  `json:"kind"`
	Generation int                     `json:"generation"`
	SelfLink   string                  `json:"selfLink"`
	Entries    LBPoolStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowPools() (error, *LBPools) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool"
	res := LBPools{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPool(pname string) (error, *LBPool) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "?expandSubcollections=true"
	res := LBPool{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPoolStats(pname string) (error, *LBObjectStats) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowAllPoolStats() (error, *LBPoolStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/stats"
	res := LBPoolStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) AddPool(body *json.RawMessage) (error, *LBPool) {
	// we use json.RawMessage so we can modify the input file without using a struct
	// use of a struct will send all available fields, some of which can't be modified

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool"
	res := LBPool{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePool(pname string, body *json.RawMessage) (error, *LBPool) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool
	res := LBPool{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeletePool(pname string) (error, *Response) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) ShowPoolMembers(pname string) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddPoolMembers(pname string, body *json.RawMessage) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePoolMembers(pname string, body *json.RawMessage) (error, *LBPoolMembers) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := LBPoolMembers{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeletePoolMembers(pname string) (error, *Response) {

	pool := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members"
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) OnlinePoolMember(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-up", "user-enabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}

func (f *Device) OfflinePoolMember(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-up", "user-disabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
func (f *Device) OfflinePoolMemberForced(pname string, mname string) (error, *Response) {

	pmember := strings.Replace(mname, "/", "~", -1)
	pool := strings.Replace(pname, "/", "~", -1)

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/pool/" + pool + "/members/" + pmember
	res := json.RawMessage{}

	/*
	   {"state": "user-down", "session": "user-disabled"} (Member Forced Offline in GUI)
	   {"state": "user-up", "session": "user-disabled"} (Member Disabled in GUI)
	   {"state": "user-up", "session": "user-enabled"}  (Member Enabled in GUI)
	*/
	body := LBPoolMemberState{"user-down", "user-disabled"}

	// put the request
	err, resp := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
