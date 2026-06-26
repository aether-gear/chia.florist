CREATE TABLE staff_memberships (
    id UUID PRIMARY KEY,

    staff_id UUID NOT NULL,
    account_id UUID NOT NULL,

    role_id UUID NOT NULL,

    created_by UUID NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_staff
        FOREIGN KEY (staff_id)
        REFERENCES staff(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id),
        
    CONSTRAINT uq_staff_account
        UNIQUE (account_id)
);
