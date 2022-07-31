package storage

import (
	"database/sql"

	"github.com/Jibaru/go-db/pkg/invoice"
	"github.com/Jibaru/go-db/pkg/invoiceheader"
	"github.com/Jibaru/go-db/pkg/invoiceitem"
)

// mySQLInvoice used for work with MySQL - invoice
type mySQLInvoice struct {
	db            *sql.DB
	storageHeader invoiceheader.Storage
	storageItems  invoiceitem.Storage
}

// newMySQLInvoice return a new pointer of MySQLInvoice
func newMySQLInvoice(
	db *sql.DB,
	h invoiceheader.Storage,
	i invoiceitem.Storage,
) *mySQLInvoice {
	return &mySQLInvoice{
		db:            db,
		storageHeader: h,
		storageItems:  i,
	}
}

// Create implement the interface invoice.Storage
func (p *mySQLInvoice) Create(m *invoice.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err := p.storageHeader.CreateTx(tx, m.Header); err != nil {
		tx.Rollback()
		return err
	}

	if err := p.storageItems.CreateTx(tx, m.Header.ID, m.Items); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
