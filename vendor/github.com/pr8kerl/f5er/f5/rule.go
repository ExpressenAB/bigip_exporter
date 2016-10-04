package f5

import (
	"bytes"
	"encoding/json"
	"strings"
)

type LBRawValues struct {
	VerificationStatus string `json:"verificationStatus"`
}

type LBRule struct {
	Name         string      `json:"name"`
	Partition    string      `json:"partition"`
	FullPath     string      `json:"fullPath"`
	Generation   int         `json:"generation"`
	ApiAnonymous string      `json:"apiAnonymous"`
	ApiRawValues LBRawValues `json:"apiRawValues"`
}

type LBRules struct {
	Items []LBRule `json:"items"`
}

type LBRuleStatsDescription struct {
	Description string `json:"description"`
}

type LBRuleStatsValue struct {
	Value int `json:"value"`
}

type LBRuleStatsInnerEntries struct {
	Aborts          LBStatsValue           `json:"aborts"`
	AvgCycles       LBStatsValue           `json:"avgCycles"`
	EventType       LBRuleStatsDescription `json:"eventType"`
	Failures        LBStatsValue           `json:"failures"`
	MaxCycles       LBStatsValue           `json:"maxCycles"`
	MinCycles       LBStatsValue           `json:"minCycles"`
	TmName          LBRuleStatsDescription `json:"tmName"`
	Priority        LBStatsValue           `json:"priority"`
	TotalExecutions LBStatsValue           `json:"totalExecutions"`
}

type LBRuleStatsNestedStats struct {
	Kind     string                  `json:"kind"`
	SelfLink string                  `json:"selfLink"`
	Entries  LBRuleStatsInnerEntries `json:"entries"`
}

type LBRuleURLKey struct {
	NestedStats LBRuleStatsNestedStats `json:"nestedStats"`
}
type LBRuleStatsOuterEntries map[string]LBRuleURLKey

type LBRuleStats struct {
	Kind       string                  `json:"kind"`
	Generation int                     `json:"generation"`
	SelfLink   string                  `json:"selfLink"`
	Entries    LBRuleStatsOuterEntries `json:"entries"`
}

func (f *Device) ShowRules() (error, *LBRules) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule"
	res := LBRules{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowRule(rname string) (error, *LBRule) {

	rule := strings.Replace(rname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/" + rule
	res := LBRule{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowRuleStats(rname string) (error, *LBObjectStats) {

	rule := strings.Replace(rname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/" + rule + "/stats"
	res := LBObjectStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowAllRuleStats() (error, *LBRuleStats) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/stats"
	res := LBRuleStats{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddRule(body *json.RawMessage) (error, *LBRule) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule"
	res := LBRule{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddRuleRaw(body *bytes.Buffer) (error, *LBRule) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule"
	res := LBRule{}

	// post the request
	err, _ := f.sendRequest(u, POSTR, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateRule(rname string, body *json.RawMessage) (error, *LBRule) {

	rule := strings.Replace(rname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/" + rule
	res := LBRule{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateRuleRaw(rname string, body *bytes.Buffer) (error, *LBRule) {

	rule := strings.Replace(rname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/" + rule
	res := LBRule{}

	// put the request
	err, _ := f.sendRequest(u, PUTR, body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteRule(rname string) (error, *Response) {

	rule := strings.Replace(rname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/rule/" + rule
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
