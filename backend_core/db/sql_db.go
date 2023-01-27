package db

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	scheduler "fermion/backend_core/controllers/scheduler/model"
	"fermion/backend_core/internal/model/accounting"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/offers"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/payments"
	"fermion/backend_core/internal/model/rating"
	"fermion/backend_core/internal/model/returns"
	"fermion/backend_core/internal/model/shipping"
	"fermion/backend_core/pkg/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
 Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Lesser General Public License v3.0 as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Lesser General Public License v3.0 for more details.
 You should have received a copy of the GNU Lesser General Public License v3.0
 along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/

var db *gorm.DB
var err error

func Init() {
	var (
		DB_HOST    = os.Getenv("DB_HOST")
		DB_USER    = os.Getenv("DB_USER")
		DB_PASS    = os.Getenv("DB_PASS")
		DB_NAME    = os.Getenv("DB_NAME")
		DB_PORT    = os.Getenv("DB_PORT")
		DB_SSLMODE = os.Getenv("DB_SSLMODE")
		DB_TZ      = os.Getenv("DB_TZ")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT, DB_SSLMODE, DB_TZ)

	gormCustomLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             1 * time.Second, // Slow SQL threshold
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			LogLevel:                  logger.Warn,     // Log level
			Colorful:                  true,            // Disable color
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormCustomLogger,
	})

	if err != nil {
		panic("failed to connect sql database")
	}
	clearDB := flag.Bool("clearDB", false, "a bool")
	seedMD := flag.Bool("seedMD", false, "a bool")
	seedTD := flag.Bool("seedTD", false, "a bool")
	migrateDB := flag.Bool("migrateDB", false, "a bool")
	seedMeta := flag.String("seedMeta", "all", "a string")
	// migrateDBVersion := flag.Bool("migrateDBVersion", false, "migrate version")
	flag.Parse()
	if *clearDB {
		fmt.Println("Deleting all data from database...")
		db.Exec("DROP SCHEMA public CASCADE")
		db.Exec("CREATE SCHEMA public")
		fmt.Println("Successfully deleted all data from database...")
	}
	if *migrateDB {

		// reset redis-db
		RedisClient := RedisManager()
		RedisClient.FlushAllAsync()
		fmt.Println("Successfully deleted all data from cache...")

		db.AutoMigrate(model_core.Tables...)
		db.AutoMigrate(model_core.CoreUsersTables...)
		db.AutoMigrate(mdm.MdmTables...)
		db.AutoMigrate(shared_pricing_and_location.SharedLocationandPricingTables...)
		db.AutoMigrate(inventory_orders.InventoryOrdersTables...)
		db.AutoMigrate(inventory_tasks.InventoryTasksTables...)
		db.AutoMigrate(orders.OrdersTables...)
		db.AutoMigrate(returns.ReturnsTables...)
		db.AutoMigrate(omnichannel.Omnichanneltables...)
		db.AutoMigrate(shipping.ShippingTables...)
		db.AutoMigrate(accounting.AccountingTables...)
		db.AutoMigrate(model_core.MetaTables...)
		db.AutoMigrate(scheduler.SchedulerTables...)
		db.AutoMigrate(model_core.AccessTables...)
		db.AutoMigrate(payments.PaymentsTables...)
		db.AutoMigrate(rating.RatingTables...)
		db.AutoMigrate(&util.MigrationVersionControl{})
		db.AutoMigrate(offers.OffersTables...)
		util.SeedMetaTable(db, seedMeta)
	}
	if *seedMD {
		util.SeedMasterData(db)
	}
	if *seedTD {
		util.SeedTestData(db)
	}
	// if *migrateDBVersion {
	util.Migrate(db)
	// }
	// if !*clearDB && (*seedMeta != "" && !*migrateDB) {
	// 	util.SeedMetaTable(db, seedMeta)
	// }
}

func DbManager() *gorm.DB {
	return db
}
