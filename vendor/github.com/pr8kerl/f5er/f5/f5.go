package f5

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	//sessn   napping.Session
	tsport            http.Transport
	clnt              http.Client
	headers           http.Header
	debug             bool
	tokenMutex        = sync.Mutex{}
	stats_path_prefix string
)

const (
	GET = iota
	POST
	POSTR
	PUT
	PUTR
	PATCH
	DELETE
)

type httperr struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

type Device struct {
	Hostname        string
	Username        string
	Password        string
	Session         napping.Session
	AuthToken       authToken
	AuthMethod      AuthMethod
	Proto           string
	StatsPathPrefix string
	StatsShowZeroes bool
}

type Response struct {
	Status  int
	Message string
}

type LBEmptyBody struct{}

type LBTransaction struct {
	TransId int    `json:"transId"`
	Timeout int    `json:"timeoutSeconds"`
	State   string `json:"state"`
}

type LBTransactionState struct {
	State string `json:"state"`
}

type AuthMethod int

const (
	TOKEN AuthMethod = iota
	BASIC_AUTH
)

type authToken struct {
	Token            string
	ExpirationMicros int64
}

func New(host string, username string, pwd string, authMethod AuthMethod) *Device {
	f := Device{Hostname: host, Username: username, Password: pwd, AuthMethod: authMethod, Proto: "https", StatsPathPrefix: "f5.", StatsShowZeroes: false}
	f.InitSession()
	return &f
}

func NewInsecure(host string, username string, pwd string, authMethod AuthMethod) *Device {
	f := Device{Hostname: host, Username: username, Password: pwd, AuthMethod: authMethod, Proto: "http", StatsPathPrefix: "f5.", StatsShowZeroes: false}
	f.InitSession()
	return &f
}

func (f *Device) InitSession() {

	// REST connection setup
	if f.Proto == "https" {
		tsport = http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		clnt = http.Client{Transport: &tsport}
	} else {
		clnt = http.Client{}
	}
	headers = make(http.Header)

	//
	// Setup HTTP Basic auth for this session (ONLY use this with SSL).  Auth
	// can also be configured on a per-request basis when using Send().
	//
	f.Session = napping.Session{
		Client: &clnt,
		Log:    debug,
		// if Userinfo is set - napping will set the basic auth header for you
		Userinfo: url.UserPassword(f.Username, f.Password),
		Header:   &headers,
	}

}

func (f *Device) SetDebug(b bool) {
	debug = b
}

func (f *Device) SetStatsPathPrefix(p string) {
	if strings.HasSuffix(p, ".") {
		f.StatsPathPrefix = p
	} else {
		f.StatsPathPrefix = p + "."
	}

}
func (f *Device) SetStatsShowZeroes(b bool) {
	f.StatsShowZeroes = b
}

func (f *Device) StartTransaction() (error, string) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/transaction"
	empty := LBEmptyBody{}
	tres := LBTransaction{}
	err, _ := f.sendRequest(u, POST, &empty, &tres)
	if err != nil {
		return err, ""
	}

	tid := fmt.Sprintf("%d", tres.TransId)
	// set the transaction header
	f.Session.Header.Set("X-F5-REST-Coordination-Id", tid)
	return nil, tid

}

func (f *Device) CommitTransaction(tid string) error {

	// remove the transaction header first
	f.Session.Header.Del("X-F5-REST-Coordination-Id")

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/transaction/" + tid
	body := LBTransaction{State: "VALIDATING"}
	tres := LBTransaction{}
	err, _ := f.sendRequest(u, PATCH, &body, &tres)
	if err != nil {
		return err
	}

	return nil

}

func (f *Device) sendRequest(u string, method int, pload interface{}, res interface{}) (error, *Response) {

	if f.AuthMethod == TOKEN {
		f.ensureValidToken()
	}

	//
	// Send request to server
	//
	e := httperr{}
	var (
		err   error
		nresp *napping.Response
	)
	f.Session.Log = debug

	switch method {
	case GET:
		nresp, err = f.Session.Get(u, nil, &res, &e)
	case POST:
		nresp, err = f.Session.Post(u, &pload, &res, &e)
	case PUT:
		nresp, err = f.Session.Put(u, &pload, &res, &e)
	case PATCH:
		nresp, err = f.Session.Patch(u, &pload, &res, &e)
	case DELETE:
		nresp, err = f.Session.Delete(u, nil, &res, &e)
	case POSTR:
		r := napping.Request{
			Method:     "POST",
			Url:        u,
			Params:     nil,
			Payload:    pload,
			RawPayload: true,
			Result:     res,
			Error:      e,
		}
		nresp, err = f.Session.Send(&r)
	case PUTR:
		r := napping.Request{
			Method:     "PUT",
			Url:        u,
			Params:     nil,
			Payload:    pload,
			RawPayload: true,
			Result:     &res,
			Error:      &e,
		}
		nresp, err = f.Session.Send(&r)
	}

	var resp Response
	if nresp != nil {
		resp = Response{Status: nresp.Status(), Message: e.Message}
	}

	if err != nil {
		return err, &resp
	}
	if nresp.Status() == 401 {
		f.PrintObject(resp)
		return errors.New("error: 401 Unauthorised - check your username and passwd"), &resp
	}
	if nresp.Status() >= 300 {
		return errors.New(e.Message), &resp
	} else {
		// all is good in the world
		return nil, &resp
	}
}

func (f *Device) PrintObject(input interface{}) {

	jsonresp, err := json.MarshalIndent(&input, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonresp))

}

// F5 Module data struct
// to show all available modules when using show without args
type LBModule struct {
	Link string `json:"link"`
}

type LBModuleRef struct {
	Reference LBModule `json:"reference"`
}

type LBModules struct {
	Items []LBModuleRef `json:"items"`
}

func (f *Device) ShowModules() (error, *LBModules) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/ltm"
	res := LBModules{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}
}

func (f *Device) GetToken() {

	type login struct {
		Token struct {
			Token            string `json:"token"`
			ExpirationMicros int64  `json:"expirationMicros"`
		} `json:"token"`
	}

	// Simply posting LoginData to the login endpoint doesn't seem to work.
	// I seem to need to set basic auth for the token request
	// after which I can disable basic auth by killing f.Session.Userinfo
	// I suspect this is a BIG-IP v11 issue - v12 docs clearly state no basic auth header is required
	// https://devcentral.f5.com/wiki/iControl.Authentication_with_the_F5_REST_API.ashx?lc=1
	// can't hurt for now
	if f.Session.Userinfo == nil {
		// turn on basic auth for this token request only
		f.Session.Userinfo = url.UserPassword(f.Username, f.Password)
	}

	// We need to remove X-F5-Auth-Token header when logging in because the BIG-IP
	// will look att it first and if it has expired it will return Unathorized
	f.Session.Header.Del("X-F5-Auth-Token")

	LoginData := map[string]string{"username": f.Username, "password": f.Password, "loginProviderName": "tmos"}
	byteLogin, err := json.Marshal(LoginData)
	body := json.RawMessage(byteLogin)
	u := f.Proto + "://" + f.Hostname + "/mgmt/shared/authn/login"
	res := login{}
	e := httperr{}

	resp, err := f.Session.Post(u, &body, &res, &e)
	if err != nil {
		log.Fatal(fmt.Errorf("error: %s, %v", err, resp))
		return
	}
	if debug {
		f.PrintObject(&resp)
	}

	f.AuthToken = authToken{
		Token:            res.Token.Token,
		ExpirationMicros: res.Token.ExpirationMicros,
	}
	f.Session.Header.Set("X-F5-Auth-Token", f.AuthToken.Token)

	// disable basic auth now
	f.Session.Userinfo = nil
}

func (f *Device) hasValidToken() bool {
	nowMicros := time.Now().UnixNano() / (int64(time.Microsecond) / int64(time.Nanosecond))
	if f.AuthToken.Token == "" || f.AuthToken.ExpirationMicros < nowMicros+int64(time.Millisecond)*100 {
		return false
	}
	return true
}

func (f *Device) ensureValidToken() {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	if !f.hasValidToken() {
		f.GetToken()
	}
}
