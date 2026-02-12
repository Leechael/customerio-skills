---
name: customerio
description: Interact with the Customer.io App API via the cio CLI. Use this skill when the user wants to manage Customer.io resources such as customers, segments, campaigns, broadcasts, newsletters, collections, exports, messages, objects, webhooks, or transactional messages. Also use when the user asks to query, filter, create, update, or delete Customer.io data from the command line.
---

# Customer.io CLI (`cio`)

A CLI tool wrapping the Customer.io App API. Requires `CUSTOMERIO_API_TOKEN` environment variable.

## Quick Reference

```bash
# Global flags
cio --region eu ...          # Use EU region (default: us)
cio ... --jq '.field'        # Filter JSON output with jq expression

# Customers
cio customers ls                                     # List customers (POST with empty filter)
cio customers ls --body '{"ids":["u1","u2"]}'        # Filter by IDs
cio customers search --email user@example.com        # Search by email
cio customers get <id>                               # Get attributes
cio customers activities <id>                        # Get activities
cio customers segments <id>                          # Get segments
cio customers messages <id>                          # Get messages

# Segments
cio segments ls                                      # List all segments
cio segments get <id>                                # Get segment details
cio segments create --body '{"segment":{"name":"VIP"}}'
cio segments rm <id>                                 # Delete segment
cio segments count <id>                              # Customer count
cio segments members <id>                            # Membership list

# Campaigns
cio campaigns ls                                     # List campaigns
cio campaigns get <id>                               # Get campaign
cio campaigns metrics <id>                           # Get metrics
cio campaigns actions <id>                           # List actions

# Broadcasts
cio broadcasts ls                                    # List broadcasts
cio broadcasts get <id>                              # Get broadcast
cio broadcasts trigger <id>                          # Trigger (optional filter body)
cio broadcasts metrics <id>                          # Get metrics

# Newsletters
cio newsletters ls                                   # List newsletters
cio newsletters get <id>                             # Get newsletter
cio newsletters metrics <id>                         # Get metrics

# Transactional
cio transactional ls                                 # List transactional messages
cio transactional get <id>                           # Get message
cio transactional deliveries <id>                    # Get deliveries

# Send messages
cio send email --body '{"to":"u@e.com","transactional_message_id":"1",...}'
cio send push --body '...'
cio send sms --body '...'

# Collections
cio collections ls                                   # List collections
cio collections get <id>                             # Get collection
cio collections create --body '{"name":"n","data":[...]}'
cio collections update <id> --body '...'
cio collections content <id>                         # GET or PUT content

# Exports
cio exports ls                                       # List exports
cio exports create-customers                         # Export all customers
cio exports create-customers --body '{"filters":{...}}'
cio exports create-deliveries                        # Export all deliveries
cio exports get <id>                                 # Get export status
cio exports download <id>                            # Get download URL

# Objects
cio objects get <type> <id>                          # Get object
cio objects search                                   # Search objects

# Other resources
cio messages get <id>                                # Get message
cio messages deliveries <id>                         # Get deliveries
cio webhooks ls                                      # List webhooks
cio webhooks get <id>                                # GET or PUT webhook
cio webhooks create --body '...'                     # Create webhook
cio sender-identities ls                             # List sender identities
cio snippets ls                                      # List snippets
cio snippets upsert --body '...'                     # Create/update snippet
cio esp-suppression ls                               # List suppressions
cio esp-suppression search                           # Search (optional filter)
cio imports create --body '...'                      # Create import
cio info ip-addresses                                # Get IP addresses
cio workspaces ls                                    # List workspaces
```

## Body Input

Commands that accept a body support two input methods:

```bash
# --body flag
cio segments create --body '{"segment":{"name":"Test"}}'

# Pipe via stdin
echo '{"segment":{"name":"Test"}}' | cio segments create
```

## JQ Filtering

Use `--jq` to extract specific fields from any response:

```bash
cio segments ls --jq '.segments[].name'
cio customers get u1 --jq '.customer.email'
cio campaigns ls --jq '.campaigns | map(select(.active)) | length'
```
