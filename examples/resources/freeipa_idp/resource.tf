resource "freeipa_idp" "goauthentik-test-client" {
  cn               = "goauthentik-test-client"
  clientid         = "b004cf15eb3cd7101b923a2a4f93ad71f946f9d5d064b1beb63f4b9ddfca2acb"
  clientsecret     = "70f98bbb70845e0e0a1948e2f64229605e20e5482c69bb427d289837bc106c55"
  authendpoint     = "https://auth.example.com/application/o/authorize/"
  devauthendpoint  = "https://auth.example.com/application/o/device/"
  tokenendpoint    = "https://auth.example.com/application/o/token/"
  userinfoendpoint = "https://auth.example.com/application/o/userinfo/"
  keysendpoint     = "https://auth.example.com/application/o/goauthentik-test-client/jwks/"
  scope            = "openid profile email"
  sub              = "email"
}
