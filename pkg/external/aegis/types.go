package aegis

type ConfigGetResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ConfigMap map[string]interface{} `json:"config_map"` // ConfigMap
	} `json:"data"`
}
