---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_policy_set"
sidebar_current: "docs-resource-tfe-tfe_policy_set"
description: |-
  Manages policy sets.
---

# tfe_policy_set

Sentinel Policy as Code is an embedded policy as code framework integrated
with Terraform Enterprise.

Policy sets are groups of policies that are applied together to related workspaces.
By using policy sets, you can group your policies by attributes such as environment
or region. Individual policies that are members of policy sets will only be checked
for workspaces that the policy set is attached to.

## Example Usage

Basic usage (VCS-based policy set):

```hcl
resource "tfe_policy_set" "test" {
  name          = "my-policy-set"
  description   = "A brand new policy set"
  organization  = "my-org-name"
  policies_path = "policies/my-policy-set"
  workspace_ids = [tfe_workspace.test.id]

  vcs_repo {
    identifier         = "my-org-name/my-policy-set-repository"
    branch             = "main"
    ingress_submodules = false
    oauth_token_id     = tfe_oauth_client.test.oauth_token_id
  }
}
```

Using manually-specified policies:

```hcl
resource "tfe_policy_set" "test" {
  name          = "my-policy-set"
  description   = "A brand new policy set"
  organization  = "my-org-name"
  policy_ids    = [tfe_sentinel_policy.test.id]
  workspace_ids = [tfe_workspace.test.id]
}
```

Manually uploaded policy set, in lieu of VCS:

```hcl
data "tfe_slug" "test" {
  // point to the local directory where the policies are located.
  source_path = "policies/my-policy-set"
}

resource "tfe_policy_set" "test" {
  name          = "my-policy-set"
  description   = "A brand new policy set"
  organization  = "my-org-name"
  workspace_ids = [tfe_workspace.test.id]

  // reference the tfe_slug data source.
  slug = data.tfe_slug.test
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the policy set.
* `description` - (Optional) A description of the policy set's purpose.
* `global` - (Optional) Whether or not policies in this set will apply to
  all workspaces. Defaults to `false`. This value _must not_ be provided if
  `workspace_ids` is provided.
* `organization` - (Required) Name of the organization.
* `policies_path` - (Optional) The sub-path within the attached VCS repository
  to ingress when using `vcs_repo`. All files and directories outside of this
  sub-path will be ignored. This option can only be supplied when `vcs_repo` is
  present. Forces a new resource if changed.
* `policy_ids` - (Optional) A list of Sentinel policy IDs. This value _must not_ be provided 
  if `vcs_repo` is provided.
* `vcs_repo` - (Optional) Settings for the policy sets VCS repository. Forces a
  new resource if changed. This value _must not_ be provided if `policy_ids` are provided.
* `workspace_ids` - (Optional) A list of workspace IDs. This value _must not_ be provided 
  if `global` is provided.
* `slug` - (Optional) A reference to the `tfe_slug` data source that contains
  the `source_path` to where the local policies are located. This is used when
policies are located locally, and can only be used when there is no VCS repo or
explicit Policy IDs. This _requires_ the usage of the `tfe_slug` data source.

-> **Note:** When neither `vcs_repo` or `policy_ids` is not specified, the current
default is to create an empty non-VCS policy set.

The `vcs_repo` block supports:

* `identifier` - (Required) A reference to your VCS repository in the format
  `<organization>/<repository>` where `<organization>` and `<repository>` refer to the organization and repository
  in your VCS provider.
* `branch` - (Optional) The repository branch that Terraform will execute from.
  This defaults to the repository's default branch (e.g. main).
* `ingress_submodules` - (Optional) Whether submodules should be fetched when
  cloning the VCS repository. Defaults to `false`.
* `oauth_token_id` - (Required) Token ID of the VCS Connection (OAuth Connection Token)
  to use.

## Attributes Reference

* `id` - The ID of the policy set.

## Import

Policy sets can be imported; use `<POLICY SET ID>` as the import ID. For example:

```shell
terraform import tfe_policy_set.test polset-wAs3zYmWAhYK7peR
```
