package outputs_kinesis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf/config"
	_ "github.com/influxdata/telegraf/migrations/outputs_kinesis" // register migration
	_ "github.com/influxdata/telegraf/plugins/outputs/kinesis"    // register plugin
	_ "github.com/influxdata/telegraf/plugins/serializers/influx" // register serializer
)

func TestPartitionKeyConflict(t *testing.T) {
	cfg := []byte(`
[[outputs.kinesis]]
  partitionkey = "test_key"
  [outputs.kinesis.partition]
    key = "test_key2"
    method = "static"
	`)
	// Migrate and check that nothing changed
	output, n, err := config.ApplyMigrations(cfg)
	require.ErrorContains(t, err, "contradicting setting for 'partitionkey' and 'partition.key'")
	require.Empty(t, output)
	require.Zero(t, n)
}

func TestPartitionMethodConflict(t *testing.T) {
	cfg := []byte(`
[[outputs.kinesis]]
  partitionkey = "test_key"
  [outputs.kinesis.partition]
    key = "test_key"
    method = "random"
	`)
	// Migrate and check that nothing changed
	output, n, err := config.ApplyMigrations(cfg)
	require.ErrorContains(t, err, "contradicting setting for 'use_random_partitionkey' and 'partition.method'")
	require.Empty(t, output)
	require.Zero(t, n)
}

func TestCases(t *testing.T) {
	// Get all directories in testdata
	folders, err := os.ReadDir("testcases")
	require.NoError(t, err)

	for _, f := range folders {
		// Only handle folders
		if !f.IsDir() {
			continue
		}

		t.Run(f.Name(), func(t *testing.T) {
			testcasePath := filepath.Join("testcases", f.Name())
			inputFile := filepath.Join(testcasePath, "telegraf.conf")
			expectedFile := filepath.Join(testcasePath, "expected.conf")

			// Read the expected output
			expected := config.NewConfig()
			require.NoError(t, expected.LoadConfig(expectedFile))
			require.NotEmpty(t, expected.Outputs)

			// Read the input data
			input, remote, err := config.LoadConfigFile(inputFile)
			require.NoError(t, err)
			require.False(t, remote)
			require.NotEmpty(t, input)

			// Migrate
			output, n, err := config.ApplyMigrations(input)
			require.NoError(t, err)
			require.NotEmpty(t, output)
			require.GreaterOrEqual(t, n, uint64(1))
			actual := config.NewConfig()
			require.NoError(t, actual.LoadConfigData(output, config.EmptySourcePath))

			// Test the output
			require.Len(t, actual.Outputs, len(expected.Outputs))
			actualIDs := make([]string, 0, len(expected.Outputs))
			expectedIDs := make([]string, 0, len(expected.Outputs))
			for i := range actual.Outputs {
				actualIDs = append(actualIDs, actual.Outputs[i].ID())
				expectedIDs = append(expectedIDs, expected.Outputs[i].ID())
			}
			require.ElementsMatch(t, expectedIDs, actualIDs, string(output))
		})
	}
}
