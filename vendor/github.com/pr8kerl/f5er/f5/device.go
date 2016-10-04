package f5

type LBDeviceRef struct {
	Link  string          `json:"selfLink"`
	Items []LBDeviceState `json":items"`
}

type LBDeviceState struct {
	Name          string `json:"name"`
	Path          string `json:"fullPath"`
	FailoverState string `json:"failoverState"`
	ManagementIP  string `json:"managementIP"`
}

func (f *Device) ShowDevice() (error, *LBDeviceRef) {

	u := f.Proto + "://" + f.Hostname + "/mgmt/tm/cm/device"
	res := LBDeviceRef{}

	err, _ := f.sendRequest(u, GET, nil, &res)
	if err != nil {
		return err, nil
	} else {
		return nil, &res
	}

}
