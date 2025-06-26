package main

import (
    "github.com/joho/godotenv"
    "go.uber.org/zap"
    "os"
)

func main_zap() { //remove zap from the name if want to run this program
    _ = godotenv.Load()

    var logger *zap.Logger
    var err error

    if os.Getenv("ENVIRONMENT") == "Development" {
        logger, err = zap.NewDevelopment()
    } else {
        logger, err = zap.NewProduction()
    }

    if err != nil {
        panic(err)
    }
    defer logger.Sync()

    logger.Info("Hello from Zap logger!")

    logger.Warn("User account is nearing the storage limit",
        zap.String("username", "john.doe"),
        zap.Float64("storageUsed", 4.5),
        zap.Float64("storageLimit", 5.0),
    )
}
