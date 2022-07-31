package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Jibaru/go-db/pkg/invoice"
	"github.com/Jibaru/go-db/pkg/invoiceheader"
	"github.com/Jibaru/go-db/pkg/invoiceitem"
	"github.com/Jibaru/go-db/pkg/product"
	"github.com/Jibaru/go-db/storage"
)

type Action string

const (
	create      Action = "CREATE"
	update      Action = "UPDATE"
	delete      Action = "DELETE"
	getAll      Action = "GETALL"
	getOne      Action = "GETONE"
	transaction Action = "TRANSACTION"
	migrate     Action = "MIGRATE"
)

func main() {
	const (
		driver        = storage.MySQL
		action Action = create
	)
	storage.New(driver)

	storageProduct, err := storage.DAOProduct(driver)
	if err != nil {
		log.Fatalf("DAOProduct: %v", err)
	}
	serviceProduct := product.NewService(storageProduct)

	storageInvoiceHeader, err := storage.DAOInvoiceHeader(driver)
	if err != nil {
		log.Fatalf("DAOInvoiceHeader %v", err)
	}
	serviceInvoiceHeader := invoiceheader.NewService(storageInvoiceHeader)

	storageInvoiceItem, err := storage.DAOInvoiceItem(driver)
	if err != nil {
		log.Fatalf("DAOInvoiceItem %v", err)
	}
	serviceInvoiceItem := invoiceitem.NewService(storageInvoiceItem)

	storageInvoice, err := storage.DAOInvoice(
		driver,
		storageInvoiceHeader,
		storageInvoiceItem,
	)
	if err != nil {
		log.Fatalf("DAOInvoice %v", err)
	}
	serviceInvoice := invoice.NewService(storageInvoice)

	switch action {
	case migrate:
		runMigrations(
			serviceProduct,
			serviceInvoiceHeader,
			serviceInvoiceItem,
		)
	case create:
		createProduct(serviceProduct)
	case getAll:
		getAllProducts(serviceProduct)
	case getOne:
		getProductByID(serviceProduct)
	case update:
		updateProduct(serviceProduct)
	case delete:
		deleteProduct(serviceProduct)
	case transaction:
		runTransaction(serviceInvoice)
	}
}

func runMigrations(
	serviceProduct *product.Service,
	serviceInvoiceHeader *invoiceheader.Service,
	serviceInvoiceItem *invoiceitem.Service,
) {
	if err := serviceProduct.Migrate(); err != nil {
		log.Fatalf("product.Migrate: %v", err)
	}

	if err := serviceInvoiceHeader.Migrate(); err != nil {
		log.Fatalf("invoiceheader.Migrate: %v", err)
	}

	if err := serviceInvoiceItem.Migrate(); err != nil {
		log.Fatalf("invoiceitem.Migrate: %v", err)
	}
}

func createProduct(serviceProduct *product.Service) {
	m := &product.Model{
		Name:         "Java Course",
		Price:        56,
		Observations: "On fire",
	}

	if err := serviceProduct.Create(m); err != nil {
		log.Fatalf("product.Create: %v", err)
	}

	fmt.Printf("%+v\n", m)
}

func getAllProducts(serviceProduct *product.Service) {
	ms, err := serviceProduct.GetAll()
	if err != nil {
		log.Fatalf("product.GetAll: %v", err)
	}
	fmt.Printf("%+v\n", ms)
}

func getProductByID(serviceProduct *product.Service) {
	const id uint = 1
	m, err := serviceProduct.GetByID(id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		fmt.Printf("there is no product with id: %v\n", id)
	case err != nil:
		log.Fatalf("product.GetByID: %v", err)
	default:
		fmt.Printf("%+v\n", m)
	}
}

func updateProduct(serviceProduct *product.Service) {
	m := &product.Model{
		ID:           1,
		Name:         "Python Course",
		Price:        56,
		Observations: "This is the python course",
	}

	if err := serviceProduct.Update(m); err != nil {
		log.Fatalf("product.Update: %v", err)
	}

	fmt.Printf("%+v\n", m)
}

func deleteProduct(serviceProduct *product.Service) {
	const id uint = 1
	if err := serviceProduct.Delete(id); err != nil {
		log.Fatalf("product.Delete: %v", err)
	}
}

func runTransaction(serviceInvoice *invoice.Service) {
	m := &invoice.Model{
		Header: &invoiceheader.Model{
			Client: "Ignacio",
		},
		Items: invoiceitem.Models{
			&invoiceitem.Model{ProductID: 2},
			&invoiceitem.Model{ProductID: 3},
		},
	}

	if err := serviceInvoice.Create(m); err != nil {
		log.Fatalf("invoice.Create: %v", err)
	}
}
