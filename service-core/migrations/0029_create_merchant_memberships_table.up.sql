CREATE TABLE merchant_memberships (
    id UUID PRIMARY KEY,

    merchant_id UUID NOT NULL,
    account_id UUID NOT NULL,

    role_id UUID NOT NULL,

    created_by UUID NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_merchant
        FOREIGN KEY (merchant_id)
        REFERENCES merchants(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id),
        
    CONSTRAINT uq_merchant_account
        UNIQUE (account_id)
);
