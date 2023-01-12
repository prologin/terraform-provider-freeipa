resource "freeipa_user" "john_doe" {
  uid       = "john.doe"
  givenname = "John"
  sn        = "Doe"
  password  = "ThisPasswordIsDefinitelyExpiredAlready"

  # Defaults to time of resource creation (default from FreeIPA API), meaning
  # that if not set, the password is immediately expired, requiring the user to
  # change it on first login
  krbpasswordexpiration = "1970-01-01T00:00:00Z"

  mail = [
    "john.doe@example.com",
    "john@example.com"
  ]
}
