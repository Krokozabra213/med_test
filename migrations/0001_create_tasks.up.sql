CREATE TABLE IF NOT EXISTS tasks (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'new',
	scheduled_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 

	CONSTRAINT tasks_status_check
        CHECK (status IN ('new', 'completed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
