CREATE TABLE jixa_master (
    id SERIAL PRIMARY KEY,
    item_name INT,
    item_head STRING(30),
    txn_type STRING(20),
    cost DECIMAL DEFAULT,
    created_at TIMESTAMPTZ DEFAULT now(),
);