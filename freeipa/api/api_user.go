package api

type User struct {
	UID                   []string `json:"uid"`       // Login
	GivenName             []string `json:"givenname"` // First name
	SN                    []string `json:"sn"`        // Last name
	KrbPasswordExpiration []struct {
		DateTime IPATime `json:"__datetime__"`
	} `json:"krbpasswordexpiration"` // Password expiration
	Mail []string `json:"mail"` // Email
}

func (c *APIClient) UserAdd(uid string, givenname string, sn string, options JSON) (*User, error) {
	if options == nil {
		options = JSON{}
	}
	options["givenname"] = givenname
	options["sn"] = sn

	return apiRequest[User, string](c, "user_add", options, uid)
}

func (c *APIClient) UserDel(uid string, options JSON) (*JSON, error) {
	return apiRequest[JSON, []string](c, "user_del", options, uid)
}

func (c *APIClient) UserMod(uid string, options JSON) (*User, error) {
	return apiRequest[User, string](c, "user_mod", options, uid)
}

func (c *APIClient) UserDisable(uid string) (*bool, error) {
	return apiRequest[bool, string](c, "user_disable", nil, uid)
}

func (c *APIClient) UserEnable(uid string) (*bool, error) {
	return apiRequest[bool, string](c, "user_enable", nil, uid)
}

func (c *APIClient) UserShow(uid string, options JSON) (*User, error) {
	return apiRequest[User, string](c, "user_show", options, uid)
}

func (c *APIClient) UserUnlock(uid string) (*bool, error) {
	return apiRequest[bool, string](c, "user_unlock", nil, uid)
}

func (c *APIClient) UserFind(criteria string, options JSON) (*[]User, error) {
	return apiRequest[[]User, string](c, "user_find", options, criteria)
}
