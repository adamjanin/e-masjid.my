package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-saas/saas"

	"github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	gorm2 "gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Dev4w4n/e-masjid.my/saas/model"
)

type Seed struct {
	dbProvider gorm.DbProvider
	connStrGen saas.ConnStrGenerator
}

func NewSeed(dbProvider gorm.DbProvider, connStrGen saas.ConnStrGenerator) *Seed {
	return &Seed{dbProvider: dbProvider, connStrGen: connStrGen}
}

func (s *Seed) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := s.dbProvider.Get(ctx, "")

	if sCtx.TenantId == "" {
		//init host
		log.Println("Initializing host.")
		err := db.Model(&model.Tenant{}).Session(&gorm2.Session{FullSaveAssociations: true}).Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches([]model.Tenant{
			{ID: "1", Name: "Host"}}, 10).Error
		if err != nil {
			return err
		}
		entities := []model.Post{
			{
				Model:       gorm2.Model{ID: 1},
				Title:       fmt.Sprintf("Host Side"),
				Description: fmt.Sprintf("Init Host"),
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}

	if sCtx.TenantId == "1" {
		entities := []model.Post{
			{
				Model:       gorm2.Model{ID: 2},
				Title:       fmt.Sprintf("Tenant %s Post 1", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are one post in this tenant. This is post 1", sCtx.TenantId),
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}

	// if its a tenant then init the database
	if sCtx.TenantId != "" {
		log.Println("Executing Sql Files.")
		err := executeSqlFiles(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPosts(db *gorm2.DB, entities []model.Post) error {
	for _, entity := range entities {
		err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Model(&model.Post{}).Create(&entity).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func executeSqlFiles(db *gorm2.DB) error {
	// List .sql files in a directory
	// Update the pattern to match your file naming convention
	sqlFiles, err := filepath.Glob("/app/*.sql")
	if err != nil {
		panic(err)
	}

	log.Println("Sorting files.")
	// Sort files by name
	sort.Strings(sqlFiles)

	log.Println("Files count: ", sqlFiles)
	// Loop through each .sql file
	for _, sqlFile := range sqlFiles {
		log.Println("Processing file: ", sqlFile)
		// Read SQL file
		sqlBytes, err := os.ReadFile(sqlFile)
		if err != nil {
			panic(err)
		}

		// Split SQL statements
		sqlStatements := strings.Split(string(sqlBytes), ";")

		// Execute SQL statements
		for _, sqlStatement := range sqlStatements {
			sqlStatement = strings.TrimSpace(sqlStatement)
			if sqlStatement == "" {
				continue
			}
			err := db.Exec(sqlStatement).Error
			if err != nil {
				panic(err)
			}
		}
	}

	log.Println("Initialization complete.")

	return nil
}
