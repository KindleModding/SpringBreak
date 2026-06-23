package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	drivedetector "github.com/penguins184/drivedetector/src"
	"github.com/schollz/progressbar/v3"
)

func cleanup(drive string) error {
	base := filepath.Join(drive, ".active_content_sandbox", "store", "resource", "cachedResources")
	if _, err := os.Stat(base); err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	if err := os.Chdir(base); err != nil {
		return err
	}

	bar := progressbar.Default(5000, "Cleaning")

	for i := 1; i <= 5000; i++ {
		if err := os.Chdir(strconv.Itoa(i)); err != nil {
			break
		}
	}

	for {
		current, _ := os.Getwd()
		if current == base || current == drive || current == drive+"\\" {
			break
		}

		os.Chdir("..")
		os.RemoveAll(current)
		bar.Add(1)
	}

	os.Chdir(drive)

	return os.RemoveAll(filepath.Join(drive, ".active_content_sandbox"))
}

func filler(drive string) error {
	base := filepath.Join(drive, ".active_content_sandbox", "store", "resource", "cachedResources")

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	if err := os.MkdirAll(base, 0755); err != nil {
		return err
	}

	if err := os.Chdir(base); err != nil {
		return err
	}

	bar := progressbar.Default(5000, "Filling ")

	for i := 1; i <= 5000; i++ {
		name := strconv.Itoa(i)

		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				if err := os.Mkdir(name, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if err := os.Chdir(name); err != nil {
			return err
		}
		bar.Add(1)
	}

	return nil
}

func helper(src, destination string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copy(drive string) error {
	local := filepath.Join(".", ".active_content_sandbox")
	remote := filepath.Join(drive, ".active_content_sandbox")

	fmt.Printf("\nCopying New Store Cache...\n\n")
	if err := os.MkdirAll(remote, 0755); err != nil {
		return fmt.Errorf("Error: Failed To Create New Sandbox: %w", err)
	}

	return filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest, _ := filepath.Rel(local, path)
		target := filepath.Join(remote, dest)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return helper(path, target)
	})
}

func main() {
	defer func() {
		fmt.Printf("\n\nPress Enter to Continue...")
		fmt.Scanln()
	}()

	//Check For Cache
	var err error

	local := filepath.Join(".", ".active_content_sandbox")
	_, err = os.Stat(local)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: SpringBreak Payload Not Found! Ensure .active_content_sandbox Is In The Same Directory as This Executable.")
		} else {
			fmt.Printf("Error: %v", err)
		}

		return
	}

	fmt.Print("SpringBreak\n")
	fmt.Print("===========\n")

	if runtime.GOOS == "linux" {
		fmt.Print("\nPsst... We Detected You're on Linux! You May Have to Mount the Kindle in your File Explorer First for It to Be Detected.")
	}
	fmt.Print("\nSearching For Devices...")

	var drives []string
	for {
		drives, err = drivedetector.Detect()
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}

		if len(drives) > 0 {
			break
		}

		fmt.Print(".")
		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\nFound %d USB Device(s):\n", len(drives))
	for i, d := range drives {
		fmt.Printf("[%d] %s\n", i+1, d)
	}

	var choice int
	for {
		fmt.Print("\nSelect The Kindle Drive: ")
		_, err := fmt.Scanln(&choice)

		if err == nil && choice >= 1 && choice <= len(drives) {
			break
		}

		// stdin closed: bail out instead of looping forever
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			fmt.Fprintln(os.Stderr, "\nstdin closed; aborting.")
			os.Exit(1)
		}

		fmt.Println("Invalid Choice! Please Try Again.")
	}

	drive := drives[choice-1]
	fmt.Printf("\nKindle Detected On %s! Proceeding...\n", drive)

	_, err = os.Stat(filepath.Join(drive, "documents"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: Documents Folder Not Found On %s. Is This A Kindle Drive?\n", drive)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	fmt.Printf("Looks Like A Kindle! Starting...\n\n")

	path := filepath.Join(drive, ".active_content_sandbox", "store", "resource", "cachedResources", "1")
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Filler Files Detected! Deleting (Post-Jailbreak Cleanup)...\n\n")
		if err := cleanup(drive); err != nil {
			fmt.Printf("Cleanup Error: %v\n", err)
			return
		}

		fmt.Printf("Done! :) Have Fun With Your Jailbreak!")
		return
	}

	os.RemoveAll(filepath.Join(drive, ".active_content_sandbox"))

	if err := filler(drive); err != nil {
		fmt.Printf("Filler Error: %v\n", err)
		return
	}

	if err := copy(drive); err != nil {
		fmt.Printf("Copy Error: %v\n", err)
		return
	}

	fmt.Printf("Done! SpringBreak Preparation Complete. You Can Now Eject Your Kindle.\nOnce You Finish Jailbreaking, *Re-Run This Utility*.")
}
