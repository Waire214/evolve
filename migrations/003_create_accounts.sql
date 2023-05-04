CREATE TABLE IF NOT EXISTS accounts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL,
    account_type VARCHAR(200) NOT NULL,
    bank_code VARCHAR(200) NOT NULL,
    account_number VARCHAR(200) NOT NULL UNIQUE,
    balance NUMERIC(10, 2) DEFAULT NULL,
    status VARCHAR(100) DEFAULT 'active' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at timestamp with time zone DEFAULT NOW() NULL,
    CONSTRAINT fk_accounts_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_accounts_id ON accounts (id);
CREATE INDEX idx_accounts_account_number ON accounts (account_number);
CREATE INDEX idx_accounts_user_id ON accounts (user_id);
CREATE INDEX idx_accounts_created_at ON accounts (created_at);