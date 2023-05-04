CREATE TABLE cashCache (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    cached_amount NUMERIC(10, 2) NOT NULL,
    cache_duration INTERVAL NOT NULL,
    cached_at TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NULL,
    CONSTRAINT fk_cashCache_account
      FOREIGN KEY (account_number)
          REFERENCES accounts (account_number)
);

CREATE INDEX idx_cashCache_id ON cashCache (id);
CREATE INDEX idx_cashCache_account ON cashCache (account_number);