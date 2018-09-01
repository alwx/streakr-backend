CREATE TABLE IF NOT EXISTS campaigns (
  id UUID NOT NULL,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  amount integer NOT NULL,
  min_price integer NOT NULL,
  exp_date date,
  prize_description text NOT NULL,
  PRIMARY KEY (id)
)