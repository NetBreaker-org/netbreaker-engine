package configInterfaces

type HostnameServerArrayData struct {
	Address        string   `json:"address"`
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

type InstanceServerConfig struct {
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

// type for instanceList
type InstancesJson struct {
	Name   string `json:"instance_identifier"`
	Dir    string `json:"instance_dir"`
	Schema string `json:"schema_used"`
}

// Index of all config objects
// exported as a variable
type ConfigObjectPointers struct {
	ServerConfig      *ServerConfig `json:"ServerConfig"`
	ServerConfigFLPTH string        `json:"ServerConfigFLPTH"`

	InstanceList      *[]InstancesJson `json:"InstanceList"`
	InstanceListFLPTH string           `json:"InstanceListFLPTH"`

	IPBlackList      *[]string `json:"ip_bl"`
	IPBlackListFLPTH string    `json:"IPBlackListFLPTH"`

	IPWhiteList      *[]string `json:"ip_wl"`
	IPWhiteListFLPTH string    `json:"IPWhiteListFLPTH"`
}

// Index for config object pointers for all instances
type InstanceConfigObjectPointers struct {
	InstanceServerConfig      *InstanceServerConfig `json:"ServerConfig"`
	InstanceServerConfigFLPTH string                `json:"ServerConfigFLPTH"`

	InstanceIPBlackList      *[]string `json:"ip_bl"`
	InstanceIPBlackListFLPTH string    `json:"IPBlackListFLPTH"`

	InstanceIPWhiteList      *[]string `json:"ip_wl"`
	InstanceIPWhiteListFLPTH string    `json:"IPWhiteListFLPTH"`

	InstanceIPMap                   *[]string `json:"ip_map"`
	InstanceAvailableHostnamesFLPTH string    `json:"AvailableHostnamesFLPTH"`
	InstanceIPMapFLPTH              string    `json:"IPMapFLPTH"`

	InstanceHostnameMapping      *[]StatusMapping `json:"hostnameMapping_map"`
	InstanceHostnameMappingFLPTH string           `json:"hostnameMappingFLPTH"`
}

// Type for array of all instances with their data
type InstanceConfigs struct {
	Metadata InstancesJson                `json:"InstanceMetadata"`
	Config   InstanceConfigObjectPointers `json:"InstanceConfig"`
}
