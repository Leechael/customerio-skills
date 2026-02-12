# Customer.io API Reference

Internal reference for Customer.io integration development.

## API Overview

Customer.io exposes **three independent APIs**, each with its own auth and base URL:

| API | Purpose | Auth | Base URL (US) | Base URL (EU) |
|-----|---------|------|---------------|---------------|
| **Track API** | Write: identify, events, delete, suppress, devices, segments | Basic Auth (`site_id:api_key`) | `https://track.customer.io` | `https://track-eu.customer.io` |
| **App API** | Read + transactional send: customers, segments, campaigns, messages, exports | Bearer token | `https://api.customer.io/v1` | `https://api-eu.customer.io/v1` |
| **Pipelines API** | CDP-style ingest (POST-only) | Bearer token | `https://cdp.customer.io/v1` | `https://cdp-eu.customer.io/v1` |

## OpenAPI Specs (machine-readable)

Downloaded from `docs.customer.io` HTML source. These are the canonical specs used for client generation.

| Spec | URL | Size | Endpoints |
|------|-----|------|-----------|
| App API (Journeys) | https://docs.customer.io/files/journeys-app.json | ~2 MB | 115 |
| Track API (Journeys) | https://docs.customer.io/files/journeys-track.json | ~835 KB | 18 |
| Pipelines API | https://docs.customer.io/files/pipelines.json | ~1.7 MB | 7 |

> **Note:** The CIO docs site is JavaScript-rendered (ReDoc SPA). WebFetch / curl on doc pages returns JS code, not content. Always use the raw JSON spec files above.

### Alternative spec references found in HTML source

```
cio_spec.json
cio_journeys_app_api.json
cio_journeys_track_api.json
/files/journeys-app.json    ← this is the one that works
/files/journeys-track.json  ← this is the one that works
/files/pipelines.json       ← this is the one that works
```

## Postman Collections

- **Public workspace:** https://www.postman.com/customerio
- **Journeys Track API:** https://www.postman.com/customerio/customer-io-s-public-workspace/collection/u78pcoo/journeys-track-api
- **Journeys App API:** https://www.postman.com/customerio/customer-io-s-public-workspace/collection/gv7unvo/journeys-app-api

> Postman pages are also JS-rendered. Use the OpenAPI specs above for programmatic access.

## Official Documentation

- **Track API docs:** https://docs.customer.io/integrations/api/track/
- **App API docs:** https://docs.customer.io/integrations/api/app/
- **API overview:** https://docs.customer.io/integrations/api/customerio-apis/

## Official SDKs (for reference)

- **Node.js:** https://github.com/customerio/customerio-node
- **Go:** https://github.com/customerio/go-customerio
- **Python (third-party):** https://github.com/customerio/customerio-python (basic, Track API only)

## Rate Limits

| Category | Limit |
|----------|-------|
| App API (general) | 10 req/s |
| App API (transactional send) | 100 req/s |
| Broadcast trigger | 1 req/10s |
| Track API | ~30 req/s (soft) |

## Authentication

### Track API — Basic Auth

```
Authorization: Basic base64(site_id:api_key)
```

Env vars: `CUSTOMERIO_SITE_ID`, `CUSTOMERIO_API_KEY`

### App API — Bearer Token

```
Authorization: Bearer {app_api_key}
```

Env var: `CUSTOMERIO_APP_API_KEY`

Obtain from: CIO dashboard → Settings → API Credentials → App API Key

## Track API Endpoints (18 total)

### V2 Endpoints (primary)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v2/entity` | Single operation (identify, event, delete, suppress, unsuppress) for person or object |
| POST | `/api/v2/batch` | Batch multiple operations in one request |

V2 uses a unified entity model with `type` (person/object) and `action` (identify/event/delete/suppress/unsuppress).

### V1 Endpoints (supplementary)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/events` | Track anonymous event |
| POST | `/api/v1/merge_customers` | Merge two customer profiles |
| PUT | `/api/v1/customers/{id}/devices` | Register device for push |
| DELETE | `/api/v1/customers/{id}/devices/{device_id}` | Remove device |
| POST | `/api/v1/segments/{id}/add_customers` | Add to manual segment |
| POST | `/api/v1/segments/{id}/remove_customers` | Remove from manual segment |

## App API Endpoints (115 total, 21 categories)

### Customers (10 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/customers/{id}/attributes` | Look up person attributes |
| GET | `/v1/customers?email={email}` | Find customers by email |
| POST | `/v1/customers` | Search customers (filter syntax) |
| POST | `/v1/customers/attributes` | List customers with attributes |
| GET | `/v1/customers/{id}/activities` | Customer activity log |
| GET | `/v1/customers/{id}/messages` | Messages sent to customer |
| GET | `/v1/customers/{id}/segments` | Segments customer belongs to |
| GET | `/v1/customers/{id}/relationships` | Customer-object relationships |
| GET | `/v1/customers/{id}/subscription_preferences` | Subscription preferences |

### Objects (4 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/object_types` | List object types |
| POST | `/v1/objects` | Search objects |
| GET | `/v1/objects/{type_id}/{object_id}/attributes` | Object attributes |
| GET | `/v1/objects/{type_id}/{object_id}/relationships` | Object relationships |

### Segments (7 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/segments` | List all segments |
| POST | `/v1/segments` | Create manual segment |
| GET | `/v1/segments/{id}` | Get segment |
| DELETE | `/v1/segments/{id}` | Delete segment |
| GET | `/v1/segments/{id}/customer_count` | Segment member count |
| GET | `/v1/segments/{id}/membership` | Segment member list |
| GET | `/v1/segments/{id}/dependencies` | What uses this segment |

### Collections (7 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/collections` | List collections |
| POST | `/v1/collections` | Create collection (inline data or URL) |
| GET | `/v1/collections/{id}` | Get collection |
| PUT | `/v1/collections/{id}` | Update collection |
| DELETE | `/v1/collections/{id}` | Delete collection |
| GET | `/v1/collections/{id}/content` | Get collection data rows |
| PUT | `/v1/collections/{id}/content` | Replace collection data |

### Campaigns (13 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/campaigns` | List campaigns |
| GET | `/v1/campaigns/{id}` | Get campaign |
| GET | `/v1/campaigns/{id}/actions` | List campaign actions |
| GET | `/v1/campaigns/{id}/actions/{action_id}` | Get campaign action |
| PUT | `/v1/campaigns/{id}/actions/{action_id}` | Update campaign action |
| GET | `/v1/campaigns/{id}/metrics` | Campaign metrics |
| GET | `/v1/campaigns/{id}/metrics/links` | Campaign link metrics |
| GET | `/v1/campaigns/{id}/actions/{action_id}/metrics` | Action metrics |
| GET | `/v1/campaigns/{id}/journey_metrics` | Journey metrics |
| GET | `/v1/campaigns/{id}/messages` | Campaign messages |
| GET | `/v1/campaigns/{id}/actions/{action_id}/language/{language}` | Action translation |
| PUT | `/v1/campaigns/{id}/actions/{action_id}/language/{language}` | Update translation |
| GET | `/v1/campaigns/{id}/actions/{action_id}/metrics/links` | Action link metrics |

### Broadcasts (16 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/broadcasts` | List broadcasts |
| GET | `/v1/broadcasts/{id}` | Get broadcast |
| POST | `/v1/campaigns/{id}/triggers` | Trigger broadcast (**rate: 1/10s**) |
| GET | `/v1/broadcasts/{id}/triggers` | List triggers |
| GET | `/v1/broadcasts/{id}/triggers/{trigger_id}` | Trigger status |
| GET | `/v1/broadcasts/{id}/triggers/{trigger_id}/errors` | Trigger errors |
| GET | `/v1/broadcasts/{id}/actions` | List broadcast actions |
| GET | `/v1/broadcasts/{id}/actions/{action_id}` | Get broadcast action |
| PUT | `/v1/broadcasts/{id}/actions/{action_id}` | Update broadcast action |
| GET | `/v1/broadcasts/{id}/metrics` | Broadcast metrics |
| GET | `/v1/broadcasts/{id}/metrics/links` | Broadcast link metrics |
| GET | `/v1/broadcasts/{id}/actions/{action_id}/metrics` | Action metrics |
| GET | `/v1/broadcasts/{id}/actions/{action_id}/metrics/links` | Action link metrics |
| GET | `/v1/broadcasts/{id}/messages` | Broadcast messages |
| GET | `/v1/broadcasts/{id}/actions/{action_id}/language/{language}` | Action translation |
| PUT | `/v1/broadcasts/{id}/actions/{action_id}/language/{language}` | Update translation |

### Newsletters (16 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/newsletters` | List newsletters |
| GET | `/v1/newsletters/{id}` | Get newsletter |
| DELETE | `/v1/newsletters/{id}` | Delete newsletter |
| GET | `/v1/newsletters/{id}/contents` | List variants |
| GET | `/v1/newsletters/{id}/contents/{content_id}` | Get variant |
| PUT | `/v1/newsletters/{id}/contents/{content_id}` | Update variant |
| GET | `/v1/newsletters/{id}/metrics` | Newsletter metrics |
| GET | `/v1/newsletters/{id}/metrics/links` | Newsletter link metrics |
| GET | `/v1/newsletters/{id}/contents/{content_id}/metrics` | Variant metrics |
| GET | `/v1/newsletters/{id}/contents/{content_id}/metrics/links` | Variant link metrics |
| GET | `/v1/newsletters/{id}/messages` | Newsletter messages |
| GET | `/v1/newsletters/{id}/language/{language}` | Newsletter translation |
| PUT | `/v1/newsletters/{id}/language/{language}` | Update translation |
| GET | `/v1/newsletters/{id}/test_groups` | List test groups |
| GET | `/v1/newsletters/{id}/test_groups/{group_id}/language/{language}` | Test group translation |
| PUT | `/v1/newsletters/{id}/test_groups/{group_id}/language/{language}` | Update test group translation |

### Messages (3 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/messages` | List messages (paginated) |
| GET | `/v1/messages/{id}` | Get message details |
| GET | `/v1/messages/{id}/archived_message` | Get archived message body |

### Transactional Messages (9 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/transactional` | List transactional templates |
| GET | `/v1/transactional/{id}` | Get transactional template |
| GET | `/v1/transactional/{id}/metrics` | Template metrics |
| GET | `/v1/transactional/{id}/metrics/links` | Template link metrics |
| GET | `/v1/transactional/{id}/content` | List template variants |
| PUT | `/v1/transactional/{id}/content` | Update template content |
| GET | `/v1/transactional/{id}/language/{language}` | Template translation |
| PUT | `/v1/transactional/{id}/language/{language}` | Update translation |
| GET | `/v1/transactional/{id}/deliveries` | Template delivery log |

### Send Messages (3 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/send/email` | Send transactional email |
| POST | `/v1/send/push` | Send transactional push |
| POST | `/v1/send/sms` | Send transactional SMS |

### Exports (5 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/exports` | List exports |
| POST | `/v1/exports/customers` | Create customer export |
| POST | `/v1/exports/deliveries` | Create deliveries export |
| GET | `/v1/exports/{id}` | Get export status |
| GET | `/v1/exports/{id}/download` | Download export file |

### Imports (2 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/imports` | Create bulk import |
| GET | `/v1/imports/{id}` | Get import status |

### Activities (1 endpoint)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/activities` | List workspace-wide activities (last 30 days) |

### Snippets (3 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/snippets` | List snippets |
| PUT | `/v1/snippets` | Create/update snippets (batch upsert) |
| DELETE | `/v1/snippets/{name}` | Delete snippet |

### Reporting Webhooks (5 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/reporting_webhooks` | List webhooks |
| POST | `/v1/reporting_webhooks` | Create webhook |
| GET | `/v1/reporting_webhooks/{id}` | Get webhook |
| PUT | `/v1/reporting_webhooks/{id}` | Update webhook |
| DELETE | `/v1/reporting_webhooks/{id}` | Delete webhook |

### Sender Identities (3 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/sender_identities` | List sender identities |
| GET | `/v1/sender_identities/{id}` | Get sender identity |
| GET | `/v1/sender_identities/{id}/used_by` | Get identity usage |

### ESP Suppression (4 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/esp_suppression/search` | Search suppressions |
| GET | `/v1/esp_suppression/{email}` | Get suppression for email |
| PUT | `/v1/esp_suppression/{email}` | Suppress email at ESP level |
| DELETE | `/v1/esp_suppression/{email}` | Unsuppress email |

### Subscription Topics (1 endpoint)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/subscription_topics` | List subscription topics |

### Data Index (2 endpoints)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/index/attributes` | List indexed attributes |
| GET | `/v1/index/events` | List indexed events |

### Workspaces (1 endpoint)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/workspaces` | List workspaces |

### Info (1 endpoint)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/info/ip_addresses` | CIO sending IP addresses |

## Pipelines API (7 endpoints, not implemented)

CDP-style POST-only ingest. Uses different auth and data model. Not implemented in our SDK.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/identify` | Identify person |
| POST | `/v1/track` | Track event |
| POST | `/v1/group` | Group (object) |
| POST | `/v1/page` | Page view |
| POST | `/v1/screen` | Screen view |
| POST | `/v1/alias` | Alias |
| POST | `/v1/batch` | Batch operations |

## Our SDK Implementation

| File | Class | Coverage |
|------|-------|---------|
| `teehouse/integrations/customerio/track.py` | `CustomerIOTrackClient` | All 18 Track API endpoints (V2 + V1) |
| `teehouse/integrations/customerio/app.py` | `CustomerIOAppClient` | All 115 App API endpoints |
| `teehouse/integrations/customerio/models.py` | 59 Pydantic models | All request/response types |
| `teehouse/integrations/customerio/__init__.py` | Re-exports | Backward-compatible imports |
| `tests/integrations/test_customerio.py` | 90 tests | Track + App methods + models |

## CIO Filter Syntax (for search_customers)

```json
{
  "and": [
    {"segment": {"id": 3}},
    {"attribute": {"field": "email", "operator": "exists"}},
    {"not": {"attribute": {"field": "unsubscribed", "operator": "eq", "value": "true"}}}
  ]
}
```

Operators: `eq`, `ne`, `exists`, `not_exists`, `starts_with`, `ends_with`, `contains`, `gt`, `gte`, `lt`, `lte`.

## Common Gotchas

1. **Email as identifier**: When looking up by email, the email is part of the URL path and must be URL-encoded (`@` → `%40`, `+` → `%2B`).
2. **Multiple profiles per email**: CIO allows multiple profiles with the same email. Use `get_customers_by_email` to find all.
3. **Broadcast trigger rate limit**: 1 request per 10 seconds. Will get 429 if exceeded.
4. **Merge is destructive**: `merge_customers` permanently deletes the secondary profile.
5. **Suppress vs Delete**: Suppress keeps the profile but stops messaging; delete removes everything.
6. **Collection schema**: The API returns `"schema"` key which is a Python reserved-ish word — our model uses `schema_fields` with `alias="schema"`.
7. **Send email `from` field**: `from` is a Python keyword — our model uses `from_field` with `alias="from"`.
