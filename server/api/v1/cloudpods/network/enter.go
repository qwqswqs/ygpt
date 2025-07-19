package network

type NetApiGroup struct {
	NetIpApi
	NetworkVpcApi
	NetworkEipApi
	NetworkSecgroupApi
}
type ErrorResponse struct {
	Error struct {
		Class   string `json:"class"`
		Code    int    `json:"code"`
		Details string `json:"details"`
		Request struct {
			Body    string            `json:"body"`
			Headers map[string]string `json:"headers"`
			Method  string            `json:"method"`
			URL     string            `json:"url"`
		} `json:"request"`
	} `json:"error"`
}
