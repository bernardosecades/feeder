package main

import (
	"context"
	"github.com/bernardosecades/feeder/pkg/logger"
	"github.com/bernardosecades/feeder/pkg/repository"
	"github.com/bernardosecades/feeder/pkg/server"
	"github.com/bernardosecades/feeder/pkg/service"
	"github.com/bernardosecades/feeder/pkg/tools/env"
	"log"
	"time"
)

func main()  {
	cf := server.Config{
		Protocol:  "tcp",
		Host:      "",
		Port:      env.GetEnvOrFallback("SVC_PORT", "3333"),
		KeepAlive: time.Second * 60,
		MaxConn:   5,
	}

	l := logger.NewFileLogger("feeder_" + time.Now().Format(time.RFC3339Nano) + ".log")

	skuRepository := repository.NewSkuPostgreSQL(
		env.GetEnvOrFallback("DB_HOST", "localhost"),
		env.GetEnvOrFallback("DB_PORT", "5416"),
		env.GetEnvOrFallback("DB_USER", "feeder"),
		env.GetEnvOrFallback("DB_PASS", "feeder"),
		env.GetEnvOrFallback("DB_NAME", "feeder"),
	)

	sku := service.NewService(skuRepository, l)

	srv := server.NewServer(cf, sku)
	if err := srv.Start(context.Background()); err != nil {
		log.Println(err)
	}
}
