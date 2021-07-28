create table reviews (
  review_id uuid primary key unique not null default (uuid_generate_v4()),
  from_user_id uuid references users on delete set null on update cascade,
  to_user_id uuid references users on delete set null on update cascade,
  text text,
  created_at timestamptz not null default (now()),
  updated_at timestamptz not null default (now())
);
