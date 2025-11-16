-- +goose Up
ALTER TABLE users
ADD CONSTRAINT fk_users_team
FOREIGN KEY (team_name) REFERENCES teams(team_name) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT fk_users_team;