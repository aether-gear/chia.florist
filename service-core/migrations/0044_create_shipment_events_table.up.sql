CREATE TABLE shipment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL,
    status TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_shipment
        FOREIGN KEY(shipment_id)
        REFERENCES shipments(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_shipment_events_shipment_id ON shipment_events(shipment_id);
