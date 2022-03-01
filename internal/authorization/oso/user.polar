actor User {}

resource User {
  permissions = ["read", "create", "modify", "delete"];
}

has_permission(user: User, "read", userResource: User) if
  user.UserID = userResource.UserID;

allow(actor, action, resource) if
  has_permission(actor, action, resource);
