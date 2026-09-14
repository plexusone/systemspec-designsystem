package dss

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadDesignSystem loads a design system from a directory or single file.
func LoadDesignSystem(path string) (*DesignSystem, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat path: %w", err)
	}

	if info.IsDir() {
		return loadFromDirectory(path)
	}
	return loadFromFile(path)
}

// LoadDesignSystemFromFS loads a design system from an fs.FS interface.
// This is useful for loading from embed.FS or other virtual filesystems.
//
// Example with embed.FS:
//
//	//go:embed spec/*
//	var specFS embed.FS
//
//	// Load from subdirectory "spec" within the embedded FS
//	sub, _ := fs.Sub(specFS, "spec")
//	ds, err := dss.LoadDesignSystemFromFS(sub)
func LoadDesignSystemFromFS(fsys fs.FS) (*DesignSystem, error) {
	// Check if there's a single design-system file at the root
	for _, ext := range []string{".json", ".yaml", ".yml"} {
		path := "design-system" + ext
		data, err := fs.ReadFile(fsys, path)
		if err == nil {
			var ds DesignSystem
			if err := unmarshal(path, data, &ds); err != nil {
				return nil, fmt.Errorf("unmarshal %s: %w", path, err)
			}
			return &ds, nil
		}
	}

	// Load from directory structure
	return loadFromDirectoryFS(fsys)
}

// loadFromDirectoryFS loads a design system from an fs.FS directory structure.
func loadFromDirectoryFS(fsys fs.FS) (*DesignSystem, error) {
	ds := &DesignSystem{}

	// Load meta
	if err := loadLayerFS(fsys, "meta", &ds.Meta); err != nil {
		return nil, fmt.Errorf("load meta: %w", err)
	}

	// Load principles
	if err := loadLayerFS(fsys, "modes", &ds.Modes); err != nil && !isNotExist(err) {
		return nil, err
	}
	if err := loadLayerFS(fsys, "principles", &ds.Principles); err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("load principles: %w", err)
	}

	// Load foundations
	if err := loadFoundationsFS(fsys, &ds.Foundations); err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("load foundations: %w", err)
	}

	// Load components
	if err := loadCollectionFS(fsys, "components", &ds.Components); err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("load components: %w", err)
	}

	// Load patterns
	if err := loadCollectionFS(fsys, "patterns", &ds.Patterns); err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("load patterns: %w", err)
	}

	// Load templates
	if err := loadCollectionFS(fsys, "templates", &ds.Templates); err != nil && !isNotExist(err) {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	// Load content
	var content Content
	if err := loadLayerFS(fsys, "content", &content); err == nil {
		ds.Content = &content
	}

	// Load accessibility
	var a11y Accessibility
	if err := loadLayerFS(fsys, "accessibility", &a11y); err == nil {
		ds.Accessibility = &a11y
	}

	// Load governance
	var gov Governance
	if err := loadLayerFS(fsys, "governance", &gov); err == nil {
		ds.Governance = &gov
	}

	return ds, nil
}

// loadLayerFS loads a single layer file from an fs.FS.
func loadLayerFS(fsys fs.FS, name string, v interface{}) error {
	for _, ext := range []string{".json", ".yaml", ".yml"} {
		path := name + ext
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			if isNotExist(err) {
				continue
			}
			return err
		}
		return unmarshal(path, data, v)
	}
	return fs.ErrNotExist
}

// loadFoundationsFS loads foundations from an fs.FS.
func loadFoundationsFS(fsys fs.FS, f *Foundations) error {
	// Try loading from foundations subdirectory
	subFS, err := fs.Sub(fsys, "foundations")
	if err != nil {
		// Try loading from single file
		return loadLayerFS(fsys, "foundations", f)
	}

	// Check if it's a valid directory by trying to read it
	if _, err := fs.ReadDir(subFS, "."); err != nil {
		return loadLayerFS(fsys, "foundations", f)
	}

	// Load individual token files
	if err := loadLayerFS(subFS, "colors", &f.Colors); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "typography", &f.Typography); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "spacing", &f.Spacing); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "elevation", &f.Elevation); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "motion", &f.Motion); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "grid", &f.Grid); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "breakpoints", &f.Breakpoints); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "border-radius", &f.BorderRadius); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "border-width", &f.BorderWidth); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "opacity", &f.Opacity); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "densities", &f.Densities); err != nil && !isNotExist(err) {
		return err
	}
	if err := loadLayerFS(subFS, "z-index", &f.ZIndex); err != nil && !isNotExist(err) {
		return err
	}

	return nil
}

// loadCollectionFS loads a collection of items from an fs.FS directory.
func loadCollectionFS[T any](fsys fs.FS, name string, items *[]T) error {
	// Try loading from subdirectory
	subFS, err := fs.Sub(fsys, name)
	if err != nil {
		// Try loading from single file
		return loadLayerFS(fsys, name, items)
	}

	// Check if it's a valid directory
	entries, err := fs.ReadDir(subFS, ".")
	if err != nil {
		return loadLayerFS(fsys, name, items)
	}

	// Load each file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		data, err := fs.ReadFile(subFS, entry.Name())
		if err != nil {
			return err
		}

		var item T
		if err := unmarshal(entry.Name(), data, &item); err != nil {
			return fmt.Errorf("unmarshal %s: %w", entry.Name(), err)
		}
		*items = append(*items, item)
	}

	return nil
}

// isNotExist checks if an error indicates a file or directory does not exist.
func isNotExist(err error) bool {
	return err == fs.ErrNotExist || os.IsNotExist(err)
}

// loadFromFile loads a design system from a single JSON or YAML file.
func loadFromFile(path string) (*DesignSystem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var ds DesignSystem
	if err := unmarshal(path, data, &ds); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &ds, nil
}

// loadFromDirectory loads a design system from a directory structure.
// Expected structure:
//
//	dir/
//	  meta.json (or meta.yaml)
//	  foundations/
//	    colors.json
//	    typography.json
//	    ...
//	  components/
//	    button.json
//	    ...
//
// Alternatively, loads from a single design-system.json or design-system.yaml file.
func loadFromDirectory(dir string) (*DesignSystem, error) {
	// First, check for a single design-system.json or design-system.yaml file
	for _, ext := range []string{".json", ".yaml", ".yml"} {
		path := filepath.Join(dir, "design-system"+ext)
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path)
		}
	}

	ds := &DesignSystem{}

	// Load meta
	if err := loadLayer(dir, "meta", &ds.Meta); err != nil {
		return nil, fmt.Errorf("load meta: %w", err)
	}

	// Load principles
	if err := loadLayer(dir, "modes", &ds.Modes); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err := loadLayer(dir, "principles", &ds.Principles); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load principles: %w", err)
	}

	// Load foundations
	if err := loadFoundations(dir, &ds.Foundations); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load foundations: %w", err)
	}

	// Load components
	if err := loadCollection(dir, "components", &ds.Components); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load components: %w", err)
	}

	// Load patterns
	if err := loadCollection(dir, "patterns", &ds.Patterns); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load patterns: %w", err)
	}

	// Load templates
	if err := loadCollection(dir, "templates", &ds.Templates); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	// Load content
	var content Content
	if err := loadLayer(dir, "content", &content); err == nil {
		ds.Content = &content
	}

	// Load accessibility
	var a11y Accessibility
	if err := loadLayer(dir, "accessibility", &a11y); err == nil {
		ds.Accessibility = &a11y
	}

	// Load governance
	var gov Governance
	if err := loadLayer(dir, "governance", &gov); err == nil {
		ds.Governance = &gov
	}

	return ds, nil
}

// loadLayer loads a single layer file (meta.json, principles.json, etc.).
func loadLayer(dir, name string, v interface{}) error {
	for _, ext := range []string{".json", ".yaml", ".yml"} {
		path := filepath.Join(dir, name+ext)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		return unmarshal(path, data, v)
	}
	return os.ErrNotExist
}

// loadFoundations loads foundations from a directory of token files.
func loadFoundations(dir string, f *Foundations) error {
	foundationsDir := filepath.Join(dir, "foundations")
	info, err := os.Stat(foundationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Try loading from single file
			return loadLayer(dir, "foundations", f)
		}
		return err
	}
	if !info.IsDir() {
		return loadLayer(dir, "foundations", f)
	}

	// Load individual token files
	if err := loadLayer(foundationsDir, "colors", &f.Colors); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "typography", &f.Typography); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "spacing", &f.Spacing); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "elevation", &f.Elevation); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "motion", &f.Motion); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "grid", &f.Grid); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "breakpoints", &f.Breakpoints); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "border-radius", &f.BorderRadius); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "border-width", &f.BorderWidth); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "opacity", &f.Opacity); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "densities", &f.Densities); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := loadLayer(foundationsDir, "z-index", &f.ZIndex); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// loadCollection loads a collection of items from a directory.
func loadCollection[T any](dir, name string, items *[]T) error {
	collectionDir := filepath.Join(dir, name)
	info, err := os.Stat(collectionDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Try loading from single file
			return loadLayer(dir, name, items)
		}
		return err
	}
	if !info.IsDir() {
		return loadLayer(dir, name, items)
	}

	// Load each file in the directory
	entries, err := os.ReadDir(collectionDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(collectionDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var item T
		if err := unmarshal(path, data, &item); err != nil {
			return fmt.Errorf("unmarshal %s: %w", entry.Name(), err)
		}
		*items = append(*items, item)
	}

	return nil
}

// unmarshal parses JSON or YAML data based on file extension.
func unmarshal(path string, data []byte, v interface{}) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, v)
	default:
		return json.Unmarshal(data, v)
	}
}

// SaveDesignSystem saves a design system to a JSON file.
func SaveDesignSystem(ds *DesignSystem, path string) error {
	data, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
