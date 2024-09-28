package transformer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/diegosz/go-archetype/log"
	"github.com/diegosz/go-archetype/operations"
	"github.com/diegosz/go-archetype/reader"
	"github.com/diegosz/go-archetype/types"
	"github.com/diegosz/go-archetype/writer"
)

type Transformer interface {
	GetName() string
	GetFilePatterns() []types.FilePattern
	Template(vars map[string]string) error
	Transform(types.File) types.File
}

// Transform performs the actual transformation of files from source to
// destination. If the destination is not empty, it will abort.
func Transform(source, destination string, transformations Transformations, logger log.Logger) error {
	return transform(source, destination, transformations, logger, false)
}

// OverlayTransform performs the actual transformation of files from source to
// destination. If the destination is not empty, it will overlay the files.
func OverlayTransform(source, destination string, transformations Transformations, logger log.Logger) error {
	return transform(source, destination, transformations, logger, true)
}

func transform(source, destination string, transformations Transformations, logger log.Logger, overlay bool) error {
	empty, err := isDirEmptyOrDoesntExist(destination)
	if err != nil {
		return err
	}
	if !overlay && !empty {
		logger.Errorf("Destination %s is not empty, aborting", destination)
		return errors.New("destination is not empty")
	}

	// Before actions
	err = before(transformations)
	if err != nil {
		return err
	}

	// All transformations
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking to file: %w", err)
		}
		sourceFile := path
		isDir, ignored, file, err := reader.ReadFile(sourceFile, info, source, transformations.IsGloballyIgnored)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		if isDir {
			return nil
		}

		if ignored {
			logger.Debugf("Ignoring file %s", path)
		} else {
			file, err = transformations.Transform(file)
			if writeErr := writer.WriteFile(destination, file, info.Mode(), logger); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			return fmt.Errorf("transforming: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// After actions
	return after(transformations)
}

func before(ts Transformations) error {
	return executeOperators(ts.before)
}

func after(ts Transformations) error {
	return executeOperators(ts.after)
}

func executeOperators(ops []operations.Operator) error {
	for _, op := range ops {
		if err := op.Operate(); err != nil {
			return err
		}
	}
	return nil
}

func isDirEmptyOrDoesntExist(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// OK, does not exist
		return true, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		// Empty dir
		return true, nil
	}

	return false, err // Either not empty or error, suits both cases
}
