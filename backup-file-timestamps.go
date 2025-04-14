package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type FileAttrs struct {
	M int64 `json:"m"`
}

type FileAttrsMu struct {
	mu   sync.Mutex
	data map[string]FileAttrs
}

func collectFileAttrs(wg *sync.WaitGroup, path string, orgPathLen int, fileAttrs *FileAttrsMu, num_of_folders *atomic.Uint32, num_of_files *atomic.Uint32, ch_Err chan error) {
	defer wg.Done()
	entries, err := os.ReadDir(path)
	if err != nil {
		ch_Err <- err
		return
	}
	var dirs []string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			ch_Err <- err
			return
		}
		if entry.IsDir() {
			fullPath := filepath.Join(path, entry.Name())
			dirs = append(dirs, fullPath)
			num_of_folders.Add(1)
		} else {
			relPath := filepath.Join(path, entry.Name())[orgPathLen:]
			fileAttrs.mu.Lock()
			fileAttrs.data[relPath] = FileAttrs{
				M: info.ModTime().UnixNano(),
			}
			fileAttrs.mu.Unlock()
			num_of_files.Add(1)
		}
	}
	for _, dir := range dirs {
		wg.Add(1)
		go collectFileAttrs(wg, dir, orgPathLen, fileAttrs, num_of_folders, num_of_files, ch_Err)
	}
	ch_Err <- nil
}

func applyFileAttrs(restorePath string, attrs map[string]FileAttrs, num_of_missing *int, num_of_updated *int, num_of_skipped *int) error {
	for relPath, attr := range attrs {
		var path = filepath.Join(restorePath, relPath)
		var fileInfo, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Skipping missing file %s\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "Skipping file %s cuz of error: %v\n", path, err)
			}
			*num_of_missing += 1
			continue
		}
		currentMtime := fileInfo.ModTime().UnixNano()
		savedMtime := attr.M
		if currentMtime != savedMtime {
			fmt.Printf("Updating mtime for %s\n", path)
			mtime := time.Unix(0, savedMtime)
			err = os.Chtimes(path, time.Time{}, mtime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't change timestamps for file %s cuz of error: %v\n", path, err)
				*num_of_missing += 1
				continue
			}
			*num_of_updated += 1
		} else {
			*num_of_skipped += 1
		}
	}
	return nil
}

func dirPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("error: invalid path to folder provided")
	}
	if !info.IsDir() {
		return "", fmt.Errorf("error: path is not a directory")
	}
	return path, nil
}

func main() {
	const attrFileName = ".saved-file-timestamps"

	var savePath string
	var restorePath string
	var defaultPath string

	flag.StringVar(&savePath, "save", "", "Save the timestamps of files in the directory tree (used by default)")
	flag.StringVar(&restorePath, "restore", "", "Restore saved file timestamps in the directory tree")
	flag.Parse()

	if len(flag.Args()) > 0 {
		var err error
		defaultPath, err = dirPath(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			time.Sleep(7 * time.Second)
			os.Exit(1)
		}
	}

	if restorePath != "" {
		restorePath, err := dirPath(restorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			time.Sleep(7 * time.Second)
			os.Exit(1)
		}

		filePath := filepath.Join(restorePath, attrFileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Timestamps file '%s' not found\n", filePath)
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}

		fmt.Print("Are you sure you want to restore timestamps to previously saved ones? (Y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				time.Sleep(7 * time.Second)
				os.Exit(1)
			}

			var attrs map[string]FileAttrs
			if err := json.Unmarshal(data, &attrs); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing timestamps file JSON: %v\n Wrong format?\n", err)
				// using old format?
				time.Sleep(7 * time.Second)
				os.Exit(1)
			}
			var num_of_missing, num_of_updated, num_of_skipped int
			var start = time.Now()
			fmt.Printf("Restoring timestamps from: '%s'\n...\n", filePath)
			applyFileAttrs(restorePath, attrs, &num_of_missing, &num_of_updated, &num_of_skipped)

			fmt.Printf("...\nResults:\n\n")

			fmt.Printf("- files skipped (missing from HDD): %d\n", num_of_updated)
			fmt.Printf("- files skipped (same timestamps) : %d\n", num_of_skipped)
			fmt.Printf("- files we updated timestamps for : %d\n", num_of_updated)

			fmt.Printf("= in %.2f sec\n\nRestore complete!\n", time.Since(start).Seconds())
		} else {
			fmt.Println("Aborted")
			os.Exit(1)
		}
	} else {
		var path string
		if savePath != "" {
			path = savePath
		} else if defaultPath != "" {
			path = defaultPath
		} else {
			path, _ = os.Getwd()
		}

		filePath := filepath.Join(path, attrFileName)
		fmt.Printf("Saving timestamps to: '%s'\nwith:\n", filePath)

		var wg sync.WaitGroup
		var fileAttrsMu = FileAttrsMu{data: make(map[string]FileAttrs)}
		var start = time.Now()
		var num_of_folders, num_of_files atomic.Uint32
		var ch_Err = make(chan error)

		wg.Add(1)
		go collectFileAttrs(&wg, path, len(path)+1, &fileAttrsMu, &num_of_folders, &num_of_files, ch_Err)
		go func() {
			wg.Wait()
			close(ch_Err)
		}()
		for err := range ch_Err {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error collecting file attributes: %v\n", err)
			}
		}
		fmt.Printf("- folders: %d\n- files: %d\n= in %.2f sec\n", num_of_folders.Load(), num_of_files.Load(), time.Since(start).Seconds())

		data, err := json.Marshal(fileAttrsMu.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			time.Sleep(7 * time.Second)
			os.Exit(1)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			time.Sleep(7 * time.Second)
			os.Exit(1)
		}

		fmt.Println("\nSave complete!")
	}

	time.Sleep(3 * time.Second)
}
