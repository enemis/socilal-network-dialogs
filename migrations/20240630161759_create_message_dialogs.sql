-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE dialog_type AS ENUM ('direct', 'private', 'published');
CREATE TABLE dialogs (
     "id" uuid DEFAULT uuid_generate_v1(),
     "name" varchar(128),
     "type" dialog_type NOT NULL default 'direct',
     "created_at" timestamp NOT NULL,
     "deleted_at" timestamp,
     PRIMARY KEY (id));
CREATE INDEX IF NOT EXISTS type_index ON dialogs (name varchar_pattern_ops);
CREATE TABLE dialog_participants (
    "dialog_id" uuid,
    "user_id" uuid NOT NULL,
--     "created_at" timestamp NOT NULL,
    "deleted_at" timestamp,
    UNIQUE (dialog_id, user_id),
    CONSTRAINT fk_dialogs_participiants FOREIGN KEY(dialog_id)
        REFERENCES dialogs(id)
);
CREATE INDEX idx_dialog_participants_dialog_id
    ON dialog_participants(dialog_id);
CREATE INDEX idx_dialog_participants_user_id
    ON dialog_participants(user_id);
CREATE TABLE messages (
      "id" uuid DEFAULT uuid_generate_v1(),
      "user_id" varchar(128) NOT NULL,
      "dialog_id" uuid NOT NULL,
      "message" text NOT NULL,
      "created_at" timestamp NOT NULL,
      "deleted_at" timestamp,
      CONSTRAINT fk_dialogs FOREIGN KEY(dialog_id)
          REFERENCES dialogs(id)
);
-- +goose StatementEnd
EXPLAIN SELECT m.* FROM dialogs d INNER JOIN messages m ON m.dialog_id = d.id
        WHERE d.id='061048b0-3ef0-4e2e-9b44-ec55c097f2b5'
          AND d.deleted_at IS NULL AND d.type = 'direct'
        ORDER BY m.created_at;


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
