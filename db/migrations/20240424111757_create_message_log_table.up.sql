CREATE TABLE message_log (
    asynq_uuid UUID NOT NULL,
    data_log TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);