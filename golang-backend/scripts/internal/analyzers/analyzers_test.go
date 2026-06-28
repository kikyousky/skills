package analyzers_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/boundaries"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/complexity"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/getenv"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/jsondecode"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/requests"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/serviceparams"
	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/transportflow"
)

func TestBoundaries(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), boundaries.Analyzer, "example.com/boundaries/internal/...")
}

func TestRequests(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), requests.Analyzer, "example.com/requests/internal/...")
}

func TestServiceParams(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), serviceparams.Analyzer, "example.com/serviceparams/internal/...")
}

func TestTransportFlow(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), transportflow.Analyzer, "example.com/transportflow/internal/...")
}

func TestJSONDecode(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), jsondecode.Analyzer, "example.com/jsondecode/internal/...")
}

func TestGetenv(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), getenv.Analyzer, "example.com/getenv/internal/...")
}

func TestComplexity(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), complexity.Analyzer, "example.com/complexity/internal/...")
}
