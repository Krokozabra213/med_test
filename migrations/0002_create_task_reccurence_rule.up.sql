CREATE TABLE task_recurrence_rules (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    recurrence_type TEXT NOT NULL,
    settings JSONB NOT NULL,
    scheduled_time TIME NOT NULL,
    timezone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE tasks 
    ADD COLUMN recurrence_rule_id BIGINT;

ALTER TABLE tasks 
    ADD CONSTRAINT fk_tasks_recurrence_rule 
    FOREIGN KEY (recurrence_rule_id) REFERENCES task_recurrence_rules(id) 
    ON DELETE SET NULL;