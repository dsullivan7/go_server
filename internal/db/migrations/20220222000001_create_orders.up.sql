create table orders (
  order_id uuid primary key unique not null default (uuid_generate_v4()),
  portfolio_id uuid references portfolios on delete set null on update cascade,
  user_id uuid references users on delete set null on update cascade,
  side text,
  amount numeric,
  alpaca_order_id text,
  created_at timestamptz not null default (now()),
  updated_at timestamptz not null default (now())
);
