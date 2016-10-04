package f5

import (
	"encoding/json"
	"strings"
)

type LBMonitorHttp struct {
	Name                     string `json:"name"`
	Partition                string `json:"partition"`
	FullPath                 string `json:"fullPath"`
	Adaptive                 string `json:"adaptive"`
	AdaptiveDivergenceType   string `json:"adaptiveDivergenceType"`
	AdaptiveDivergenceValue  int    `json:"adaptiveDivergenceValue"`
	AdaptiveLimit            int    `json:"adaptiveLimit"`
	AdaptiveSamplingTimespan int    `json:"adaptiveSamplingTimespan"`
	DefaultsFrom             string `json:"defaultsFrom"`
	Destination              string `json:"destination"`
	Interval                 int    `json:"interval"`
	IpDscp                   int    `json:"ipDscp"`
	ManualResume             string `json:"manualResume"`
	Recv                     string `json:"recv"`
	Reverse                  string `json:"reverse"`
	Send                     string `json:"send"`
	TimeUntilUp              int    `json:"timeUntilUp"`
	Timeout                  int    `json:"timeout"`
	Transparent              string `json:"transparent"`
	UpInterval               int    `json:"upInterval"`
}

type LBMonitorHttpRef struct {
	Items []LBMonitorHttp `json":items"`
}

func (f *Device) ShowMonitorsHttp() (error, *LBMonitorHttpRef) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http"
	res := LBMonitorHttpRef{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowMonitorHttp(vname string) (error, *LBMonitorHttp) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname + "?expandSubcollections=true"
	res := LBMonitorHttp{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddMonitorHttp(body *json.RawMessage) (error, *LBMonitorHttp) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http"
	res := LBMonitorHttp{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) UpdateMonitorHttp(vname string, body *json.RawMessage) (error, *LBMonitorHttp) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname
	res := LBMonitorHttp{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteMonitorHttp(vname string) (error, *Response) {

	vname = strings.Replace(vname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/monitor/http/" + vname
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
