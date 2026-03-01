---
title: API Reference
weight: 3
---

KitaManager provides a REST API with interactive OpenAPI/Swagger documentation available at `/swagger/index.html` when running the application. All endpoints except login and token refresh require JWT authentication. Mutating requests (POST, PUT, DELETE) require a CSRF token via the `X-CSRF-Token` header.

## Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/login` | Authenticate and receive access + refresh tokens |
| POST | `/api/v1/refresh` | Refresh an expired access token |
| POST | `/api/v1/logout` | Invalidate the current session |
| GET | `/api/v1/me` | Get the current user's profile |
| PUT | `/api/v1/me/password` | Change the current user's password |

### Login Example

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "admin123"}'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Using the Token

Include the token in the `Authorization` header for all subsequent requests:

```bash
curl http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Organizations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/organizations` | List organizations |
| POST | `/api/v1/organizations` | Create organization (superadmin) |
| GET | `/api/v1/organizations/{orgId}` | Get organization |
| PUT | `/api/v1/organizations/{orgId}` | Update organization |
| DELETE | `/api/v1/organizations/{orgId}` | Delete organization (superadmin) |

## Sections

All section endpoints are scoped to an organization: `/api/v1/organizations/{orgId}/sections`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../sections` | List sections |
| POST | `.../sections` | Create section |
| GET | `.../sections/{sectionId}` | Get section |
| PUT | `.../sections/{sectionId}` | Update section |
| DELETE | `.../sections/{sectionId}` | Delete section |

## Employees

All employee endpoints are scoped to an organization: `/api/v1/organizations/{orgId}/employees`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../employees` | List employees |
| POST | `.../employees` | Create employee |
| GET | `.../employees/{id}` | Get employee |
| PUT | `.../employees/{id}` | Update employee |
| DELETE | `.../employees/{id}` | Delete employee |
| GET | `.../employees/export/excel` | Export employees to Excel |
| GET | `.../employees/export/yaml` | Export employees to YAML |
| POST | `.../employees/import` | Import employees from YAML |
| GET | `.../employees/step-promotions` | Get step promotion eligibility |

### Employee Contracts

Nested under an employee: `.../employees/{id}/contracts`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../contracts` | List contracts |
| POST | `.../contracts` | Create contract |
| GET | `.../contracts/current` | Get current active contract |
| GET | `.../contracts/{contractId}` | Get contract |
| PUT | `.../contracts/{contractId}` | Update contract |
| DELETE | `.../contracts/{contractId}` | Delete contract |

## Children

All child endpoints are scoped to an organization: `/api/v1/organizations/{orgId}/children`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../children` | List children |
| POST | `.../children` | Create child |
| GET | `.../children/{id}` | Get child |
| PUT | `.../children/{id}` | Update child |
| DELETE | `.../children/{id}` | Delete child |
| GET | `.../children/export/excel` | Export children to Excel |
| GET | `.../children/export/yaml` | Export children to YAML |
| POST | `.../children/import` | Import children from YAML |
| GET | `.../children/attendance` | Org-wide attendance by date |
| GET | `.../children/attendance/summary` | Daily attendance summary |

### Child Contracts

Nested under a child: `.../children/{id}/contracts`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../contracts` | List contracts |
| POST | `.../contracts` | Create contract |
| GET | `.../contracts/current` | Get current active contract |
| GET | `.../contracts/{contractId}` | Get contract |
| PUT | `.../contracts/{contractId}` | Update contract |
| DELETE | `.../contracts/{contractId}` | Delete contract |

### Child Attendance

Nested under a child: `.../children/{id}/attendance`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `.../attendance` | Create attendance record |
| GET | `.../attendance` | List child's attendance records |
| GET | `.../attendance/{attendanceId}` | Get attendance record |
| PUT | `.../attendance/{attendanceId}` | Update attendance record |
| DELETE | `.../attendance/{attendanceId}` | Delete attendance record |

## Government Funding Rates

Global resource managed by superadmins.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/government-funding-rates` | List funding configurations |
| POST | `/api/v1/government-funding-rates` | Create funding configuration |
| GET | `/api/v1/government-funding-rates/{id}` | Get funding configuration |
| PUT | `/api/v1/government-funding-rates/{id}` | Update funding configuration |
| DELETE | `/api/v1/government-funding-rates/{id}` | Delete funding configuration |
| POST | `/api/v1/government-funding-rates/import` | Import funding rates from YAML |

### Funding Periods

Nested under a funding rate: `.../government-funding-rates/{id}/periods`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `.../periods` | Create period |
| GET | `.../periods/{periodId}` | Get period |
| PUT | `.../periods/{periodId}` | Update period |
| DELETE | `.../periods/{periodId}` | Delete period |

### Funding Properties

Nested under a period: `.../periods/{periodId}/properties`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `.../properties` | Create property |
| GET | `.../properties/{propertyId}` | Get property |
| PUT | `.../properties/{propertyId}` | Update property |
| DELETE | `.../properties/{propertyId}` | Delete property |

## Government Funding Bills

Scoped to an organization: `/api/v1/organizations/{orgId}/government-funding-bills`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../government-funding-bills` | List bills |
| POST | `.../government-funding-bills` | Upload ISBJ bill |
| GET | `.../government-funding-bills/{billId}` | Get bill |
| GET | `.../government-funding-bills/{billId}/compare` | Compare calculated vs. billed amounts |
| DELETE | `.../government-funding-bills/{billId}` | Delete bill |

## Pay Plans

Scoped to an organization: `/api/v1/organizations/{orgId}/pay-plans`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../pay-plans` | List pay plans |
| POST | `.../pay-plans` | Create pay plan |
| GET | `.../pay-plans/{id}` | Get pay plan |
| PUT | `.../pay-plans/{id}` | Update pay plan |
| DELETE | `.../pay-plans/{id}` | Delete pay plan |
| GET | `.../pay-plans/{id}/export` | Export pay plan to YAML |
| POST | `.../pay-plans/import` | Import pay plan from YAML |

### Pay Plan Periods

Nested under a pay plan: `.../pay-plans/{id}/periods`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `.../periods` | Create period |
| GET | `.../periods/{periodId}` | Get period |
| PUT | `.../periods/{periodId}` | Update period |
| DELETE | `.../periods/{periodId}` | Delete period |

### Pay Plan Entries

Nested under a period: `.../periods/{periodId}/entries`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `.../entries` | Create entry |
| GET | `.../entries/{entryId}` | Get entry |
| PUT | `.../entries/{entryId}` | Update entry |
| DELETE | `.../entries/{entryId}` | Delete entry |

## Budget Items

Scoped to an organization: `/api/v1/organizations/{orgId}/budget-items`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../budget-items` | List budget items |
| POST | `.../budget-items` | Create budget item |
| GET | `.../budget-items/{id}` | Get budget item |
| PUT | `.../budget-items/{id}` | Update budget item |
| DELETE | `.../budget-items/{id}` | Delete budget item |

### Budget Item Entries

Nested under a budget item: `.../budget-items/{id}/entries`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../entries` | List entries |
| POST | `.../entries` | Create entry |
| GET | `.../entries/{entryId}` | Get entry |
| PUT | `.../entries/{entryId}` | Update entry |
| DELETE | `.../entries/{entryId}` | Delete entry |

## Statistics

Scoped to an organization: `/api/v1/organizations/{orgId}/statistics`. All statistics endpoints require `from` and `to` query parameters specifying a date range (format: `YYYY-MM-DD`).

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../statistics/staffing-hours` | Staffing hours summary |
| GET | `.../statistics/staffing-hours/employees` | Per-employee staffing detail |
| GET | `.../statistics/financials` | Financial overview |
| GET | `.../statistics/occupancy` | Occupancy statistics |
| GET | `.../statistics/age-distribution` | Age distribution |
| GET | `.../statistics/contract-properties` | Contract property distribution |
| GET | `.../statistics/funding` | Funding statistics |

## Users

Global user management endpoints.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List users |
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users/{id}` | Get user |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Delete user |
| GET | `/api/v1/users/{id}/memberships` | Get user's organization memberships |
| POST | `/api/v1/users/{id}/organizations` | Add user to organization |
| PUT | `/api/v1/users/{id}/organizations/{orgId}` | Update user's role in organization |
| DELETE | `/api/v1/users/{id}/organizations/{orgId}` | Remove user from organization |
| PUT | `/api/v1/users/{id}/password` | Reset user's password (admin) |
| PUT | `/api/v1/users/{id}/superadmin` | Set superadmin status |

### Organization Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/organizations/{orgId}/users` | List users in an organization |

## Pagination

List endpoints support pagination via query parameters:

```bash
curl "http://localhost:8080/api/v1/organizations?page=1&limit=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

Response:

```json
{
  "data": [],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

## Error Responses

Errors are returned with the appropriate HTTP status code and a JSON body:

```json
{
  "error": "Description of the error"
}
```

| Status | Meaning |
|--------|---------|
| 400 | Bad Request -- Invalid input or missing required parameters |
| 401 | Unauthorized -- Missing or invalid authentication token |
| 403 | Forbidden -- Insufficient permissions for the requested action |
| 404 | Not Found -- The requested resource does not exist |
| 500 | Internal Server Error -- An unexpected error occurred |
