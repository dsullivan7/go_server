create table user_industries (
  user_industry_id uuid primary key unique not null default (uuid_generate_v4()),
  user_id uuid references users on delete cascade on update cascade,
  industry_id uuid references industries on delete cascade on update cascade,
  created_at timestamptz not null default (now()),
  updated_at timestamptz not null default (now())
);
