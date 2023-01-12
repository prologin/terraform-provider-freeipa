resource "freeipa_group" "important_users" {
  cn          = "important_users"
  description = "A group of important users"
}
