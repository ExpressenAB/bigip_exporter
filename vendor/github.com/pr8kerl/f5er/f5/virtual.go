package f5

import (
	"encoding/json"
	"strings"
)

type LBVirtualPolicy struct {
	Name      string `json:"name"`
	Partition string `json:"partition"`
	FullPath  string `json:"fullPath"`
}

type LBVirtualPoliciesRef struct {
	Items []LBVirtualPolicy `json":items"`
}

type LBVirtualProfile struct {
	Name      string `json:"name"`
	Partition string `json:"partition"`
	FullPath  string `json:"fullPath"`
	Context   string `json:"context"`
}

type LBVirtualPersistProfile struct {
	Name      string `json:"name"`
	Partition string `json:"partition"`
	TmDefault string `json:"tmDefault"`
}

type LBVirtualProfileRef struct {
	Items []LBVirtualProfile `json":items"`
}

type LBVirtual struct {
	Name             string                    `json:"name"`
	FullPath         string                    `json:"fullPath"`
	Partition        string                    `json:"partition"`
	Destination      string                    `json:"destination"`
	Pool             string                    `json:"pool"`
	AddressStatus    string                    `json:"addressStatus"`
	AutoLastHop      string                    `json:"autoLasthop"`
	CmpEnabled       string                    `json:"cmpEnabled"`
	ConnectionLimit  int                       `json:"connectionLimit"`
	Enabled          bool                      `json:"enabled"`
	IpProtocol       string                    `json:"ipProtocol"`
	Source           string                    `json:"source"`
	SourcePort       string                    `json:"sourcePort"`
	SynCookieStatus  string                    `json:"synCookieStatus"`
	TranslateAddress string                    `json:"translateAddress"`
	TranslatePort    string                    `json:"translatePort"`
	Profiles         LBVirtualProfileRef       `json:"profilesReference"`
	Policies         LBVirtualPoliciesRef      `json:"policiesReference"`
	Rules            []string                  `json:"rules"`
	Persist          []LBVirtualPersistProfile `json:"persist"`
}

type LBVirtuals struct {
	Items []LBVirtual
}

type LBVirtualStatsDescription struct {
	Description string `json:"description"`
}

type LBVirtualStatsValue struct {
	Value int `json:"value"`
}

type LBVirtualStatsInnerEntries struct {
	Clientside_bitsIn             LBStatsValue              `json:"clientside.bitsIn"`
	Clientside_bitsOut            LBStatsValue              `json:"clientside.bitsOut"`
	Clientside_curConns           LBStatsValue              `json:"clientside.curConns"`
	Clientside_evictedConns       LBStatsValue              `json:"clientside.evictedConns"`
	Clientside_maxConns           LBStatsValue              `json:"clientside.maxConns"`
	Clientside_pktsIn             LBStatsValue              `json:"clientside.pktsIn"`
	Clientside_pktsOut            LBStatsValue              `json:"clientside.pktsOut"`
	Clientside_slowKilled         LBStatsValue              `json:"clientside.slowKilled"`
	Clientside_totConns           LBStatsValue              `json:"clientside.totConns"`
	CmpEnableMode                 LBVirtualStatsDescription `json:"cmpEnableMode"`
	CmpEnabled                    LBVirtualStatsDescription `json:"cmpEnabled"`
	CsMaxConnDur                  LBStatsValue              `json:"csMaxConnDur"`
	CsMeanConnDur                 LBStatsValue              `json:"csMeanConnDur"`
	CsMinConnDur                  LBStatsValue              `json:"csMinConnDur"`
	Destination                   LBVirtualStatsDescription `json:"destination"`
	Ephemeral_bitsIn              LBStatsValue              `json:"ephemeral.bitsIn"`
	Ephemeral_bitsOut             LBStatsValue              `json:"ephemeral.bitsOut"`
	Ephemeral_curConns            LBStatsValue              `json:"ephemeral.curConns"`
	Ephemeral_evictedConns        LBStatsValue              `json:"ephemeral.evictedConns"`
	Ephemeral_maxConns            LBStatsValue              `json:"ephemeral.maxConns"`
	Ephemeral_pktsIn              LBStatsValue              `json:"ephemeral.pktsIn"`
	Ephemeral_pktsOut             LBStatsValue              `json:"ephemeral.pktsOut"`
	Ephemeral_slowKilled          LBStatsValue              `json:"ephemeral.slowKilled"`
	Ephemeral_totConns            LBStatsValue              `json:"ephemeral.totConns"`
	FiveMinAvgUsageRatio          LBStatsValue              `json:"fiveMinAvgUsageRatio"`
	FiveSecAvgUsageRatio          LBStatsValue              `json:"fiveSecAvgUsageRatio"`
	TmName                        LBVirtualStatsDescription `json:"tmName"`
	OneMinAvgUsageRatio           LBStatsValue              `json:"oneMinAvgUsageRatio"`
	Status_availabilityState      LBVirtualStatsDescription `json:"status.availabilityState"`
	Status_enabledState           LBVirtualStatsDescription `json:"status.enabledState"`
	Status_statusReason           LBVirtualStatsDescription `json:"status.statusReason"`
	SyncookieStatus               LBVirtualStatsDescription `json:"syncookieStatus"`
	Syncookie_accepts             LBStatsValue              `json:"syncookie.accepts"`
	Syncookie_hwAccepts           LBStatsValue              `json:"syncookie.hwAccepts"`
	Syncookie_hwSyncookies        LBStatsValue              `json:"syncookie.hwSyncookies"`
	Syncookie_hwsyncookieInstance LBStatsValue              `json:"syncookie.hwsyncookieInstance"`
	Syncookie_rejects             LBStatsValue              `json:"syncookie.rejects"`
	Syncookie_swsyncookieInstance LBStatsValue              `json:"syncookie.swsyncookieInstance"`
	Syncookie_syncacheCurr        LBStatsValue              `json:"syncookie.syncacheCurr"`
	Syncookie_syncacheOver        LBStatsValue              `json:"syncookie.syncacheOver"`
	Syncookie_syncookies          LBStatsValue              `json:"syncookie.syncookies"`
	TotRequests                   LBStatsValue              `json:"totRequests"`
}

type LBVirtualStatsNestedStats struct {
	Kind     string                     `json:"kind"`
	SelfLink string                     `json:"selfLink"`
	Entries  LBVirtualStatsInnerEntries `json:"entries"`
}

type LBVirtualURLKey struct {
	NestedStats LBVirtualStatsNestedStats `json:"nestedStats"`
}
type LBVirtualStatsOuterEntries map[string]LBVirtualURLKey

type LBVirtualStats struct {
	Kind       string                     `json:"kind"`
	Generation int                        `json:"generation"`
	SelfLink   string                     `json:"selfLink"`
	Entries    LBVirtualStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowVirtuals() (error, *LBVirtuals) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual"
	res := LBVirtuals{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowVirtual(vname string) (error, *LBVirtual) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname + "?expandSubcollections=true"
	res := LBVirtual{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowVirtualStats(vname string) (error, *LBObjectStats) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowAllVirtualStats() (error, *LBVirtualStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/stats"
	res := LBVirtualStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) AddVirtual(virt *json.RawMessage) (error, *LBVirtual) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual"
	res := LBVirtual{}

	// post the request
	err, _ := f.sendRequest(u, POST, virt, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateVirtual(vname string, body *json.RawMessage) (error, *LBVirtual) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname
	res := LBVirtual{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteVirtual(vname string) (error, *Response) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/virtual/" + vname
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
