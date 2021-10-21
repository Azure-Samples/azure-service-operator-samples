package azvotes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

type Client struct {
	db *sql.DB
}

func NewClient(server string, user string, password string, port int, database string) (*Client, error) {
	// postgres://asoadmin:{your_password}@asodemo-postgres.postgres.database.azure.com/postgres?sslmode=require
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", user, password, server, port, database)
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to database")
	}

	return &Client{db}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	subctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return c.db.PingContext(subctx)
}

func (c *Client) Init(ctx context.Context) error {
	tableQuery := `create table if not exists votes (
		id       serial primary key,
		vote     varchar(50)
	);`

	_, err := c.db.ExecContext(ctx, tableQuery)
	if err != nil {
		return errors.Wrapf(err, "Creating table votes")
	}

	return nil
}

// CreateVote inserts an employee record
func (c *Client) CreateVote(val string) (int64, error) {
	ctx := context.Background()
	var err error

	if c.db == nil {
		err = errors.New("CreateVote: db is null")
		return -1, err
	}

	stmt, err := c.db.Prepare("INSERT INTO votes (vote) VALUES ($1) RETURNING id;")
	if err != nil {
		return -1, errors.Wrapf(err, "preparing insert staement")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, val)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, errors.Wrapf(err, "executing insert to votes table")
	}

	return newID, nil
}

func (c *Client) CountVote(vote string) (int, error) {
	ctx := context.Background()

	var count int
	stmt, err := c.db.Prepare("SELECT COUNT (id) FROM votes WHERE vote = $1;")
	if err != nil {
		return -1, errors.Wrapf(err, "preparing count votes statement")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		vote,
	)

	err = row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (c *Client) DeleteVotes() (int64, error) {
	ctx := context.Background()

	result, err := c.db.ExecContext(ctx, "DELETE FROM votes;")
	if err != nil {
		return -1, errors.Wrapf(err, "failed to delete votes")
	}

	return result.RowsAffected()
}
