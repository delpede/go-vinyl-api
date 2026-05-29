package main

import (
	"context"

	"github.com/delpede/go-vinyl-api/internal/api/generated"
	"github.com/delpede/go-vinyl-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func (s *Server) GetRecords(c *gin.Context, params generated.GetRecordsParams) {
	query := `
		SELECT id, title, artist, created_at, updated_at
		FROM records
	`
	args := []any{}

	// Optional full-text-ish search via the `q` query parameter.
	// TODO: extend the match to `label` and `notes` once those columns exist.
	if params.Q != nil && *params.Q != "" {
		query += `
		WHERE title ILIKE $1 OR artist ILIKE $1
		`
		args = append(args, "%"+*params.Q+"%")
	}

	query += `
		ORDER BY created_at DESC
	`

	rows, err := s.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	records := []generated.Record{}
	for rows.Next() {
		var record generated.Record

		err := rows.Scan(
			&record.Id,
			&record.Title,
			&record.Artist,
			&record.CreatedAt,
			&record.UpdatedAt,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records = append(records, record)
	}

	c.JSON(200, records)
}

func (s *Server) PostRecords(c *gin.Context) {
	var req generated.CreateRecordRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var record generated.Record
	err := s.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO records (title, artist)
		VALUES ($1, $2)
		RETURNING id, title, artist, created_at, updated_at
		`,
		req.Title,
		req.Artist,
	).Scan(
		&record.Id,
		&record.Title,
		&record.Artist,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, record)
}

func (s *Server) GetRecordsId(c *gin.Context, id generated.RecordId) {
	var record generated.Record

	err := s.DB.QueryRow(
		context.Background(),
		`
		SELECT id, title, artist, created_at, updated_at
		FROM records
		WHERE id = $1
		`,
		id,
	).Scan(
		&record.Id,
		&record.Title,
		&record.Artist,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(404, gin.H{"error": "record not found"})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, record)
}

func (s *Server) PutRecordsId(c *gin.Context, id generated.RecordId) {
	var req generated.UpdateRecordRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var record generated.Record
	err := s.DB.QueryRow(
		context.Background(),
		`
		UPDATE records
		SET title = $1, artist = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, artist, created_at, updated_at
		`,
		req.Title,
		req.Artist,
		id,
	).Scan(
		&record.Id,
		&record.Title,
		&record.Artist,
		&record.CreatedAt,
		&record.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(404, gin.H{"error": "record not found"})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, record)
}

func (s *Server) DeleteRecordsId(c *gin.Context, id generated.RecordId) {
	commandTag, err := s.DB.Exec(
		context.Background(),
		`
		DELETE FROM records
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if commandTag.RowsAffected() == 0 {
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}

	c.Status(204)
}

func main() {
	router := gin.Default()

	dbpool, err := db.New()
	if err != nil {
		panic(err)
	}

	server := &Server{
		DB: dbpool,
	}
	generated.RegisterHandlers(router, server)

	router.Run(":8080")
}
