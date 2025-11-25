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
