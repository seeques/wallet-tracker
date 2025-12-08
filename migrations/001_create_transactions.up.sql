CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    val NUMERIC NOT NULL,
    tmstp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);