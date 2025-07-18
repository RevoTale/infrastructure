# NGINX Requirements Summary for Claude Sonnet 4

## Project Overview
**Universal Edge Cache + TLS Terminator** for Next.js applications  
**Date**: July 2025  
**Purpose**: Production-ready shared cache proxy that fronts all sub-domains

---

## Critical Requirements (NON-NEGOTIABLE)

### 1. UNIVERSAL LOCATION BLOCK
- ✅ **ONLY ONE location block** (`location /`)
- ❌ **NO separate locations** for `/api/`, static files, or any other paths
- **Rationale**: All requests must be handled uniformly by upstream server

### 2. TRUE SHARED CACHE (Anti-Abuse)
- ✅ **IGNORE all client request headers** (Cache-Control, Pragma from browsers)
- ✅ **IGNORE all cookies** (auth-token, session-id, any cookies)
- ✅ **IGNORE Authorization headers**
- ❌ **NO proxy_cache_bypass for client-side conditions**
- **Rationale**: Prevent cache abuse, maximize cache hit ratio

### 3. UPSTREAM CACHE CONTROL RESPECT
- ✅ **ONLY respect upstream response headers**:
  - `Cache-Control: private` → don't cache
  - `Cache-Control: no-cache` → don't cache
  - `Cache-Control: no-store` → don't cache
  - `Cache-Control: max-age=0` → don't cache
- ✅ **NO proxy_cache_valid directives** (let upstream control duration)
- **Rationale**: Upstream server has full control over caching behavior

### 4. POST REQUEST CACHING
- ✅ **Allow POST caching** if upstream sends cacheable headers
- ✅ **Remove POST from bypass conditions**
- ✅ **proxy_cache_methods GET HEAD POST**
- **Rationale**: Modern APIs use POST for complex queries that can be cached

### 5. SCHEMA-AWARE CACHE KEYS
- ✅ **Include $scheme in cache key** for HTTP/HTTPS separation
- ✅ **Cache key format**: `$scheme$host$request_uri$http_accept_encoding`
- ✅ **Standard Vary handling**: nginx processes upstream Vary headers according to HTTP standards
- **Rationale**: Simple, standards-compliant cache key with upstream-controlled content negotiation

### 6. WEBSOCKET SUPPORT
- ✅ **Proper WebSocket upgrade headers**
- ✅ **Long timeouts** (3600s) for persistent connections
- ✅ **Connection upgrade mapping**
- **Rationale**: Next.js development and real-time features

### 7. INITIAL LOAD OPTIMIZATION
- ✅ **Rate limiting**: 50 req/s with 150+ burst
- ✅ **Support 50+ files** in initial page load without blocking
- **Rationale**: Modern web apps require many resources on first load

---

## Cache Bypass Logic (MINIMAL)

### ONLY These Conditions Bypass Cache:
1. **Upstream response headers** (`$no_cache`):
   - `Cache-Control: private`
   - `Cache-Control: no-cache`
   - `Cache-Control: no-store`
   - `Cache-Control: max-age=0`

2. **HTTP methods** (`$bypass_cache`):
   - `PUT` (destructive)
   - `DELETE` (destructive)
   - `PATCH` (destructive)
   - ❌ **NOT POST** (POST can be cached)

3. **WebSocket upgrades** (`$http_upgrade`):
   - When `Upgrade: websocket` header present

### EXPLICITLY IGNORED (NO BYPASS):
- ❌ Client `Cache-Control` headers
- ❌ Client `Pragma` headers
- ❌ Any cookies (auth, session, etc.)
- ❌ `Authorization` headers
- ❌ Any client-side cache-busting attempts

---

## Content Negotiation Support

### Standard Cache Key Strategy:
Using nginx's standard caching with upstream Vary header handling:

```
proxy_cache_key "$scheme$host$request_uri$http_accept_encoding";
```

**Cache Key Components:**
- `$scheme` - HTTP vs HTTPS separation
- `$host` - Domain separation  
- `$request_uri` - Path and query parameters
- `$http_accept_encoding` - Always included for compression support

**How nginx handles Vary headers:**

nginx automatically processes Vary headers from upstream responses according to HTTP standards. According to the [official nginx documentation](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_valid), when upstream sends a `Vary` header, nginx processes it as follows:

1. **Vary: * (wildcard)**: Response is not cached at all
2. **Vary: Accept, Accept-Language**: Response is cached, but nginx considers the corresponding request headers when serving from cache
3. **No Vary header**: Response is cached once and served to all clients regardless of their headers

**From nginx docs**: *"If the header includes the "Vary" field with the special value "*", such a response will not be cached. If the header includes the "Vary" field with another value, such a response will be cached taking into account the corresponding request header fields."*

This automatic Vary handling means you don't need to configure anything - nginx handles content negotiation automatically based on upstream response headers.

**Example Flow:**

```
# Static asset (no Vary header from upstream)
GET /style.css
Accept: text/css,*/*
→ Upstream: (no Vary header)  
→ Cache key: https://example.com/style.css$gzip
→ All future requests served from this cache entry

# Image endpoint with content negotiation
GET /api/image.jpg
Accept: image/webp,*/*
→ Upstream: Vary: Accept
→ Cached with consideration of Accept header
→ Different Accept values may result in cache misses

# API endpoint rejecting caching
GET /api/personalized
→ Upstream: Vary: *
→ Not cached at all
```

**Benefits:**
- **Standards compliant**: Uses nginx's built-in Vary handling
- **Upstream controlled**: Server decides caching behavior via Vary headers
- **Simple configuration**: No complex mapping needed
- **Automatic**: Works with any upstream that follows HTTP caching standards

**Limitations:**
- nginx's Vary handling is limited compared to dedicated CDNs
- For complex content negotiation, you might get cache misses
- No fine-grained control over which headers are considered

This approach prioritizes simplicity and standards compliance while still allowing upstream servers to control content negotiation via standard HTTP headers.

---

## Performance Configuration

### Cache Settings:
- **Size**: 5GB max cache, 200MB metadata
- **Hierarchy**: 2-level directory structure
- **Inactive**: 180 minutes
- **Background updates**: Enabled for seamless refreshes
- **Lock timeout**: 5 seconds

### Buffers:
- **Proxy buffer size**: 128k
- **Proxy buffers**: 4 × 256k
- **Busy buffers**: 256k

### Timeouts:
- **Connect**: 10 seconds (quick connection)
- **Send/Read**: 3600 seconds (WebSocket compatible)

---

## Security Requirements

### TLS Configuration:
- **Protocols**: TLS 1.2, TLS 1.3 only
- **Cipher preference**: Client preference (modern approach)
- **Session cache**: 50MB shared
- **OCSP stapling**: Enabled

### Security Headers:
- **HSTS**: 2 years with preload
- **X-Content-Type-Options**: nosniff
- **X-Frame-Options**: DENY
- **X-XSS-Protection**: enabled
- **Referrer-Policy**: strict-origin-when-cross-origin

---

## Rate Limiting Strategy

### Zone Configuration:
- **Zone**: `general` (single zone for all requests)
- **Rate**: 50 requests/second
- **Burst**: 150 requests (accommodates initial loads)
- **Behavior**: `nodelay` (immediate processing within burst)

### Rationale:
- Initial page loads can require 50+ resources
- Burst allows immediate loading without delays
- Rate limit prevents abuse while allowing normal usage

---

## Monitoring & Debugging

### Headers Added:
- `X-Cache-Status` - Shows HIT/MISS/BYPASS/EXPIRED for debugging
- Security headers (HSTS, etc.) for protection

### Access Logs:
- **Disabled** for performance (can be enabled for debugging)

---

## Architecture Decisions

### Why Universal Location:
- Upstream server controls routing logic
- Simplified nginx configuration
- Consistent caching behavior across all endpoints
- No nginx-level path-based logic

### Why Ignore Client Headers:
- Prevents cache abuse by malicious users
- Maximizes cache efficiency
- CDN-like behavior
- Only server decides what should be cached

### Why Allow POST Caching:
- Modern APIs use POST for complex queries
- GraphQL typically uses POST
- Search/filter APIs can benefit from caching
- Server controls with Cache-Control headers

---

## Common Mistakes to Avoid

❌ **Adding separate location blocks** for static files or API  
❌ **Adding cookie-based cache bypassing**  
❌ **Adding client header-based bypassing**  
❌ **Using proxy_cache_valid to override upstream**  
❌ **Blocking POST requests from caching**  
❌ **Aggressive rate limiting that blocks initial loads**  
❌ **Short timeouts that break WebSocket connections**

---

## Testing Scenarios

### Expected Behaviors:
1. **Client sends `Cache-Control: no-cache`** → Served from cache if available
2. **Client has auth cookies** → Served from cache if available  
3. **Upstream sends `Cache-Control: private`** → Not cached
4. **Upstream sends `Cache-Control: max-age=3600`** → Cached for 1 hour
5. **POST with cacheable response** → Cached if upstream allows
6. **WebSocket upgrade** → Proxied with long timeouts
7. **Initial page load (50+ files)** → No rate limiting blocks

---

## File Structure
```
/workspaces/backend/nginx/
├── nginx.conf                    # Main configuration
└── NGINX_REQUIREMENTS_SUMMARY.md # This requirements document
```

---

**Last Updated**: July 4, 2025  
**Configuration Status**: Production Ready ✅  
**Validation Status**: Syntax Valid ✅
