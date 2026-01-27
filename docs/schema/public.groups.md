# public.groups

## Description

## Columns

| Name            | Type                     | Default                            | Nullable | Children                                    | Parents                                         | Comment |
| --------------- | ------------------------ | ---------------------------------- | -------- | ------------------------------------------- | ----------------------------------------------- | ------- |
| id              | bigint                   | nextval('groups_id_seq'::regclass) | false    | [public.user_groups](public.user_groups.md) |                                                 |         |
| name            | varchar(255)             |                                    | false    |                                             |                                                 |         |
| organization_id | bigint                   |                                    | false    |                                             | [public.organizations](public.organizations.md) |         |
| is_default      | boolean                  | false                              | true     |                                             |                                                 |         |
| active          | boolean                  | true                               | true     |                                             |                                                 |         |
| created_at      | timestamp with time zone |                                    | true     |                                             |                                                 |         |
| created_by      | varchar(255)             |                                    | true     |                                             |                                                 |         |
| updated_at      | timestamp with time zone |                                    | true     |                                             |                                                 |         |

## Constraints

| Name                            | Type        | Definition                                                 |
| ------------------------------- | ----------- | ---------------------------------------------------------- |
| groups_id_not_null              | n           | NOT NULL id                                                |
| groups_name_not_null            | n           | NOT NULL name                                              |
| groups_organization_id_not_null | n           | NOT NULL organization_id                                   |
| fk_organizations_groups         | FOREIGN KEY | FOREIGN KEY (organization_id) REFERENCES organizations(id) |
| groups_pkey                     | PRIMARY KEY | PRIMARY KEY (id)                                           |

## Indexes

| Name        | Definition                                                        |
| ----------- | ----------------------------------------------------------------- |
| groups_pkey | CREATE UNIQUE INDEX groups_pkey ON public.groups USING btree (id) |

## Relations

```mermaid
erDiagram

"public.user_groups" }o--|| "public.groups" : "FOREIGN KEY (group_id) REFERENCES groups(id)"
"public.user_groups" }o--|| "public.users" : "FOREIGN KEY (user_id) REFERENCES users(id)"
"public.groups" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.employees" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.children" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.government_fundings" : "FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id)"

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
"public.users" {
  bigint id
  varchar_255_ name
  varchar_255_ email
  varchar_255_ password
  boolean active
  boolean is_superadmin
  timestamp_with_time_zone last_login
  timestamp_with_time_zone created_at
  varchar_255_ created_by
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
"public.employees" {
  bigint id
  bigint organization_id FK
  varchar_255_ first_name
  varchar_255_ last_name
  date birthdate
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
