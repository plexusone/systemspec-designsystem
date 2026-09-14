// dss is the systemspec-designsystem CLI tool.
// It generates code artifacts and validates implementations against DSS specifications.
package main

import (
	"os"

	"github.com/plexusone/systemspec-designsystem/cmd/dss/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
