# Security & Architectural Improvements TODO

This document tracks security and architectural improvements identified during code audit.
Items are prioritized by severity and effort required.

## Backend (REST API)

### Critical Priority

- [ ] **Token Refresh Mechanism**
  - Location: `internal/handlers/auth.go`
  - Issue: No refresh token flow, tokens valid for 24h with no revocation
  - Solution: Implement refresh tokens with shorter-lived access tokens (15 min)
  - Add token blacklist/revocation table for logout

- [ ] **Context Timeout for DB Queries**
  - Location: `internal/service/*.go`
  - Issue: Context passed to services but not used for DB query cancellation
  - Solution: Pass context to GORM queries, use `context.WithTimeout()` for long operations

- [ ] **Distributed Rate Limiting**
  - Location: `internal/middleware/ratelimit.go`
  - Issue: In-memory rate limiter only works for single-instance deployments
  - Solution: Implement Redis-backed rate limiting for production

### High Priority

- [ ] **Audit Logging**
  - Location: `internal/handlers/user.go`, `internal/handlers/auth.go`
  - Issue: No audit trail for sensitive operations (superadmin changes, failed logins)
  - Solution: Add audit log table, log IP address, timestamp, user agent

- [ ] **N+1 Query Optimization**
  - Location: `internal/store/child.go:46`, `internal/store/employee.go`
  - Issue: Preloading contracts then fetching separately
  - Solution: Review if preloads are used in response, consider lazy loading

### Medium Priority

- [ ] **Soft Deletes**
  - Location: Model definitions
  - Issue: Physical deletes prevent data recovery and compliance auditing
  - Solution: Add `deleted_at` field, use GORM soft delete features

- [ ] **Explicit Cascade Delete Documentation**
  - Location: `internal/store/employee.go:73-78`
  - Issue: Cascade delete order not documented, property deletion implicit
  - Solution: Add comments documenting delete order, verify FK constraints

---

## Frontend (Web UI)

### Critical Priority

- [ ] **Migrate JWT to httpOnly Cookies**
  - Location: `web/src/stores/auth.ts`, backend
  - Issue: JWT in localStorage vulnerable to XSS attacks
  - Solution: Store JWT in httpOnly, secure, sameSite=Strict cookie managed by backend
  - Requires backend changes for cookie handling

- [ ] **Token Refresh Handling**
  - Location: `web/src/api/client.ts`
  - Issue: No refresh token mechanism, 401 just logs out
  - Solution: Implement automatic token refresh on 401 with retry

### High Priority

- [ ] **Session Timeout**
  - Location: `web/src/stores/auth.ts`
  - Issue: No activity-based session timeout
  - Solution: Track last activity, auto-logout after idle period (15-30 min)

- [ ] **CSRF Token Handling**
  - Location: `web/src/api/client.ts`
  - Issue: No CSRF protection for state-changing requests
  - Solution: Add CSRF token header for POST/PUT/DELETE requests
  - Requires backend to provide CSRF tokens

- [ ] **Rate Limit Response Handling**
  - Location: `web/src/api/client.ts`
  - Issue: No handling for 429 (rate limited) responses
  - Solution: Add response interceptor for 429, show user-friendly message

### Medium Priority

- [ ] **Request Correlation IDs**
  - Location: `web/src/api/client.ts`, backend
  - Issue: No request tracking for debugging
  - Solution: Generate correlation ID per request, pass in header

- [ ] **Error Boundary Component**
  - Location: `web/src/` (new component)
  - Issue: No global error handling for unexpected errors
  - Solution: Create Vue error boundary component

- [ ] **AsyncState Pattern**
  - Location: Views with async operations
  - Issue: Inconsistent loading/error state handling
  - Solution: Implement standard AsyncState type for all async operations

---

## Implementation Notes

### Token Refresh Flow (Recommended Implementation)

```
1. Access token: 15 minute expiry
2. Refresh token: 7 day expiry, stored in httpOnly cookie
3. On 401, attempt refresh automatically
4. On refresh failure, logout user
5. Implement token rotation on refresh
```

### CSRF Implementation

```
1. Backend generates CSRF token on login, returns in response
2. Frontend stores CSRF token in memory (not localStorage)
3. Frontend sends CSRF token in X-CSRF-Token header for mutations
4. Backend validates CSRF token matches session
```

### Audit Log Schema

```sql
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id INTEGER,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id INTEGER,
    ip_address INET,
    user_agent TEXT,
    details JSONB
);
```
