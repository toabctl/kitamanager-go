# public.employees

## Description

## Columns

| Name            | Type                     | Default                               | Nullable | Children                                                  | Parents                                         | Comment |
| --------------- | ------------------------ | ------------------------------------- | -------- | --------------------------------------------------------- | ----------------------------------------------- | ------- |
| id              | bigint                   | nextval('employees_id_seq'::regclass) | false    | [public.employee_contracts](public.employee_contracts.md) |                                                 |         |
| organization_id | bigint                   |                                       | false    |                                                           | [public.organizations](public.organizations.md) |         |
| first_name      | varchar(255)             |                                       | false    |                                                           |                                                 |         |
| last_name       | varchar(255)             |                                       | false    |                                                           |                                                 |         |
| birthdate       | date                     |                                       | false    |                                                           |                                                 |         |
| created_at      | timestamp with time zone |                                       | true     |                                                           |                                                 |         |
| updated_at      | timestamp with time zone |                                       | true     |                                                           |                                                 |         |

## Constraints

| Name                               | Type        | Definition                                                 |
| ---------------------------------- | ----------- | ---------------------------------------------------------- |
| employees_birthdate_not_null       | n           | NOT NULL birthdate                                         |
| employees_first_name_not_null      | n           | NOT NULL first_name                                        |
| employees_id_not_null              | n           | NOT NULL id                                                |
| employees_last_name_not_null       | n           | NOT NULL last_name                                         |
| employees_organization_id_not_null | n           | NOT NULL organization_id                                   |
| fk_employees_organization          | FOREIGN KEY | FOREIGN KEY (organization_id) REFERENCES organizations(id) |
| employees_pkey                     | PRIMARY KEY | PRIMARY KEY (id)                                           |

## Indexes

| Name                          | Definition                                                                                   |
| ----------------------------- | -------------------------------------------------------------------------------------------- |
| employees_pkey                | CREATE UNIQUE INDEX employees_pkey ON public.employees USING btree (id)                      |
| idx_employees_organization_id | CREATE INDEX idx_employees_organization_id ON public.employees USING btree (organization_id) |

## Relations

```mermaid
erDiagram

"public.employee_contracts" }o--|| "public.employees" : "FOREIGN KEY (employee_id) REFERENCES employees(id)"
"public.employees" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.groups" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.children" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.government_fundings" : "FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id)"

"public.employees" {
  bigint id
  bigint organization_id FK
  varchar_255_ first_name
  varchar_255_ last_name
  date birthdate
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.employee_contracts" {
  bigint id
  bigint employee_id FK
  date from_date
  date to_date
  varchar_255_ position
  numeric weekly_hours
  bigint salary
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.organizations" {
  bigint id
  varchar_255_ name
  boolean active
  bigint payplan_id FK
  timestamp_with_time_zone created_at
  varchar_255_ created_by
  timestamp_with_time_zone updated_at
  bigint funding_id
  bigint government_funding_id FK
}
"public.groups" {
  bigint id
  varchar_255_ name
  bigint organization_id FK
  boolean is_default
  boolean active
  timestamp_with_time_zone created_at
  varchar_255_ created_by
  timestamp_with_time_zone updated_at
}
"public.children" {
  bigint id
  bigint organization_id FK
  varchar_255_ first_name
  varchar_255_ last_name
  date birthdate
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.payplans" {
  bigint id
  varchar_255_ name
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.government_fundings" {
  bigint id
  varchar_255_ name
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
