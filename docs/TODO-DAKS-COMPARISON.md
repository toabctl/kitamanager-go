# KitaManager Feature Comparison with DAKS Kalkulator

This document compares KitaManager features with the DAKS Kalkulator software and identifies missing features.

## Current KitaManager Features

### ✅ Already Implemented
- [x] Organizations management (CRUD)
- [x] Groups management (CRUD, per organization)
- [x] Users management with RBAC (superadmin, admin, manager, member)
- [x] Children management with contracts
  - Name, birthdate
  - Care hours per week
  - Group assignment
  - Meals included flag
  - Special needs notes
  - Attributes (e.g., ganztag, ndh)
  - Contract periods (from/to dates)
- [x] Employees management with contracts
  - Name, birthdate
  - Position
  - Weekly hours
  - Salary (in cents)
  - Contract periods (from/to dates)
- [x] Government Funding (Landeszuschüsse Berlin)
  - Funding periods with age ranges
  - Payment amounts per category
  - FTE requirements per category
  - Assignment to organizations

---

## Missing Features (Compared to DAKS Kalkulator)

### 🔴 High Priority - Core Business Features

#### 1. Care Type/Extent (Umfang)
DAKS has predefined care types for children:
- [ ] `halbtags ohne Mittag` (half-day without lunch)
- [ ] `halbtags` (half-day)
- [ ] `teilzeit` (part-time)
- [ ] `ganztags` (full-day)
- [ ] `ganztags erweitert` (extended full-day)

**Current state:** KitaManager uses `Attributes` array (e.g., `["ganztags", "ndh"]`) to store care type and extras.
**Recommendation:** The Attributes approach is flexible and can encode both care type and surcharges.

#### 2. Children Surcharges/Extras (Zuschlag)
DAKS tracks special surcharges for children:
- [ ] `Integration A` - Integration support level A
- [ ] `Integration B` - Integration support level B
- [ ] `ndH` - Non-German speaking household (nicht-deutsche Herkunftssprache)
- [ ] `QM` - Quality management/enhancement

**Current state:** KitaManager uses `Attributes` array on `ChildContract` to store surcharges (e.g., `["integration_a", "ndh"]`).
**Recommendation:** The Attributes approach covers surcharges. Consider validation to ensure valid attribute values.

#### 3. Employee Contract Extensions
DAKS has additional employee contract fields:
- [ ] `Gehalt voll` - Full salary amount
- [ ] `Stunden voll` - Full-time hours (contractual)
- [ ] `Stunden real` - Real/actual working hours
- [ ] `Weihnachtsgeld` - Christmas bonus (amount or percentage)
- [ ] `Arbeitgeberanteil` - Employer contribution type (Normal vs Minijob)
- [ ] `Anrechenbar auf Fachpersonalschlüssel` - Counts towards qualified staff ratio

**Current state:** KitaManager has `salary` and `weekly_hours` but missing these Berlin-specific fields.

#### 4. Dashboard / Financial Overview (Übersicht)
DAKS has a monthly overview showing:
- [ ] Monthly income totals (Einnahmen)
- [ ] Monthly expense totals (Ausgaben)
- [ ] Monthly coverage/balance (Deckung)
- [ ] Monthly staff hours (Soll/Ist/Deckung)

**Current state:** No dashboard in KitaManager.

### 🟡 Medium Priority - Reporting & Calculations

#### 5. Staff Hours Tracking (Personalstunden)
Monthly tracking of staff hours:
- [ ] Total hours (Gesamt)
- [ ] Pedagogical hours actual (päd. Stunden Ist)
- [ ] Pedagogical hours required (päd. Stunden Soll)
- [ ] Coverage calculation (Deckung)

**Current state:** Not implemented.

#### 6. Staff Ratio Calculation (Personalschlüssel)
Calculate required staff based on:
- [ ] Children count by age group
- [ ] Care type distribution
- [ ] Integration children count
- [ ] Compare with actual qualified staff

**Current state:** Government funding entries exist but no calculation logic.

#### 7. Occupancy Overview (Belegung)
- [ ] Children count by care type
- [ ] Children count by age group
- [ ] Monthly occupancy changes

#### 8. Income Overview (Übersicht Einnahmen)
- [ ] Government funding per child
- [ ] Additional fees per child
- [ ] Total income calculation

### 🟢 Lower Priority - Budget & Configuration

#### 9. Budget Module (Haushalt)
Comprehensive budget tracking:

**Income categories (Einnahmen):**
- [ ] Senat und TKBG (Government funding)
- [ ] Zusatzbeitrag (Additional parental contribution)
- [ ] Vereinsbeitrag (Association membership fee)
- [ ] Erstattungen U1 (U1 reimbursements)
- [ ] Sonstige Einnahmen (Other income)

**Expense categories (Ausgaben):**
- [ ] Gehälter (Salaries)
- [ ] Aufwand, Honorare (Fees and honorariums)
- [ ] Fortbildung, Qualität (Training, quality)
- [ ] Verwaltung (Administration)
- [ ] Wirtschaft, Betreuung (Operations, care supplies)
- [ ] Miete, Energie (Rent, utilities)
- [ ] Sonstiges (Other expenses)

#### 10. Basic Data Configuration (Grunddaten)
Organization-level settings:
- [ ] `durchschnittl. Zusatzbeitrag pro Kind` - Average additional fee per child
- [ ] `Arbeitgebernebenkosten Normal` - Normal employer side costs (%)
- [ ] `Arbeitgebernebenkosten Minijob` - Minijob employer side costs (%)
- [ ] Custom income category names
- [ ] Custom expense category names

### 🔵 Nice to Have - Export & UX

#### 11. Export Features
- [ ] Export to PDF (lists, reports)
- [ ] Export to Excel (data tables)
- [ ] Download manual/help document

#### 12. Inline Help System
- [ ] Contextual help buttons per section
- [ ] Help modal with explanations

---

## Implementation Recommendations

### Phase 1: Core Data Model Extensions
1. Add `care_type` enum to ChildContract
2. Add structured surcharge fields to ChildContract
3. Add employee contract extensions (Christmas bonus, employer type, qualified flag)

### Phase 2: Calculations & Dashboard
1. Implement government funding calculation based on children data
2. Implement staff ratio calculation
3. Create dashboard with overview

### Phase 3: Budget & Reporting
1. Add budget module with income/expense tracking
2. Add organization settings (Grunddaten)
3. Implement staff hours tracking

### Phase 4: Polish
1. Add PDF/Excel export
2. Add inline help system
3. Improve UX based on DAKS patterns

---

## Data Model Mapping

### Child Care Types → Government Funding
| DAKS Care Type | Hours/Week | Government Funding Category |
|----------------|------------|----------------------------|
| halbtags ohne Mittag | ~20h | TBD |
| halbtags | ~25h | TBD |
| teilzeit | ~35h | TBD |
| ganztags | ~45h | TBD |
| ganztags erweitert | ~50h+ | TBD |

### Employee Fields Mapping
| DAKS Field | KitaManager Field | Status |
|------------|-------------------|--------|
| Vorname | first_name | ✅ |
| Nachname | last_name | ✅ |
| Gehalt voll | salary | ✅ (partial) |
| Stunden voll | weekly_hours | ✅ |
| Stunden real | - | ❌ Missing |
| Weihnachtsgeld | - | ❌ Missing |
| Vertragsanfang | from | ✅ |
| Vertragsende | to | ✅ |
| Arbeitgeberanteil | - | ❌ Missing |
| Fachpersonal | - | ❌ Missing |

### Child Fields Mapping
| DAKS Field | KitaManager Field | Status |
|------------|-------------------|--------|
| Vorname | first_name | ✅ |
| Nachname | last_name | ✅ |
| Geburtstag | birthdate | ✅ |
| Umfang | attributes (e.g., "ganztags") | ✅ Via Attributes |
| Zuschlag | attributes (e.g., "ndh", "integration_a") | ✅ Via Attributes |
| Vertragsanfang | from | ✅ |
| Vertragsende | to | ✅ |

---

*Generated: 2026-01-27*
*Based on investigation of DAKS Kalkulator (https://kalkulator.daks-berlin.de)*
