# Amplience Plugin for Mach Composer 

This repository contains the Sentry plugin for Mach Composer. It requires Mach Composer 3.x

It uses the Terraform provider for Amplience, see https://github.com/labd/terraform-provider-amplience/

## Usage

```yaml
mach_composer:
  plugins:
    amplience:
      version: latest
      
global:
  # ...
  
sites:
  - identifier: my-site
    # ...
    amplience:
      hub_id: hub-id
      client_id: client-id
      client_secret: client-secret
```
