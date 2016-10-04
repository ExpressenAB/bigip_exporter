package f5

import (
	"encoding/json"
	"strings"
)

type LBServerSsl struct {
	Name                         string   `json:"name"`
	Partition                    string   `json:"partition"`
	FullPath                     string   `json:"fullPath"`
	Generation                   int      `json:"generation"`
	UntrustedCertResponseControl string   `json:"untrustedCertResponseControl"`
	UncleanShutdown              string   `json:"uncleanShutdown"`
	StrictResume                 string   `json:"strictResume"`
	SslSignHash                  string   `json:"sslSignHash"`
	SslForwardProxyBypass        string   `json:"sslForwardProxyBypass"`
	SslForwardProxy              string   `json:"sslForwardProxy"`
	SniRequire                   string   `json:"sniRequire"`
	SniDefault                   string   `json:"sniDefault"`
	ExpireCertResponseControl    string   `json:"expireCertResponseControl"`
	DefaultsFrom                 string   `json:"defaultsFrom"`
	Ciphers                      string   `json:"ciphers"`
	Chain                        string   `json:"chain"`
	Cert                         string   `json:"cert"`
	Key                          string   `json:"key"`
	CacheTimeout                 int      `json:"cacheTimeout"`
	CacheSize                    int      `json:"cacheSize"`
	AuthenticateDepth            int      `json:"authenticateDepth"`
	AlertTimeout                 string   `json:"alertTimeout"`
	SelfLink                     string   `json:"selfLink"`
	Authenticate                 string   `json:"authenticate"`
	GenericAlert                 string   `json:"genericAlert"`
	HandshakeTimeout             string   `json:"handshakeTimeout"`
	ModSslMethods                string   `json:"modSslMethods"`
	Mode                         string   `json:"mode"`
	TmOptions                    []string `json:"tmOptions"`
	PeerCertMode                 string   `json:"peerCertMode"`
	ProxySsl                     string   `json:"proxySsl"`
	ProxySslPassthrough          string   `json:"proxySslPassthrough"`
	RenegotiatePeriod            string   `json:"renegotiatePeriod"`
	RenegotiateSize              string   `json:"renegotiateSize"`
	Renegotiation                string   `json:"renegotiation"`
	RetainCertificate            string   `json:"retainCertificate"`
	SecureRenegotiation          string   `json:"secureRenegotiation"`
	SessionMirroring             string   `json:"sessionMirroring"`
	SessionTicket                string   `json:"sessionTicket"`
}

type LBServerSsls struct {
	Items []LBServerSsl `json:"items"`
}

func (f *Device) ShowServerSsls() (error, *LBServerSsls) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/server-ssl"
	res := LBServerSsls{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

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
