package main

import (
	"flag"
	"fmt"
	"log"

	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
)

func main() {
	apply := flag.Bool("apply", false, "apply destructive hard-delete cleanup (default is dry-run)")
	flag.Parse()

	config.LoadConfig()
	cfg := config.Get()
	if cfg == nil {
		log.Fatal("failed to load config")
	}

	if err := models.InitDB(cfg); err != nil {
		log.Fatal(err)
	}

	summary, errCode := services.RunHardDeleteCleanupWithOptions(!*apply)
	if errCode != protocol.Success {
		log.Fatalf("hard delete cleanup failed: %s", errCode)
	}

	fmt.Println("Hard delete cleanup completed.")
	if summary.DryRun {
		fmt.Println("Mode: DRY-RUN (no rows were deleted). Use --apply to execute.")
	}
	fmt.Println(summary.String())
}
