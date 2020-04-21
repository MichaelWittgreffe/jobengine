package main

import (
	"fmt"

	"github.com/MichaelWittgreffe/jobengine/api"
	"github.com/MichaelWittgreffe/jobengine/crypto"
	"github.com/MichaelWittgreffe/jobengine/database"
	"github.com/MichaelWittgreffe/jobengine/filesystem"
	"github.com/MichaelWittgreffe/jobengine/logger"
)

func main() {
	logger := logger.NewLogger("std")
	fileHandler := filesystem.NewFileSystem("os")
	dbFile := database.NewDBFile()

	dbPath := fileHandler.GetEnv("DB_PATH")
	if len(dbPath) <= 0 {
		logger.Info("DB_PATH Not Defined, Using Default")
		dbPath = "/jobengine/database.queuedb"
	}

	apiPort := fileHandler.GetEnv("API_PORT")
	if len(dbPath) <= 0 {
		logger.Info("aPI_PORT Not Defined, Using Default")
		apiPort = "80"
	}

	secretKey := fileHandler.GetEnv("SECRET")
	if len(secretKey) <= 0 {
		logger.Fatal("SECRET Not Defined")
	}

	dbFileHandler := database.NewDBFileHandler(
		"fs",
		crypto.NewEncryptionHandler(secretKey, "AES", crypto.NewHashHandler("md5")),
		database.NewDBDataHandler("json"),
		fileHandler,
	)

	if dbFileHandler == nil {
		logger.Fatal("Unable To Create File Handler")
	}

	if exists, err := fileHandler.FileExists(dbPath); err == nil {
		if exists {
			if err = dbFileHandler.LoadFromFile(dbFile, dbPath); err == nil {
				logger.Info("Database Loaded")
			}
		} else {
			if err = dbFileHandler.SaveToFile(dbFile, dbPath); err == nil {
				logger.Info("Database Created")
			}
		}
	} else {
		logger.Fatal(fmt.Sprintf("Error Locating DB File: %s", err.Error()))
	}

	dbFileMonitor := database.NewDBFileMonitor(dbFile, dbPath, dbFileHandler, logger)
	if dbFileMonitor == nil {
		logger.Fatal("Failed Creating Monitor")
	}
	go dbFileMonitor.Start()

	httpAPI := api.NewHTTPAPI(logger, dbFileMonitor, database.NewQueryController(dbFile, crypto.NewHashHandler("sha512")))
	httpAPI.ListenAndServe(apiPort)
}
