CREATE TABLE asynq_job (
    app_uuid UUID NOT NULL,
    asynq_uuid UUID UNIQUE NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    status VARCHAR(255) NULL,
    message_text TEXT NULL,
    order_number BIGINT NULL,
    payload JSONB NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- PRIMARY KEY (app_uuid)
);

-- Create an index for the order_number column in the asynq_job table
CREATE INDEX idx_async_job_order_number ON asynq_job (order_number);