---
page_title: "Provider: azname"
description: |-
  Provider for generating standardized Azure resource names following naming conventions.
---

# azname Provider

Provider for generating standardized Azure resource names following naming conventions.

This provider helps maintain consistent resource naming across your Azure infrastructure by providing configurable templates
and functions for generating resource names. It supports both global resources and child resources, with configurable
separators, random suffixes, and instance numbers.

#### Why is this provider needed?

When working with Azure, maintaining consistent resource naming is critical. This provider makes it easy to align resource names to organizational standards with minimal input from developers. 

A key advantage of this provider is that generated resource names are **stored in Terraform state**. This means your resource names remain stable even when naming conventions or input parameters change. Since Azure treats resource names as immutable identifiers, simple changes to name generation logic can trigger destructive rebuilds that cascade across your infrastructure. By persisting names in state, this provider protects you from unintended resource recreation and the associated downtime and data loss.

#### Why not just a terraform module?

While a terraform module can encapsulate most of the required functionality, it was:
a. Not possible to store generated names in state without using resources.
b. Require all parameters to be passed each time a name is needed, leading to repetitive code. This provider allows you to configure commmon settings at the provider level, or even as environment variables - which can be useful in CI/CD pipelines.

#### What about azurecaf_name?

The [azurecaf](https://registry.terraform.io/providers/aztfmod/azurecaf/latest) provider actually inspired the creation of this provider! However I wasn't satisfied with the functionality and ease of use of the `azurecaf_name` data source. I found that it was not flexible enough for my needs and I wanted to have more control over the naming process. It also does not persist naming in state, which is a key advantage of this provider.

## Example Usage

```terraform
# Basic provider configuration with default settings
provider "azname" {
  # Separator character used between name components (default: "-")
  separator = "-"

  # Length of random suffix appended to names (1-6, default: 3)
  random_length = 4

  # Length of instance number padding (1-6, default: 3)
  # With instance_length = 3, instance 1 becomes "001"
  instance_length = 3

  # Remove special characters to ensure Azure naming compliance (default: true)
  clean_output = true

  # Trim names to fit within Azure resource length limits (default: true)
  trim_output = true
}

# Advanced provider configuration with prefixes and suffixes
# These will be applied to all generated names
provider "azname" {
  separator     = "-"
  random_length = 4

  # Global prefixes applied to all resource names
  prefixes = ["corp", "az"]

  # Global suffixes applied to all resource names
  suffixes = ["001"]

  # Custom templates for name generation (advanced usage)
  # Use ~ as placeholder for the separator character
  # Default template: {prefix}~{resource_type}~{workload}~{environment}~{service}~{location}{instance}{rand}~{suffix}
  template = "{prefix}~{resource_type}~{workload}~{environment}~{location}{rand}"

  # Template for child resources (default: {parent_name}~{resource_type}{instance}~{rand})
  template_child = "{parent_name}~{resource_type}~{workload}{rand}"
}
```

## Customizing Resource Slugs and Regions (Overrides)

The provider supports customization of resource abbreviations (slugs), region short names, and even adding completely new resource types or regions that aren't built into the provider. This is done via an `azname_overrides.yaml` file.

### Auto-Discovery

The provider will automatically search for and load an overrides file from:

- `./azname_overrides.yaml` (current working directory)

If an overrides file is found, it will be loaded automatically. If the file contains errors, a warning will be displayed but the provider will continue to work with default values.

### Override File Structure

The overrides file supports four main sections:

#### 1. Resource Slug Overrides

Override the default CAF abbreviations for existing Azure resource types:

```yaml
resource_slug_overrides:
  azurerm_resource_group: "resourcegroup"    # Change "rg" to "resourcegroup"
  azurerm_storage_account: "storage"         # Change "st" to "storage"
  azurerm_key_vault: "vault"                 # Change "kv" to "vault"
```

#### 2. Region Shortname Overrides

Override the short names used for existing Azure regions (matched by CLI name):

```yaml
region_shortname_overrides:
  eastus: "use"           # Change "eus" to "use"
  westus2: "usw2"         # Change "wus2" to "usw2"
  australiaeast: "aue"    # Change "ae" to "aue"
```

#### 3. New Resources

Define custom resource types not yet supported by the provider:

```yaml
new_resources:
  azurerm_custom_resource:
    slug: "custom"              # Resource abbreviation
    min_length: 1               # Minimum name length
    max_length: 63              # Maximum name length
    scope: "resourceGroup"      # "global", "resourceGroup", or "parent"
    dashes: true                # Whether dashes are allowed
    lowercase: true             # Whether name should be lowercase
```

**Note:** Custom resources do not have regex validation applied, giving you full flexibility in naming.

#### 4. New Regions

Define custom regions not yet supported by the provider:

```yaml
new_regions:
  customregion:
    cli_name: "customregion"    # CLI name (used in lookups)
    full_name: "Custom Region"  # Full display name
    short_name: "cr"            # Short name for name generation
```

### Complete Example

See the full example override file: [examples/azname_overrides.yaml](https://github.com/BHoggs/terraform-provider-azname/blob/main/examples/azname_overrides.yaml)

### Usage Notes

- **Thread-safe**: Overrides are applied once during provider initialization using `sync.Once`
- **No provider code changes needed**: Overrides are merged into the base datasets, so all existing functionality works seamlessly
- **Graceful degradation**: If the overrides file has errors, a warning is displayed and the provider continues with defaults
- **Precedence**: Override values take precedence over built-in defaults
- **Validation**: The file is validated on load to catch configuration errors early

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `clean_output` (Boolean) Remove special characters from generated names to ensure compatibility with Azure naming rules. Can be set via `AZNAME_CLEAN_OUTPUT` environment variable (1 for true, 0 for false).
- `environment` (String) Default environment name (e.g., dev, test, prod) to use in resource names. Can be overridden at resource/data source level. Can be set via `AZNAME_ENVIRONMENT` environment variable.
- `instance_length` (Number) Length of instance number padding in generated names. Must be between 1 and 6. Can be set via `AZNAME_INSTANCE_LENGTH` environment variable.
- `location` (String) Default location (e.g., eastus, westeurope) to use in resource names. Can be overridden at resource/data source level. Can be set via `AZNAME_LOCATION` environment variable.
- `prefixes` (List of String) List of prefixes to prepend to resource names. These will be joined using the separator character. Can be set via `AZNAME_PREFIX` environment variable (comma-separated).
- `random_length` (Number) Length of random suffix to append to generated names. Must be between 1 and 6. Can be set via `AZNAME_RANDOM_LENGTH` environment variable.
- `separator` (String) Character to use as separator in resource names. Must be a single character. Can be set via `AZNAME_SEPARATOR` environment variable.
- `suffixes` (List of String) List of suffixes to append to resource names. These will be joined using the separator character. Can be set via `AZNAME_SUFFIX` environment variable (comma-separated).
- `template` (String) Global template for resource name generation. Uses ~ as a placeholder for the separator character. Can be set via `AZNAME_TEMPLATE` environment variable.
- `template_child` (String) Template for child resource name generation. Uses ~ as a placeholder for the separator character. Can be set via `AZNAME_TEMPLATE_CHILD` environment variable.
- `trim_output` (Boolean) Trim generated names to fit Azure resource length limits while preserving important parts. Can be set via `AZNAME_TRIM_OUTPUT` environment variable (1 for true, 0 for false).
