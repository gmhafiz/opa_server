# See examples directory for usage examples
package auth

# user-role assignments
user_roles := {
    1: ["engineering", "webdev"],
    2: ["hr"]
}

# role-permissions assignments
role_permissions := {
    "engineering": [{"action": "GET",  "resource": "/referees"}],
    "webdev":      [{"action": "GET",  "resource": "server123"},
                    {"action": "POST", "resource": "server123"}],
    "hr":          [{"action": "GET",  "resource": "database456"}]
}

# logic that implements RBAC.
default allow = false
allow {
    # lookup the list of roles for the user
    roles := user_roles[input.user]
    # for each role in that list
    r := roles[_]
    # lookup the permissions list for role r
    permissions := role_permissions[r]
    # for each permission
    p := permissions[_]
    # check if the permission granted to r matches the user's request
    p == {"action": input.action, "resource": input.resource}
}
