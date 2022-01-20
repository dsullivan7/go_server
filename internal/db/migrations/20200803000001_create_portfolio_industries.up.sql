create table portfolio_industries (
  portfolio_industry_id uuid primary key unique not null default (uuid_generate_v4()),
  portfolio_id uuid references portfolios on delete cascade on update cascade,
  industry_id uuid references industries on delete cascade on update cascade,
  created_at timestamptz not null default (now()),
  updated_at timestamptz not null default (now())
);
