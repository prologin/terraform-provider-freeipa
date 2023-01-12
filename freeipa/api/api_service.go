package api

type Service struct {
	KrbCanonicalName []string `json:"krbcanonicalname"` // Service name
}

func (c *APIClient) ServiceAdd(krbcanonicalname string, options JSON) (*Service, error) {
	return apiRequest[Service, string](c, "service_add", options, krbcanonicalname)
}

func (c *APIClient) ServiceDel(krbcanonicalname string, options JSON) (*JSON, error) {
	return apiRequest[JSON, []string](c, "service_del", options, krbcanonicalname)
}

func (c *APIClient) ServiceMod(krbcanonicalname string, options JSON) (*Service, error) {
	return apiRequest[Service, string](c, "service_mod", options, krbcanonicalname)
}

func (c *APIClient) ServiceShow(krbcanonicalname string, options JSON) (*Service, error) {
	return apiRequest[Service, string](c, "service_show", options, krbcanonicalname)
}

func (c *APIClient) ServiceFind(criteria string, options JSON) (*[]Service, error) {
	return apiRequest[[]Service, string](c, "service_find", options, criteria)
}
