# public.children

## Description

## Columns

| Name            | Type                     | Default                              | Nullable | Children                                            | Parents                                         | Comment |
| --------------- | ------------------------ | ------------------------------------ | -------- | --------------------------------------------------- | ----------------------------------------------- | ------- |
| id              | bigint                   | nextval('children_id_seq'::regclass) | false    | [public.child_contracts](public.child_contracts.md) |                                                 |         |
| organization_id | bigint                   |                                      | false    |                                                     | [public.organizations](public.organizations.md) |         |
| first_name      | varchar(255)             |                                      | false    |                                                     |                                                 |         |
| last_name       | varchar(255)             |                                      | false    |                                                     |                                                 |         |
| birthdate       | date                     |                                      | false    |                                                     |                                                 |         |
| created_at      | timestamp with time zone |                                      | true     |                                                     |                                                 |         |
| updated_at      | timestamp with time zone |                                      | true     |                                                     |                                                 |         |

## Constraints

| Name                              | Type        | Definition                                                 |
| --------------------------------- | ----------- | ---------------------------------------------------------- |
| children_birthdate_not_null       | n           | NOT NULL birthdate                                         |
| children_first_name_not_null      | n           | NOT NULL first_name                                        |
| children_id_not_null              | n           | NOT NULL id                                                |
| children_last_name_not_null       | n           | NOT NULL last_name                                         |
| children_organization_id_not_null | n           | NOT NULL organization_id                                   |
| fk_children_organization          | FOREIGN KEY | FOREIGN KEY (organization_id) REFERENCES organizations(id) |
| children_pkey                     | PRIMARY KEY | PRIMARY KEY (id)                                           |

## Indexes

| Name                         | Definition                                                                                 |
| ---------------------------- | ------------------------------------------------------------------------------------------ |
| children_pkey                | CREATE UNIQUE INDEX children_pkey ON public.children USING btree (id)                      |
| idx_children_organization_id | CREATE INDEX idx_children_organization_id ON public.children USING btree (organization_id) |

## Relations

```mermaid
erDiagram

"public.child_contracts" }o--|| "public.children" : "FOREIGN KEY (child_id) REFERENCES children(id)"
"public.children" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.groups" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.employees" }o--|| "public.organizations" : "FOREIGN KEY (organization_id) REFERENCES organizations(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.payplans" : "FOREIGN KEY (payplan_id) REFERENCES payplans(id)"
"public.organizations" }o--o| "public.government_fundings" : "FOREIGN KEY (government_funding_id) REFERENCES government_fundings(id)"

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
"public.employees" {
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
