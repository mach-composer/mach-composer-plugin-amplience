# Amplience Plugin for Mach Composer 

This repository contains the Amplience plugin for Mach Composer. It requires Mach Composer 3.x

It uses the Terraform provider for Amplience, see https://github.com/labd/terraform-provider-amplience/

## Usage

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
      hub_id: your-hub-id
      client_id: your-client-id # override
      client_secret: your-client-secret # override
```
