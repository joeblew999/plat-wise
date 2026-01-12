# Wise API Coverage

This document tracks the implementation status of the [Wise API](https://docs.wise.com/api-reference).

## Legend

- [x] Implemented
- [ ] Not implemented
- [-] Partially implemented

---

## Authentication

| Feature | Status | Notes |
|---------|--------|-------|
| API Token (Bearer) | [x] | Implemented in client.go |
| SCA (Strong Customer Authentication) | [ ] | Not implemented |
| OAuth 2.0 | [ ] | Not implemented |
| Webhook Signatures | [ ] | Not implemented |

---

## Profiles API

**Endpoint:** `/v1/profiles`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v1/profiles` | [x] | `Profiles.List()` |
| GET | `/v1/profiles/{profileId}` | [x] | `Profiles.Get()` |
| POST | `/v1/profiles` | [x] | `Profiles.CreatePersonal()`, `Profiles.CreateBusiness()` |
| PUT | `/v1/profiles/{profileId}` | [ ] | Update profile |

---

## Quotes API

**Endpoints:** `/v2/quotes`, `/v3/profiles/{profileId}/quotes`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| POST | `/v2/quotes` | [x] | `Quotes.CreateV2()` |
| GET | `/v2/quotes/{quoteId}` | [x] | `Quotes.GetV2()` |
| POST | `/v3/profiles/{profileId}/quotes` | [x] | `Quotes.Create()` |
| GET | `/v3/profiles/{profileId}/quotes/{quoteId}` | [x] | `Quotes.Get()` |
| PATCH | `/v3/profiles/{profileId}/quotes/{quoteId}` | [x] | `Quotes.Update()` |

---

## Recipients (Accounts) API

**Endpoint:** `/v1/accounts`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| POST | `/v1/accounts` | [x] | `Recipients.Create()` |
| GET | `/v1/accounts/{accountId}` | [x] | `Recipients.Get()` |
| GET | `/v1/accounts` | [x] | `Recipients.List()` |
| DELETE | `/v1/accounts/{accountId}` | [x] | `Recipients.Delete()` |
| GET | `/v1/account-requirements` | [x] | `Recipients.GetRequirements()` |
| POST | `/v1/account-requirements` | [ ] | Refresh requirements |
| GET | `/v1/quotes/{quoteId}/account-requirements` | [ ] | Quote-specific requirements |

---

## Transfers API

**Endpoint:** `/v1/transfers`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| POST | `/v1/transfers` | [x] | `Transfers.Create()` |
| GET | `/v1/transfers/{transferId}` | [x] | `Transfers.Get()` |
| GET | `/v1/transfers` | [x] | `Transfers.List()` |
| PUT | `/v1/transfers/{transferId}/cancel` | [x] | `Transfers.Cancel()` |
| POST | `/v3/profiles/{profileId}/transfers/{transferId}/payments` | [x] | `Transfers.Fund()` |
| GET | `/v1/transfers/{transferId}/issues` | [x] | `Transfers.GetIssues()` |
| GET | `/v1/delivery-estimates/{transferId}` | [x] | `Transfers.GetDeliveryTime()` |
| GET | `/v1/transfers/{transferId}/receipt.pdf` | [ ] | Download receipt |
| GET | `/v3/profiles/{profileId}/transfers/{transferId}/activities` | [ ] | Transfer activities |

---

## Exchange Rates API

**Endpoint:** `/v1/rates`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v1/rates` | [x] | `ExchangeRates.List()` |
| GET | `/v1/rates?source={}&target={}` | [x] | `ExchangeRates.Get()` |
| GET | `/v1/rates?time={}` | [x] | `ExchangeRates.GetHistorical()` |

---

## Balances API

**Endpoint:** `/v4/profiles/{profileId}/balances`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v4/profiles/{profileId}/balances` | [x] | `Balances.List()` |
| GET | `/v4/profiles/{profileId}/balances/{balanceId}` | [x] | `Balances.Get()` |
| POST | `/v2/profiles/{profileId}/balance-movements` | [x] | `Balances.Convert()` |
| GET | `/v1/profiles/{profileId}/balance-statements/{balanceId}/statement.json` | [x] | `Balances.GetStatement()` |
| POST | `/v3/profiles/{profileId}/balances` | [ ] | Create balance |
| DELETE | `/v3/profiles/{profileId}/balances/{balanceId}` | [ ] | Delete balance |

---

## Borderless Accounts API

**Endpoint:** `/v1/borderless-accounts`

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v1/borderless-accounts` | [ ] | List borderless accounts |
| GET | `/v1/borderless-accounts/{accountId}` | [ ] | Get borderless account |
| GET | `/v1/borderless-accounts/{accountId}/statement.json` | [ ] | Get statement |

---

## Bank Details API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v1/profiles/{profileId}/account-details` | [ ] | Get account details |
| POST | `/v1/profiles/{profileId}/account-details` | [ ] | Create account details |

---

## Cards API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v3/profiles/{profileId}/cards` | [ ] | List cards |
| GET | `/v3/profiles/{profileId}/cards/{cardId}` | [ ] | Get card |
| POST | `/v3/profiles/{profileId}/cards` | [ ] | Order card |
| PUT | `/v3/profiles/{profileId}/cards/{cardId}/status` | [ ] | Update card status |
| GET | `/v3/profiles/{profileId}/cards/{cardId}/sensitive-details` | [ ] | Get card details |

---

## Webhooks API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| POST | `/v3/profiles/{profileId}/subscriptions` | [ ] | Create subscription |
| GET | `/v3/profiles/{profileId}/subscriptions` | [ ] | List subscriptions |
| DELETE | `/v3/profiles/{profileId}/subscriptions/{subscriptionId}` | [ ] | Delete subscription |
| GET | `/v3/profiles/{profileId}/subscriptions/{subscriptionId}/events` | [ ] | List events |

---

## Multi-Currency Account API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v2/profiles/{profileId}/multi-currency-account` | [ ] | Get MCA |
| POST | `/v2/profiles/{profileId}/multi-currency-account` | [ ] | Create MCA |

---

## Batch Payments API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| POST | `/v3/profiles/{profileId}/batch-payments` | [ ] | Create batch |
| GET | `/v3/profiles/{profileId}/batch-payments/{batchId}` | [ ] | Get batch |
| GET | `/v3/profiles/{profileId}/batch-payments` | [ ] | List batches |

---

## Direct Debits API

| Method | Endpoint | Status | Function |
|--------|----------|--------|----------|
| GET | `/v1/profiles/{profileId}/direct-debit-mandates` | [ ] | List mandates |
| POST | `/v1/profiles/{profileId}/direct-debit-mandates` | [ ] | Create mandate |

---

## Summary

### Implemented (Core APIs)

| Service | Endpoints | Coverage |
|---------|-----------|----------|
| Profiles | 3/4 | 75% |
| Quotes | 5/5 | 100% |
| Recipients | 5/7 | 71% |
| Transfers | 7/9 | 78% |
| Exchange Rates | 3/3 | 100% |
| Balances | 5/7 | 71% |

### Not Implemented

- Borderless Accounts API
- Bank Details API
- Cards API
- Webhooks API
- Multi-Currency Account API
- Batch Payments API
- Direct Debits API
- SCA/OAuth Authentication

---

## Environment

| Variable | Description |
|----------|-------------|
| `WISE_API_TOKEN` | API token from Wise |

## Endpoints

| Environment | Base URL |
|-------------|----------|
| Production | `https://api.wise.com` |
| Sandbox | `https://api.sandbox.transferwise.tech` |
