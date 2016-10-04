package f5

import (
	"encoding/json"
	"strings"
)

type LBClientSsl struct {
	Name                            string           `json:"name"`
	Partition                       string           `json:"partition"`
	FullPath                        string           `json:"fullPath"`
	Generation                      int              `json:"generation"`
	AlertTimeout                    string           `json:"alertTimeout"`
	AllowNonSsl                     string           `json:"allowNonSsl"`
	Authenticate                    string           `json:"authenticate"`
	AuthenticateDepth               int              `json:"authenticateDepth"`
	CacheSize                       int              `json:"cacheSize"`
	CacheTimeout                    int              `json:"cacheTimeout"`
	Cert                            string           `json:"cert"`
	CertExtensionIncludes           []string         `json:"certExtensionIncludes"`
	CertLifespan                    int              `json:"certLifespan"`
	CertLookupByIpaddrPort          string           `json:"certLookupByIpaddrPort"`
	Chain                           string           `json:"chain"`
	Ciphers                         string           `json:"ciphers"`
	DefaultsFrom                    string           `json:"defaultsFrom"`
	ForwardProxyBypassDefaultAction string           `json:"forwardProxyBypassDefaultAction"`
	GenericAlert                    string           `json:"genericAlert"`
	HandshakeTimeout                string           `json:"handshakeTimeout"`
	InheritCertkeychain             string           `json:"inheritCertkeychain"`
	Key                             string           `json:"key"`
	MaxRenegotiationsPerMinute      int              `json:"maxRenegotiationsPerMinute"`
	ModSslMethods                   string           `json:"modSslMethods"`
	Mode                            string           `json:"mode"`
	TmOptions                       []string         `json:"tmOptions"`
	PeerCertMode                    string           `json:"peerCertMode"`
	PeerNoRenegotiateTimeout        string           `json:"peerNoRenegotiateTimeout"`
	ProxySsl                        string           `json:"proxySsl"`
	ProxySslPassthrough             string           `json:"proxySslPassthrough"`
	RenegotiateMaxRecordDelay       string           `json:"renegotiateMaxRecordDelay"`
	RenegotiatePeriod               string           `json:"renegotiatePeriod"`
	RenegotiateSize                 string           `json:"renegotiateSize"`
	Renegotiation                   string           `json:"renegotiation"`
	RetainCertificate               string           `json:"retainCertificate"`
	SecureRenegotiation             string           `json:"secureRenegotiation"`
	SessionMirroring                string           `json:"sessionMirroring"`
	SessionTicket                   string           `json:"sessionTicket"`
	SniDefault                      string           `json:"sniDefault"`
	SniRequire                      string           `json:"sniRequire"`
	SslForwardProxy                 string           `json:"sslForwardProxy"`
	SslForwardProxyBypass           string           `json:"sslForwardProxyBypass"`
	SslSignHash                     string           `json:"sslSignHash"`
	StrictResume                    string           `json:"strictResume"`
	UncleanShutdown                 string           `json:"uncleanShutdown"`
	CertKeyChain                    []LBCertKeyChain `json:"certKeyChain"`
}

type LBCertKeyChain struct {
	Name  string `json:"name"`
	Cert  string `json:"cert"`
	Chain string `json:"chain"`
	Key   string `json:"key"`
}

type LBClientSsls struct {
	Items []LBClientSsl `json:"items"`
}

func (f *Device) ShowClientSsls() (error, *LBClientSsls) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/client-ssl"
	res := LBClientSsls{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) ShowClientSsl(cname string) (error, *LBClientSsl) {

	client := strings.Replace(cname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/client-ssl/" + client
	res := LBClientSsl{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) AddClientSsl(body *json.RawMessage) (error, *LBClientSsl) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/client-ssl"
	res := LBClientSsl{}

	err, _ := f.sendRequest(u, POST, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) UpdateClientSsl(cname string, body *json.RawMessage) (error, *LBClientSsl) {

	client := strings.Replace(cname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/client-ssl/" + client
	res := LBClientSsl{}

	// put the request
	err, _ := f.sendRequest(u, PUT, &body, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}

func (f *Device) DeleteClientSsl(cname string) (error, *Response) {

	client := strings.Replace(cname, "/", "~", -1)
	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm/profile/client-ssl/" + client
	res := json.RawMessage{}

	err, resp := f.sendRequest(u, DELETE, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, resp
	}

}
