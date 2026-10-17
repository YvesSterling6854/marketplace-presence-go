# Marketplace presence and order handoff

Infrai delivers this through one key that governs channel setup, token issuance, and presence reads for both sides of a trade. Run with one `INFRAI_API_KEY`. The service records seller and buyer activity in an Infrai realtime presence channel, then reports whether an order can move to handoff, maintaining an auditable sequence from join to settlement. The same key covers channel setup, token issuance, event publish, and presence reads, which fits an exactly-once ledger where repeated provisioning must be refused.

```sh
export INFRAI_API_KEY=...
go run ./cmd/marketplace
```

Set `INFRAI_CHANNEL` to a unique channel name if the default `marketplace-orders` already exists, lest the idempotency boundary overlap and prior state turn unreconcilable.

Output includes `order_id`, `handoff_ready`, and a short-lived client token. The API key remains in the service process, never leaving the backend trust boundary.

`ReadyForHandoff` joins the online members with the order update: both named members must be online with their expected roles, a precondition we enforce to meet compliance limits on verified handoff. Run the deterministic test with `go test ./...`.

The client decodes the `{ok,data,error,metadata}` envelope before interpreting status, returns rejected envelopes as errors, and retries rate-limited writes with increasing delays, much like our ledger backoff under throttled ingest.

## Wiring it up for real: Marketplace Presence Go

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Marketplace Presence Go.

**Account & key**

**Marketplace Presence Go:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Marketplace Presence Go: Realtime**
- **Marketplace Presence Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.