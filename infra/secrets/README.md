# Chia Florist — Secrets Management Guide

This document describes how to provision and inject secrets for each deployment platform.

## Required Secrets Per Module

### service-core
| Variable | Description | Sensitive |
|---|---|---|
| `JWT_SECRET` | JWT signing key | ✅ |
| `POSTGRES_DB_HOST` | Supabase pooler host | ✅ |
| `POSTGRES_DB_NAME` | Database name | — |
| `POSTGRES_DB_USER` | Database user | ✅ |
| `POSTGRES_DB_PASSWORD` | Database password | ✅ |
| `POSTGRES_DB_PORT` | Database port (6543) | — |
| `POSTGRES_DB_SSLMODE` | `disable` or `require` | — |
| `SUPABASE_PROJECT_URL` | Supabase project URL | — |
| `SUPABASE_SUPA_KEY` | Supabase service role JWT | ✅ |
| `STORAGE_BUCKET` | Storage bucket name | — |
| `STORAGE_SIGNED_URL_EXPIRY` | URL expiry seconds | — |
| `RAJAONGKIR_DESTINATION_URL` | Komerce destination endpoint | — |
| `RAJAONGKIR_CALCULATE_URL` | Komerce calculate endpoint | — |
| `RAJAONGKIR_SHIPPING` | Komerce shipping API key | ✅ |
| `RAJAONGKIR_PAYMENT` | Komerce payment API key | ✅ |
| `SMTP_HOST` | SMTP server host | — |
| `SMTP_PORT` | SMTP port (587) | — |
| `SMTP_USERNAME` | SMTP user email | — |
| `SMTP_PASSWORD` | SMTP app password | ✅ |
| `SMTP_FROM` | From display name + email | — |

### e-commerce
| Variable | Description |
|---|---|
| `SERVICE_CORE_API_URL` | Internal URL of service-core |
| `SUPABASE_CHIA_URL` | Supabase project URL |
| `SUPABASE_CHIA_KEY` | Supabase anon/service key |

### control-panel
| Variable | Description |
|---|---|
| `VITE_API_BASE_URL` | Public URL of service-core API |

## Platform Guides

### Railway (Primary)
