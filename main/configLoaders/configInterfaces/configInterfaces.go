package configInterfaces

type HostnameServerArrayData struct {
	Address        int      `json:"address"`
	Status         int      `json:"status"`
	LoadAssignment []string `json:"loadAssignment"`
}

// Miscellaneous data of the server instances running, proxied by a specific proxy server
type MiscellaneousData struct {
	Views int `json:"views"`
	Reqs  int `json:"reqs"`
}

// Data of the server instances running, proxied by a specific proxy server, to be used as an array []StatusMapping
type StatusMapping struct {
	Hostname          string                    `json:"hostname"`
	StatusData        []HostnameServerArrayData `json:"StatusData"`
	MiscellaneousData MiscellaneousData         `json:"MiscellaneousData"`
}

// Data of the proxy servers running somewhere
type ProxyServices struct {
	Address     string `json:"proxy_server_address"`
	ProxyEngine string `json:"proxy_engine"`
	Version     string `json:"version"`
}

// Array of hostnames available
type IPMap struct {
	Key []string `json:"key"`
}

// Data struct for main config
type ServerConfig struct {
	Status                              string `json:"status"`
	Port                                string `json:"port"`
	Cipherkey                           string `json:"cipherkey"`
	CookieEncryptOption                 string `json:"cookieEncryptOption"`
	MaxRequestRateLimit                 string `json:"maxRequestRateLimit"`
	InspectOriginMode                   string `json:"inspectOriginMode"`
	AdminMode                           string `json:"adminMode"`
	AdminDefinedPath                    string `json:"adminDefinedPath"`
	AdminDefinedPathRequiredReloadCount string `json:"adminDefinedPathRequiredReloadCount"`
	AdminDefinedPathReloadTimeWindow    string `json:"adminDefinedPathReloadTimeWindow"`
	Protocol                            string `json:"protocol"`
	SSLKeyPath                          string `json:"sslKeyPath"`
	SSLCertPath                         string `json:"sslCertPath"`
	ErrorPagePort                       string `json:"errorPagePort"`
	Workers                             string `json:"workers"`
}

// Index of all config objects
type ConfigObjectPointers struct {
	ServerConfig      *ServerConfig `json:"ServerConfig"`
	ServerConfigFLPTH string        `json:"ServerConfigFLPTH"`

	IPBlackList      *[]string `json:"ip_bl"`
	IPBlackListFLPTH string    `json:"IPBlackListFLPTH"`

	IPWhiteList      *[]string `json:"ip_wl"`
	IPWhiteListFLPTH string    `json:"IPWhiteListFLPTH"`
}
