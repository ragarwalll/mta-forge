package forger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ragarwalll/mta-forge.git/pkg/constants"
)

// Generate handles the generation of either a descriptor or deployment based on user selection
func (f *Forger) Generate(generationType string) error {
	switch generationType {
	case "deployment":
		return f.generateDeployment()
	case "extension":
		return f.generateExtension()
	default:
		return fmt.Errorf("invalid generation type: %s", generationType)
	}
}

func (f *Forger) generateDeployment() error {
	slog.Info("Generating mta.yaml")

	mta, err := f.CreateDescriptor("")
	if err != nil {
		return fmt.Errorf("error creating descriptor: %w", err)
	}

	// Write the mta.yaml file
	slog.Info("Writing to", "path", filepath.Join(f.OutputDir, "mta.yaml"))
	err = os.WriteFile(filepath.Join(f.OutputDir, "mta.yaml"), []byte(mta), constants.DefaultFilePermissions)

	if err != nil {
		return fmt.Errorf("error writing mta.yaml: %w", err)
	}

	return nil
}

func (f *Forger) generateExtension() error {
	// Create 'descriptors' folder if it doesn't exist
	descriptorsDir := filepath.Join(f.OutputDir, "descriptors")

	slog.Info("Creating descriptors directory", "path", descriptorsDir)

	if err := os.MkdirAll(descriptorsDir, constants.DefaultDirPermissions); err != nil {
		slog.Error("Failed to create descriptors directory", "error", err)
		return fmt.Errorf("error creating descriptors directory: %w", err)
	}

	// List all folders in 'extensions' folder
	extensionsDir := filepath.Join(f.BaseDir, "extensions")

	slog.Info("Reading extensions directory", "path", extensionsDir)
	entries, err := os.ReadDir(extensionsDir)

	if err != nil {
		slog.Error("Failed to read extensions directory", "error", err)
		return fmt.Errorf("error reading extensions directory: %w", err)
	}

	slog.Debug("Found entries in extensions directory", "count", len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			slog.Debug("Skipping non-directory entry", "name", entry.Name())
			continue
		}

		slog.Info("Creating descriptor for extension", "name", entry.Name())
		descriptor, err := f.CreateDescriptor(filepath.Join("extensions", entry.Name()))

		if err != nil {
			slog.Error("Failed to create descriptor", "extension", entry.Name(), "error", err)
			return fmt.Errorf("error creating descriptor for %s: %w", entry.Name(), err)
		}

		// Write the descriptor to the file
		path := filepath.Join(descriptorsDir, fmt.Sprintf("%s.mtaext", entry.Name()))
		slog.Info("Writing descriptor to file", "path", path)
		err = os.WriteFile(path, []byte(descriptor), constants.DefaultFilePermissions)

		if err != nil {
			slog.Error("Failed to write descriptor file", "path", path, "error", err)
			return fmt.Errorf("error writing descriptor for %s: %w", entry.Name(), err)
		}

		slog.Debug("Successfully wrote descriptor file", "path", path)
	}

	slog.Info("Extension generation completed successfully")

	return nil
}
