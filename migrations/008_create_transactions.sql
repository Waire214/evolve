CREATE TABLE transactions (
     id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
     link_id UUID NULL,
     source_user_id UUID NULL,
     source_account_id UUID NULL,
     source_account_number VARCHAR(50) NOT NULL,
     target_user_id UUID NULL,
     target_account_id UUID NULL,
     target_account_number VARCHAR(50) NOT NULL,
     transaction_type VARCHAR(20) NOT NULL,
     source_transaction_amount NUMERIC(10, 2) NOT NULL,
     target_transaction_amount NUMERIC(10, 2) NOT NULL,
     source_balance_after_transaction NUMERIC(10, 2) NOT NULL,
     target_balance_after_transaction NUMERIC(10, 2) NOT NULL,
     status VARCHAR(20) NOT NULL,
     transaction_date TIMESTAMP NOT NULL,
     updated_at timestamp with time zone DEFAULT NOW() NULL,
     CONSTRAINT fk_transactions_source_user FOREIGN KEY (source_user_id) REFERENCES users (id),
     CONSTRAINT fk_transactions_source_account FOREIGN KEY ("source_account_id") REFERENCES accounts (id)
);

CREATE INDEX idx_transactions_id ON transactions (id);
CREATE INDEX idx_transactions_source_user_id ON transactions (source_user_id);