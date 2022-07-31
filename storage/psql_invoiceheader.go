package storage

import (
	"database/sql"
	"fmt"

	"github.com/Jibaru/go-db/pkg/invoiceheader"
)

const (
	psqlMigrateInvoiceHeader = `
	CREATE TABLE IF NOT EXISTS invoice_headers(
		id SERIAL NOT NULL,
		client VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT invoice_headers_id_pk PRIMARY KEY (id)
	)`
	psqlCreateInvoiceHeader = `
	INSERT INTO invoice_headers(client) 
	VALUES($1) 
	RETURNING id, created_at`
)

// psqlInvoiceHeader used for work with postgres - invoiceHeader
type psqlInvoiceHeader struct {
	db *sql.DB
}

// newPsqlInvoiceHeader return a new pointer of PsqlInvoiceHeader
func newPsqlInvoiceHeader(db *sql.DB) *psqlInvoiceHeader {
	return &psqlInvoiceHeader{db}
}

// Migrate implement the interface invoiceheader.Storage
func (p *psqlInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	fmt.Println("invoice_headers table migrated successfully")
	return nil
}

// CreateTx implement the interface invoiceHeader.Storage
func (p *psqlInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(psqlCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(m.Client).Scan(&m.ID, &m.CreatedAt)
}
