package azvotes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var tableQuery = `CREATE TABLE VoteSchema.Votes (
	Id       INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
	Vote     NVARCHAR(50)
);`

type Client struct {
	db *sql.DB
}

func NewClient(server string, user string, password string, port int, database string) (*Client, error) {
	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}
	return &Client{db}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	subctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return c.db.PingContext(subctx)
}

func (c *Client) Init() error {

	stmt, err := c.db.Prepare("CREATE SCHEMA VoteSchema;")
	if err != nil {
		return err
	}

	_, _ = stmt.Exec()

	stmt, err = c.db.Prepare(tableQuery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		if !strings.Contains(err.Error(), "already an object") {
			return err
		}
	}

	return nil
}

// CreateVote inserts an employee record
func (c *Client) CreateVote(val string) (int64, error) {
	ctx := context.Background()
	var err error

	if c.db == nil {
		err = errors.New("CreateEmployee: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = c.db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := "INSERT INTO VoteSchema.Votes (Vote) VALUES (@Vote); select convert(bigint, SCOPE_IDENTITY());"

	stmt, err := c.db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Vote", val),
	)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}

// ReadEmployees reads all employee records
func (c *Client) CountVote(vote string) (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := c.db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := "SELECT COUNT (Id) FROM VoteSchema.Votes WHERE Vote = @Name;"

	var count int

	stmt, err := c.db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", vote),
	)

	err = row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (c *Client) DeleteVotes() (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := c.db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("DELETE FROM VoteSchema.Votes;")

	// Execute non-query with named parameters
	result, err := c.db.ExecContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
