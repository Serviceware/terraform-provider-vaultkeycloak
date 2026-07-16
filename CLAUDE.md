# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Terraform provider (built on `terraform-plugin-sdk/v2`) that configures the [Serviceware vault-plugin-secrets-keycloak](https://github.com/Serviceware/vault-plugin-secrets-keycloak) Vault plugin. It does not talk to Keycloak directly at runtime — it writes configuration into a Vault server via the Vault API, and the Vault plugin talks to Keycloak.

## Commands

```sh
make build                  # go build -o terraform-provider-vaultkeycloak
make install                # build + copy into ~/.terraform.d/plugins/... (local terraform use)
make test                   # unit tests (30s timeout, parallel=4)
make testacc                # acceptance tests (TF_ACC=1, 120m timeout)

# Single acceptance test (requires Docker; testcontainers spins up Vault + Keycloak):
TF_ACC=1 go test ./vaultkeycloak/ -v -run TestAccVaultKeycloakBasic26_1_4 -timeout 30m
```

Acceptance tests (`TestAcc*`) manage their own infrastructure via testcontainers-go compose using `testing/docker-compose.yaml`: a Keycloak container (version parameterized via `KEYCLOAK_VERSION` env) and a custom Vault image (`testing/vault/Dockerfile`) with the keycloak secrets plugin pre-installed, dev root token `root`. The tests set `VAULT_ADDR`/`VAULT_TOKEN` themselves and bootstrap Keycloak admin clients with gocloak before running Terraform test steps. CI (`.github/workflows/test.yaml`) runs them across a Terraform version matrix.

For manual end-to-end testing, `demo/` plus `docker-compose -f testing/docker-compose.yaml up` gives a live Vault+Keycloak; see README for the flow and for using Terraform dev_overrides to test a locally built provider against an existing project.

## Architecture

- `main.go` — standard SDKv2 `plugin.Serve` entrypoint.
- `vaultkeycloak/provider.go` — provider schema (`vault_address`, `vault_token`, `add_address_to_env`). `providerConfigure` returns a `*vault/api.Client` as the provider meta; token resolution (`GetToken`) mirrors hashicorp/terraform-provider-vault: explicit attribute → `VAULT_TOKEN` env → Vault CLI token helper (`~/.vault-token`).
- Two resources, both thin CRUD wrappers around `client.Logical().Write/Read/Delete` against the mounted plugin:
  - `vaultkeycloak_secret_backend` → writes `{path}/config/connection`; resource ID is the mount `path`.
  - `vaultkeycloak_secret_backend_per_realm_config` → writes `{path}/config/realms/{realm}/connection`; resource ID encodes both path and realm (see `pathAndRealmFromId`).
- The Vault mount itself is not created by the provider — tests mount `vault-plugin-secrets-keycloak` explicitly, and users do `vault secrets enable` (or use the vault provider).

Docs in `docs/` and examples in `examples/` follow Terraform Registry layout; update them when changing resource schemas.
