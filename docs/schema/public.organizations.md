# public.organizations

## Description

## Columns

| Name                  | Type                     | Default                                   | Nullable | Children                                                                                                        | Parents                                                     | Comment |
| --------------------- | ------------------------ | ----------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------- | ------- |
| id                    | bigint                   | nextval('organizations_id_seq'::regclass) | false    | [public.groups](public.groups.md) [public.employees](public.employees.md) [public.children](public.children.md) |                                                             |         |
| name                  | varchar(255)             |                                           | false    |                                                                                                                 |                                                             |         |
| active                | boolean                  | true                                      | true     |                                                                                                                 |                                                             |         |
| payplan_id            | bigint                   |                                           | true     |                                                                                                                 | [public.payplans](public.payplans.md)                       |         |
| created_at            | timestamp with time zone |                                           | true     |                                                                                                                 |                                                             |         |
| created_by            | varchar(255)             |                                           | true     |                                                                                                                 |                                                             |         |
| updated_at            | timestamp with time zone |                                           | true     |                                                                                                                 |                                                             |         |
| funding_id            | bigint                   |                                           | true     |                                                                                                                 |                                                             |         |
| government_funding_id | bigint                   |                                           | true     |                                                                                                                 | [public.government_fundings](public.government_fundings.md) |         |

## Constraints

| Name                                | Type        | Definition                                                             |
| ----------------------------------- | ----------- | ---------------------------------------------------------------------- |
| organizations_id_not_null           | n           | NOT NULL id                                                            |
| organizations_name_not_null         | n           | NOT NULL name                                                          |
| fk_organizations_funding            | FOREIGN KEY | FOREIGN KEY (payplan_id) REFERENCES payplans(id)                       |
| fk_organizations_payplan            | FOREIGN KEY | FOREIGN KEY (payplan_id) REFERENCES payplans(id)                       |
| organizations_pkey                  | PRIMARY KEY | PRIMARY KEY (id)                                                       |
| fk_organizations_government_funding | FOREIGN KEY | FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id) |

## Indexes

| Name               | Definition                                                                      |
| ------------------ | ------------------------------------------------------------------------------- |
| organizations_pkey | CREATE UNIQUE INDEX organizations_pkey ON public.organizations USING btree (id) |

## Relations

```mermaid
erDiagram

"public.groups" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.user_groups" }o--|| "public.groups" : "FOREIGN KEY (group_id) REFERENCES groups(id)"
"public.employees" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.employee_contracts" }o--|| "public.employees" : "FOREIGN KEY (employee_id) REFERENCES employees(id)"
"public.children" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.child_contracts" }o--|| "public.children" : "FOREIGN KEY (child_id) REFERENCES children(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.payplan_periods" }o--|| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id) ON DELETE CASCADE"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.government_fundings" : "FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id)"
"public.government_funding_periods" }o--|| "public.government_fundings" : "FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id) ON DELETE CASCADE"

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
"public.user_groups" {
  bigint user_id FK
  bigint group_id FK
  varchar_50_ role
  timestamp_with_time_zone created_at
  varchar_255_ created_by
}
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
"public.children" {
  bigint id
  bigint organization_id FK
  varchar_255_ first_name
  varchar_255_ last_name
  date birthdate
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.child_contracts" {
  bigint id
  bigint child_id FK
  date from_date
  date to_date
  numeric care_hours_per_week
  bigint group_id
  boolean meals_included
  varchar_1000_ special_needs
  text__ attributes
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.payplans" {
  bigint id
  varchar_255_ name
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.payplan_periods" {
  bigint id
  bigint payplan_id FK
  date from_date
  date to_date
  varchar_1000_ comment
  timestamp_with_time_zone created_at
}
"public.government_fundings" {
  bigint id
  varchar_255_ name
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
}
"public.government_funding_periods" {
  bigint id
  bigint government_funding_id FK
  date from_date
  date to_date
  varchar_1000_ comment
  timestamp_with_time_zone created_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
