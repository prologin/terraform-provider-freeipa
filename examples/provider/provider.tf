provider "freeipa" {
  server   = "https://freeipa.example.com" # optionally use FREEIPA_SERVER env var
  user     = "admin"                       # optionally use FREEIPA_USER env var
  password = "password"                    # optionally use FREEIPA_PASSWORD env var
}
