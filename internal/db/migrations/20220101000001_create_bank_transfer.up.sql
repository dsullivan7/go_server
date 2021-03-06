create table bank_transfers (
  bank_transfer_id uuid primary key unique not null default (uuid_generate_v4()),
  user_id uuid references users on delete set null on update cascade,
  status text not null,
  amount integer not null,
  alpaca_transfer_id text,
  plaid_transfer_id text,
  created_at timestamptz not null default (now()),
  updated_at timestamptz not null default (now())
);
