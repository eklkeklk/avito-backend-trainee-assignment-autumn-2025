package postgres

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/team"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type TeamRepositoryPostgres struct {
	db *sql.DB
}

func NewTeamRepositoryPostgres(db *sql.DB) *TeamRepositoryPostgres {
	return &TeamRepositoryPostgres{
		db: db,
	}
}

func (t TeamRepositoryPostgres) Create(ctx context.Context, team *team.Team) (*team.Team, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			return
		}
	}(tx)
	queryTeam := `
        INSERT INTO teams (team_name) 
        VALUES ($1)
        ON CONFLICT (team_name) DO NOTHING
    `
	_, err = tx.ExecContext(ctx, queryTeam, team.Name)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	err = t.insertTeamMembers(ctx, tx, team.Name, team.Members)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t.GetByName(ctx, team.Name)
}

func (t TeamRepositoryPostgres) GetByName(ctx context.Context, name string) (*team.Team, error) {
	query := `
	SELECT t.team_name, u.user_id, u.username, u.is_active
	FROM teams t
	LEFT JOIN users u ON t.team_name = u.team_name
	WHERE t.team_name = $1
	`
	rows, err := t.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	var foundTeam *team.Team
	for rows.Next() {
		var (
			teamName string
			userId   string
			username string
			isActive bool
		)
		err := rows.Scan(&teamName, &userId, &username, &isActive)
		if err != nil {
			return nil, fmt.Errorf("internal error")
		}
		if foundTeam == nil {
			foundTeam = &team.Team{
				Name:    teamName,
				Members: make([]*user.TeamMember, 0),
			}
		}
		if userId != "" {
			member := &user.TeamMember{
				UserId:   userId,
				Username: username,
				IsActive: isActive,
			}
			foundTeam.Members = append(foundTeam.Members, member)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("internal error")
	}
	if foundTeam == nil {
		return nil, fmt.Errorf("resource not found")
	}

	return foundTeam, nil
}

func (t TeamRepositoryPostgres) insertTeamMembers(ctx context.Context, tx *sql.Tx, team string, members []*user.TeamMember) error {
	allUsersExist := true
	if len(members) > 0 {
		checkUsersQuery := `
            SELECT COUNT(*) = $1
            FROM users 
            WHERE user_id = ANY($2) AND team_name = $3
        `
		userIDs := make([]string, len(members))
		for i, member := range members {
			userIDs[i] = member.UserId
		}
		err := tx.QueryRowContext(
			ctx,
			checkUsersQuery,
			len(members),
			pq.Array(userIDs),
			team,
		).Scan(&allUsersExist)

		if err != nil {
			return fmt.Errorf("internal error")
		}
	}
	if allUsersExist {
		return fmt.Errorf("already exists")
	}
	if len(members) > 0 {
		queryUsers := `
            INSERT INTO users (user_id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) 
            DO UPDATE SET 
                username = EXCLUDED.username,
                team_name = EXCLUDED.team_name,
                is_active = EXCLUDED.is_active
        `
		for _, member := range members {
			_, err := tx.ExecContext(
				ctx,
				queryUsers,
				member.UserId,
				member.Username,
				team,
				member.IsActive,
			)
			if err != nil {
				return fmt.Errorf("internal error")
			}
		}
	}

	return nil
}
