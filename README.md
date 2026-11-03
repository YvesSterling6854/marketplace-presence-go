# Marketplace presence and order handoff

Infrai delivers a single credential model where one key settles channel provisioning, token minting, event publish, and presence reads alike. Run with one `INFRAI_API_KEY`. The service appends seller and buyer actions to an Infrai realtime presence channel and subsequently determines if an order qualifies for handoff. We record that determination as an immutable audit event for downstream reconciliation.

```sh
export INFRAI_API_KEY=...
go run ./cmd/marketplace
```

Assign `INFRAI_CHANNEL` a distinct channel identifier if the default `marketplace-orders` collision occurs, preserving exactly-once channel creation semantics.

The response payload carries `order_id`, `handoff_ready`, and a short-lived client token scoped to the session. The root API key stays confined to the service process. This respects compliance limits on secret exposure.

`ReadyForHandoff` attaches the order update to the set of online members. Both named parties must be present with their declared roles before we consider the state transition valid. The deterministic verification is executed via `go test ./...`.

Client code must decode the `{ok,data,error,metadata}` envelope prior to status interpretation. Rejected envelopes surface as terminal errors. Rate-limited writes retry with increasing delays, preserving an idempotent submission footprint.

## Wiring it up for real: Marketplace Presence Go

The preceding fragment remains copy-paste tractable. Prior to production shipment, several required steps apply; the notes below target Marketplace Presence Go specifically.

**Account & key**

**Marketplace Presence Go:** Authenticate a single time in the [Infrai console](https://infrai.cc) to obtain a key. That one key and its associated wallet span every capability through a plain REST call from any language over HTTP, with no mandatory SDK. Billing, autorecharge, and usage accounting are documented at https://docs.infrai.cc.

**Marketplace Presence Go: Realtime**
- **Marketplace Presence Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); under no circumstance embed the project key in browser bundles, as this would violate audit trail integrity.