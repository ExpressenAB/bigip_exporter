package f5

import "encoding/json"

type LBProfileRef struct {
	Link string `json:"link"`
}

type LBProfile struct {
	Reference LBProfileRef `json:"reference"`
}

type LBProfiles struct {
	Items []LBProfile `json:"items"`
}

func (f *Device) ShowProfiles() (error, *LBProfiles) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile"
	res := LBProfiles{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) ShowProfile(profile string) (error, *json.RawMessage) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/" + profile
	res := json.RawMessage{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

/*
func (f *Device) ShowServerSsl(sname string) (error, *LBServerSsl) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := LBServerSsl{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddServerSsl(body *json.RawMessage) (error, *LBServerSsl) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl"
	res := LBServerSsl{}

	// post the request
	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateServerSsl(sname string, body *json.RawMessage) (error, *LBServerSsl) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := LBServerSsl{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteServerSsl(sname string) (error, *Response) {

	server := strings.Replace(sname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl/" + server
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
*/
