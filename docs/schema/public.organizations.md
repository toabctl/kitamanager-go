# public.organizations

## Description

## Columns

| Name       | Type                     | Default                                   | Nullable | Children                                                                                                                                                                                            | Parents | Comment |
| ---------- | ------------------------ | ----------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------- |
| id         | bigint                   | nextval('organizations_id_seq'::regclass) | false    | [public.group_organizations](public.group_organizations.md) [public.user_organizations](public.user_organizations.md) [public.employees](public.employees.md) [public.children](public.children.md) |         |         |
| name       | varchar(255)             |                                           | false    |                                                                                                                                                                                                     |         |         |
| active     | boolean                  | true                                      | true     |                                                                                                                                                                                                     |         |         |
| created_at | timestamp with time zone |                                           | true     |                                                                                                                                                                                                     |         |         |
| created_by | varchar(255)             |                                           | true     |                                                                                                                                                                                                     |         |         |
| updated_at | timestamp with time zone |                                           | true     |                                                                                                                                                                                                     |         |         |

## Constraints

| Name               | Type        | Definition       |
| ------------------ | ----------- | ---------------- |
| organizations_pkey | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name               | Definition                                                                      |
| ------------------ | ------------------------------------------------------------------------------- |
| organizations_pkey | CREATE UNIQUE INDEX organizations_pkey ON public.organizations USING btree (id) |

## Relations

```mermaid
erDiagram

"public.group_organizations" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.group_organizations" }o--|| "public.groups" : "FOREIGN KEY (group_id) REFERENCES groups(id)"
"public.user_organizations" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.user_organizations" }o--|| "public.users" : "FOREIGN KEY (user_id) REFERENCES users(id)"
"public.employees" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.employee_contracts" }o--|| "public.employees" : "FOREIGN KEY (employee_id) REFERENCES employees(id)"
"public.children" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.child_contracts" }o--|| "public.children" : "FOREIGN KEY (child_id) REFERENCES children(id)"

"public.organizations" {
  bigint id
  varchar_255_ name
  boolean active
  timestamp_with_time_zone created_at
  varchar_255_ created_by
  timestamp_with_time_zone updated_at
}
"public.group_organizations" {
  bigint group_id FK
  bigint organization_id FK
}
"public.groups" {
  bigint id
  varchar_255_ name
  boolean active
  timestamp_with_time_zone created_at
  varchar_255_ created_by
  timestamp_with_time_zone updated_at
}
"public.user_organizations" {
  bigint user_id FK
  bigint organization_id FK
}
"public.users" {
  bigint id
  varchar_255_ name
  varchar_255_ email
  varchar_255_ password
  boolean active
  timestamp_with_time_zone created_at
  varchar_255_ created_by
  timestamp_with_time_zone updated_at
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
  timestamp_with_time_zone created_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
