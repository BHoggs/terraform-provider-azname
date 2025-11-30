# Terraform Provider Azname - AI Coding Guide

## Project Overview
This is a Terraform provider built with the [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) that generates standardized Azure resource names following Azure Cloud Adoption Framework (CAF) conventions. The provider includes a resource (`azname_name`), data source (`azname_name`), and three provider functions for region name conversion.

**Key Architecture Principle**: Names are persisted in Terraform state to prevent destructive rebuilds when naming logic changes, since Azure treats resource names as immutable identifiers.

## Critical Developer Workflows

### Building & Testing
```powershell
# Build and install locally
make install  # or: go install

# Run unit tests (fast, no TF_ACC required)
make test     # or: go test -v -cover -timeout=120s -parallel=10 ./...

# Run acceptance tests (creates real resources)
make testacc  # or: TF_ACC=1 go test -v -cover -timeout 120m ./...
```

### Code Generation (Required After Changes)
```powershell
# Generate resource definitions from resourceDefinition.json
make genresources  # or: cd internal/resources; go generate ./...

# Generate provider documentation
make gendocs       # or: terraform fmt -recursive examples/; go tool tfplugindocs generate --provider-name azname

# Generate all
make genall
```

**IMPORTANT**: 
- Run `make genresources` after editing `internal/resources/resourceDefinition.json`. This regenerates `models_generated.go` via `gen.go` using Go templates.
- **NEVER edit markdown files in `docs/` directly** - they are generated from `templates/` and `examples/`. Update those source files instead, then run `make gendocs`.
- Always check if documentation updates are needed after code changes. Documentation is generated from:
  - Template files in `templates/` directory
  - Example configurations in `examples/` directory
  - Code comments and schema descriptions

## Project-Specific Patterns

### Resource Definition System
- **Source of Truth**: `internal/resources/resourceDefinition.json` contains all Azure resource type metadata (5700+ lines)
- **Code Generation**: `internal/resources/gen.go` (go:build ignore) reads JSON and generates `models_generated.go` via `templates/model.tmpl`
- **Key Fields in JSON**:
  - `scope`: "global" (requires random suffix), "resourceGroup", or "parent" (child resources)
  - `dashes`: whether separator character is allowed
  - `slug`: CAF prefix abbreviation (e.g., "rg", "st", "kv")
  - `validation_regex`: final name must match this
  - `regex`: characters to strip when `clean_output=true`

### Name Generation Logic (`internal/provider/name_generator.go`)
1. Template expansion with `{prefix}~{resource_type}~{workload}~...` placeholders
2. Separator replacement (`~` → configured separator, default `-`)
3. Character cleaning via regex (if `clean_output=true`)
4. Length trimming (if `trim_output=true`)
5. Validation against `validation_regex`

**Random Suffix Behavior**:
- Global-scope resources need random suffixes for uniqueness
- Without `random_seed`, names show "(known after apply)" in plans (see `NeedsRandomGeneration()`)
- With `random_seed`, deterministic random values appear in plans using `rand.NewPCG(seed, seed)`

### Provider Configuration Pattern
Provider accepts config via:
1. HCL provider block attributes
2. Environment variables (`AZNAME_*` prefix)
3. Hardcoded defaults in `provider.go`

Priority: HCL > Env Vars > Defaults. See `Configure()` method in `provider.go`.

### Region Handling (`internal/regions/regions.go`)
- Maintains list of Azure regions with CLI name, full name, short name, and paired region
- `GetRegionByAnyName()` accepts any format and returns standardized struct
- Provider functions (`region_cli_name`, `region_full_name`, `region_short_name`) expose this

## Testing Conventions

### Test Structure
- Unit tests: `resource.UnitTest()` with `PreCheck`, no `TF_ACC` needed
- Acceptance tests: `resource.Test()` with `TF_ACC=1`, creates real resources
- Test provider config in `provider_test.go`: `providerConfig` constant, `testAccProtoV6ProviderFactories`

### Example Test Pattern
```go
resource.UnitTest(t, resource.TestCase{
    PreCheck:                 func() { testAccPreCheck(t) },
    ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
        {
            Config: providerConfig + `
                data "azname_name" "test" {
                    name = "myapp"
                    resource_type = "azurerm_resource_group"
                }
            `,
            Check: resource.ComposeAggregateTestCheckFunc(
                resource.TestCheckResourceAttr("data.azname_name.test", "result", "expected-name"),
            ),
        },
    },
})
```

## Key Files & Directories

- `internal/provider/name_generator.go`: Core name generation logic
- `internal/provider/name_resource.go`: Resource with state persistence + `ModifyPlan()` for plan-time name generation
- `internal/resources/resourceDefinition.json`: Azure resource metadata (edit this, then run `make genresources`)
- `internal/resources/gen.go`: Code generator (go:build ignore)
- `internal/regions/regions.go`: Azure region conversion utilities
- `examples/`: Used for documentation generation by tfplugindocs

## Framework-Specific Details

### Plugin Framework (Not SDK)
This provider uses `terraform-plugin-framework` (v1.16+), NOT the older `terraform-plugin-sdk`. Key differences:
- Use `types.String`, `types.Int64`, etc., not primitives
- Check `IsNull()` and `IsUnknown()` for optional/computed values
- Implement `ModifyPlan()` for plan-time computation (see `name_resource.go`)
- Provider functions use `function.Function` interface (see `region_function.go`)

### State Persistence Pattern
The `azname_name` resource preserves generated names in state via `ModifyPlan()`:
```go
// If result exists in state and is known, preserve it
if !state.Result.IsNull() && !state.Result.IsUnknown() {
    plan.Result = state.Result
    return
}
```
This prevents name regeneration on subsequent applies, protecting against destructive rebuilds.

## Common Pitfalls

1. **Forgetting Code Generation**: Always run `make genresources` after JSON edits
2. **Test Environment Variable**: Set `TF_ACC=1` for acceptance tests, omit for unit tests
3. **Random Seed in Tests**: Use `random_seed` attribute in tests to get deterministic results
4. **Provider vs Resource Config**: Resource-level attributes override provider-level (environment, location, separator, etc.)
5. **Editing Generated Documentation**: Never edit `docs/` markdown files directly - they're generated by tfplugindocs. Update `templates/` or `examples/` instead and run `make gendocs`
