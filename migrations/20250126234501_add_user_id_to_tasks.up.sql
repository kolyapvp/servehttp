ALTER TABLE tasks
    ADD COLUMN user_id INTEGER,
    ADD CONSTRAINT fk_tasks_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;