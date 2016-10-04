package f5

import (
	"encoding/json"
	"strings"
)

type LBPolicyConditions struct {
	Name            string   `json:"name"`
	FullPath        string   `json:"fullPath"`
	Generation      int      `json:"generation"`
	All             bool     `json:"all"`
	CaseInsensitive bool     `json:"caseInsensitive"`
	External        bool     `json:"external"`
	HttpUri         bool     `json:"httpUri"`
	Index           int      `json:"index"`
	Present         bool     `json:"present"`
	Remote          bool     `json:"remote"`
	Request         bool     `json:"request"`
	StartsWith      bool     `json:"startsWith"`
	Values          []string `json:"values"`
}

type LBPolicyActions struct {
	Name       string `json:"name"`
	FullPath   string `json:"fullPath"`
	Generation int    `json:"generation"`
	Code       int    `json:"code"`
	Forward    bool   `json:"forward"`
	Pool       string `json:"pool"`
	Port       int    `json:"port"`
	Request    bool   `json:"request"`
	Select     bool   `json:"select"`
	Status     int    `json:"status"`
	VlanId     int    `json:"vlanId"`
}

type LBPolicyConditionsRef struct {
	Items []LBPolicyConditions `json:"items"`
}

type LBPolicyActionsRef struct {
	Items []LBPolicyActions `json:"items"`
}

type LBPolicyRules struct {
	Name          string                `json:"name"`
	FullPath      string                `json:"fullPath"`
	Generation    int                   `json:"generation"`
	Ordinal       int                   `json:"ordinal"`
	ActionsRef    LBPolicyActionsRef    `json:"actionsReference"`
	ConditionsRef LBPolicyConditionsRef `json:"conditionsReference"`
}

type LBPolicyRulesRef struct {
	Items []LBPolicyRules `json:"items"`
}

type LBPolicy struct {
	Name       string           `json:"name"`
	Partition  string           `json:"partition"`
	FullPath   string           `json:"fullPath"`
	Generation int              `json:"generation"`
	Controls   []string         `json:"controls"`
	Requires   []string         `json:"requires"`
	Strategy   string           `json:"strategy"`
	RulesRef   LBPolicyRulesRef `json:"rulesReference"`
}

type LBPolicies struct {
	Items []LBPolicy `json:"items"`
}

func (f *Device) ShowPolicies() (error, *LBPolicies) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy"
	res := LBPolicies{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowPolicy(pname string) (error, *LBPolicy) {

	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy + "?expandSubcollections=true"
	res := LBPolicy{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddPolicy(body *json.RawMessage) (error, *LBPolicy) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy"
	res := LBPolicy{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdatePolicy(pname string, body *json.RawMessage) (error, *LBPolicy) {

	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy
	res := LBPolicy{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeletePolicy(pname string) (error, *Response) {

	//u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/~" + partition + "~" + pname + "?expandSubcollections=true"
	policy := strings.Replace(pname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/policy/" + policy
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
