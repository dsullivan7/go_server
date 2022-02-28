actor User {}

resource User {
  permissions = ["read", "create", "modify", "delete"];
}

has_permission(actor: User, "read", user: UserResource) if
  User.UserID = UserResource.UserID;

allow(actor, action, resource) if
  has_permission(actor, action, resource);
