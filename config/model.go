package config

type ConfigServer struct {
	Port   int    `json:"port"`
	Debug  bool   `json:"debug"`
	Status string `json:"status"`
}

type ConfigGRpc struct {
	Port int `json:"port"`
}

type ConfigRedis struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	Passwd string `json:"passwd"`
	DB     int    `json:"db"`
}

type ConfigMySQL struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DBName string `json:"dbname"`
}

type ConfigFOne struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Passwd  string `json:"passwd"`
	Timeout int    `json:"timeout"`
}

type ConfigNameServer []string

type ConfigSaml struct {
	Domain                      string `json:"domain"`
	IDPSSOURL                   string `json:"idpssourl"`
	IDPSSODescriptorURL         string `json:"idpsso_descriptor_url"`
	AssertionConsumerServiceURL string `json:"assertion_consumer_service_url"`
}
