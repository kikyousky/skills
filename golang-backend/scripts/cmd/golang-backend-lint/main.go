package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/boundaries"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/complexity"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/getenv"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/jsondecode"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/requests"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/serviceparams"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/transportflow"
)

func main() {
	multichecker.Main(
		boundaries.Analyzer,
		complexity.Analyzer,
		getenv.Analyzer,
		jsondecode.Analyzer,
		requests.Analyzer,
		serviceparams.Analyzer,
		transportflow.Analyzer,
	)
}
