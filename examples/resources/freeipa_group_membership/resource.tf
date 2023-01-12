resource "freeipa_group_membership" "managers" {
  group  = "ipausers"
  member = "root"
  type   = "user"

  manager = true
}
