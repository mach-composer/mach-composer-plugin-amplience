# Amplience Plugin for Mach Composer

This repository contains the Amplience plugin for Mach Composer. It requires
Mach Composer 3.x

It uses the Terraform provider for Amplience,
see https://github.com/labd/terraform-provider-amplience/

## Options

### `global.amplience`

- `client_id` (optional): The client ID for the Amplience API. This is used as
  the default client ID for all sites. Must be set if site is left empty
- `client_secret` (optional): The client secret for the Amplience API. This is
  used as the default client secret for all. Must be set if site is left empty

### `sites.[*].amplience`

- `hub_id` (optional): The hub ID for the Amplience API. Either this or `hubs` is required.
- `client_id` (optional): The client ID for the Amplience API. This is used as
  the default client ID for all sites. Must be set if global is left empty
- `client_secret` (optional): The client secret for the Amplience API. This is
  used as the default client secret for all. Must be set if global is left empty
- `hubs` (optional): A list of extra hubs to use. Either this or `hub_id` is required.
  Each hub has the following properties:
  - `name` (required): The name of the hub
  - `hub_id` (required): The hub ID for the Amplience API.
  - `client_id` (required): The client ID for the extra Amplience API.
  - `client_secret` (required): The client secret for the extra Amplience API.


## Examples

### Single hub

```yaml
mach_composer:
  plugins:
    amplience:
      source: mach-composer/amplience
      version: 0.1.3

global:
  # ...
  amplience:
    client_id: your-client-id
    client_secret: your-client-secret

sites:
  - identifier: my-site
    # ...
    amplience:
      hub_id: "hub-default"
```

#### Component usage

```hcl
terraform {
  required_providers {
    amplience = {
      source  = "labd/amplience"
      version = "~> 0.3.7"
    }
  }
}

variable "amplience_client_id" {}
variable "amplience_client_secret" {}
variable "amplience_hub_id" {}

resource "amplience_content_repository" "my-content-repository" {
  label = "my-label-primary"
  name  = "my-name-primary-1"
}
```

### Multiple hubs

```yaml
mach_composer:
  plugins:
    amplience:
      source: mach-composer/amplience
      version: 0.1.3

# ...

sites:
  - identifier: my-site
    # ...
    amplience:
      hubs:
        - name: hub_1 
          client_id: "id-default"
          client_secret: "secret-default"
          hub_id: "hub-default"
        - name: hub_2
          client_id: "id-alternate"
          client_secret: "secret-alternate"
          hub_id: "hub-alternate"
```

#### Component usage

```hcl
terraform {
  required_providers {
    amplience = {
      source  = "labd/amplience"
      version = "~> 0.3.7"
      configuration_aliases = [ amplience.hub_1, amplience.hub_2 ]
    }
  }
}

variable "amplience_hub_1_client_id" {}
variable "amplience_hub_1_client_secret" {}
variable "amplience_hub_1_hub_id" {}
variable "amplience_hub_2_client_id" {}
variable "amplience_hub_2_client_secret" {}
variable "amplience_hub_2_hub_id" {}

resource "amplience_content_repository" "my-content-repository_1" {
  provider = amplience.hub_1
  label = "my-label-primary"
  name  = "my-name-primary-1"
}

resource "amplience_content_repository" "my-content-repository_2" {
  provider = amplience.hub_2
  label = "my-label-alternate"
  name  = "my-name-alternate-1"
}
```
