DROP TABLE IF EXISTS task_recurrence_rules;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_recurrence_rule;
ALTER TABLE tasks DROP COLUMN IF EXISTS recurrence_rule_id;
