// internal/store/document.go
package store

import (
	"context"
	"database/sql"
	"time"
)

// Document represents a shared document.
type Document struct {
	ID        int64     `json:"id"`
	DocID     string    `json:"doc_id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DocumentStore defines methods for document operations.
type DocumentStore struct {
	db *sql.DB
}

// NewDocumentStore creates a new DocumentStore.
func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

// CreateDocument inserts a new document.
func (ds *DocumentStore) CreateDocument(ctx context.Context, docID, content string) error {
	query := `
		INSERT INTO documents (doc_id, content)
		VALUES ($1, $2)
	`
	_, err := ds.db.ExecContext(ctx, query, docID, content)
	return err
}

// GetDocumentByDocID retrieves a document by its docID.
func (ds *DocumentStore) GetDocumentByDocID(ctx context.Context, docID string) (*Document, error) {
	query := `
		SELECT id, doc_id, content, updated_at
		FROM documents
		WHERE doc_id = $1
	`
	doc := &Document{}
	err := ds.db.QueryRowContext(ctx, query, docID).Scan(&doc.ID, &doc.DocID, &doc.Content, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// UpdateDocument updates a document's content.
func (ds *DocumentStore) UpdateDocument(ctx context.Context, docID, content string) error {
	query := `
		UPDATE documents
		SET content = $1, updated_at = NOW()
		WHERE doc_id = $2
	`
	res, err := ds.db.ExecContext(ctx, query, content, docID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
