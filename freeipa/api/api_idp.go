package api

type IdentityProvider struct {
	CN []string `json:"cn"` // IDP name

	ClientID     []string `json:"ipaidpclientid"` // Client ID
	ClientSecret []struct {
		Secret Base64EncodedSecret `json:"__base64__"`
	} `json:"ipaidpclientsecret"` // Client secret

	Scope []string `json:"ipaidpscope"` // Scope
	Sub   []string `json:"ipaidpsub"`   // External IdP user identifier attribute

	AuthEndpoint     []string `json:"ipaidpauthendpoint"`     // Authentication endpoint
	DevAuthEndpoint  []string `json:"ipaidpdevauthendpoint"`  // Device authentication endpoint
	TokenEndpoint    []string `json:"ipaidptokenendpoint"`    // Token endpoint
	UserInfoEndpoint []string `json:"ipaidpuserinfoendpoint"` // User info endpoint
	KeysEndpoint     []string `json:"ipaidpkeysendpoint"`     // JWKS endpoint
	IssuerURL        []string `json:"ipaidpissuerurl"`        // OIDC URL
}

func (c *APIClient) IdentityProviderAddGeneric(cn string, clientID string, clientSecret string, authEndpoint string, devAuthEndpoint string, tokenEndpoint string, userInfoEndpoint string, keysEndpoint string, options JSON) (*IdentityProvider, error) {
	if options == nil {
		options = JSON{}
	}
	options["ipaidpclientid"] = clientID
	options["ipaidpclientsecret"] = clientSecret
	options["ipaidpauthendpoint"] = authEndpoint
	options["ipaidpdevauthendpoint"] = devAuthEndpoint
	options["ipaidptokenendpoint"] = tokenEndpoint
	options["ipaidpuserinfoendpoint"] = userInfoEndpoint
	options["ipaidpkeysendpoint"] = keysEndpoint

	return c.IdentityProviderAdd(cn, options)
}

// Usage of IdentityProviderAddGeneric is recommended over this function
func (c *APIClient) IdentityProviderAdd(cn string, options JSON) (*IdentityProvider, error) {
	return apiRequest[IdentityProvider, string](c, "idp_add", options, cn)
}

func (c *APIClient) IdentityProviderDel(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, []string](c, "idp_del", options, cn)
}

func (c *APIClient) IdentityProviderMod(cn string, options JSON) (*IdentityProvider, error) {
	return apiRequest[IdentityProvider, string](c, "idp_mod", options, cn)
}

func (c *APIClient) IdentityProviderShow(cn string, options JSON) (*IdentityProvider, error) {
	return apiRequest[IdentityProvider, string](c, "idp_show", options, cn)
}

func (c *APIClient) IdentityProviderFind(criteria string, options JSON) (*[]IdentityProvider, error) {
	return apiRequest[[]IdentityProvider, string](c, "idp_find", options, criteria)
}
