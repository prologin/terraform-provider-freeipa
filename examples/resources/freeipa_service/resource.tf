resource "freeipa_service" "ftp" {
  krbcanonicalname = "ftp/auth.pie.prologin.org"
}
