package main

import (
	"github.com/fseldow/AKS_daily_cleanup/cleanup/app"
)

func main() {
	framework := app.NewFramework()
	framework.CleanupNamespace()
}
