# Design Spec: Hybrid Cache Implementation (L1 + L2 + HTTP Headers)

**Date:** 2026-07-04  
**Author:** Antigravity  

This document outlines the design specification for implementing a hybrid, multi-layer caching system for Jikan API requests in the `koko` Go backend to prevent rate limiting, reduce backend resource usage, and improve response times.

---

## 1. Problem Statement & Objectives

Currently, the backend proxies requests to the public Jikan API (MyAnimeList) and employs a simple, non-persistent, in-memory cache for 5 minutes. 
### Problems:
1. **Non-persistent:** If the backend server restarts, the cache is lost. Every restart causes cold cache fetches to Jikan, leading to potential `429 Too Many Requests` API rate limits.
2. **Memory Leak Risk:** The current in-memory cache map does not clean up expired keys, which causes memory footprint to grow infinitely over time as new queries are made.
3. **Cache Stampede (Dogpiling):** If multiple client requests hit the backend for the same expired or missing Jikan URL simultaneously, the backend will perform duplicate concurrent Jikan HTTP queries.
4. **No Client-side Cache:** The backend does not send HTTP `Cache-Control` headers, meaning the Nuxt frontend makes requests to the backend for every page navigation even if the data was loaded seconds ago.

### Objectives:
* Implement **L2 Persistent Cache** in PostgreSQL (Supabase) to keep cache across server restarts.
* Implement **L1 In-Memory Cache** (RAM) to handle high-frequency requests without hitting the DB.
* Add **Cache Stampede Protection** (Singleflight) so concurrent requests only trigger one fetch.
* Implement **HTTP Cache-Control Headers** to allow browser and client caching.
* Implement **Periodic Cleanups** to evict expired cache entries from RAM and Database.

---

## 2. System Architecture

```
                 +-----------------------------------------+
                 |            Nuxt.js Frontend             |
                 +-----------------------------------------+
                                      |
                                      | HTTP GET Request
                                      v
                 +-----------------------------------------+
                 |            Go Backend (Gin)             |
                 +-----------------------------------------+
                                      |
                                      | 1. Check L1 Memory Cache
                                      v
                               [L1 RAM Cache] - - -> (Hit: Return < 1ms)
                                      |
                                      | 2. (Miss) Check L2 Postgres Cache
                                      v
                              [Postgres DB Cache] - - > (Hit: Cache in L1, Return 10-30ms)
                                      |
                                      | 3. (Miss) singleflight.Do()
                                      v
                             [External Jikan API] - - > (Success: Cache L2 & L1, Return)
```

---

## 3. Data Design (L2 DB Schema)

We will introduce a `jikan_cache` table to store Jikan responses:

```sql
CREATE TABLE IF NOT EXISTS jikan_cache (
    key TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index to optimize cleanup queries
CREATE INDEX IF NOT EXISTS idx_jikan_cache_expires_at ON jikan_cache(expires_at);
```

### TTL (Time to Live) Rules:
* **Anime Details (`/api/anime/:id`):** 24 Hours (highly static)
* **Genres (`/api/genres`):** 24 Hours (highly static)
* **Anime List / Searches (`/api/anime?q=...`):** 1 Hour
* **Recommendations (`/api/recommendations`):** 6 Hours

---

## 4. Go Implementation Details

### A. L1 Cache Cleanup
We will run a background goroutine in `main.go` that runs every 5 minutes to delete expired items from the L1 map.

```go
func (c *Cache) StartCleanupLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            c.mu.Lock()
            now := time.Now()
            for k, v := range c.items {
                if now.After(v.ExpiresAt) {
                    delete(c.items, k)
                }
            }
            c.mu.Unlock()
        }
    }()
}
```

### B. L2 Database Cache Cleanup
We will periodically delete expired records from the `jikan_cache` table (e.g., every 1 hour or on server startup).

```go
func CleanExpiredDBCache() {
    query := `DELETE FROM jikan_cache WHERE expires_at < CURRENT_TIMESTAMP`
    _, err := DB.Exec(query)
    if err != nil {
        log.Printf("Error cleaning expired DB cache: %v\n", err)
    }
}
```

### C. Singleflight (Stampede Protection)
We will use Go's standard library/semi-standard `golang.org/x/sync/singleflight` to ensure that duplicate concurrent requests only hit Jikan or the Database once.

```go
import "golang.org/x/sync/singleflight"

var sfGroup singleflight.Group
```

In `fetchAndCache`:
```go
// Wrap the DB/Jikan fetch with singleflight to avoid dogpiling
v, err, _ := sfGroup.Do(targetURL, func() (interface{}, error) {
    // 1. Check L2 DB Cache
    // 2. Fetch from Jikan API if DB Miss
    // 3. Save to L2 DB Cache
    return result, nil
})
```

### D. HTTP Headers
All successful GET Jikan responses will return the header:
```go
c.Writer.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
```

---

## 5. Verification & Testing

### Automated Checks
* **Go Tests:** Add a test suite verifying that cache hits do not result in subsequent HTTP calls.
* **Self-Test Verification:** Update `runSelfTests()` in `main.go` to test:
  1. L1 RAM Cache Hit (Fast execution)
  2. L2 Database Cache Hit (Restarting cache simulator)
  3. `Cache-Control` header correctness

### Manual Check
* Check browser Network tab to verify that `Cache-Control` is set and subsequent navigations show `(disk cache)` or `(memory cache)`.
