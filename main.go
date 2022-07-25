package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/clfs/renew/internal"
)

func main() {
	log.SetFlags(0) // Disable log prefix.

	var (
		listFlag      = flag.Bool("list", false, "list go install-ed binaries")
		updateFlag    = flag.String("update", "", "update a single binary to latest")
		updateAllFlag = flag.Bool("update-all", false, "update all binaries to latest")
		skipFlag      = flag.Bool("skip", false, "skip binaries that fail to update")
	)
	flag.Parse()

	var err error
	switch {
	case *listFlag:
		err = list()
	case *updateFlag != "":
		err = update(*updateFlag, *skipFlag)
	case *updateAllFlag:
		err = updateAll(*skipFlag)
	default:
		flag.Usage()
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}

// list prints a list of all installed binaries.
func list() error {
	bins, err := internal.InstalledBinaries()
	if err != nil {
		return err
	}

	for _, bin := range bins {
		fmt.Printf("%s\t%s\n", bin.Name, bin.ImportPath)
	}
	return nil
}

// runUpdate is a helper function that runs updates.
func runUpdate(u *internal.Updater, bin internal.Binary, skip bool) error {
	fmt.Printf("==== %s\n", bin.Name)

	if err := u.Update(bin); err != nil {
		if skip {
			fmt.Printf("[-] skipped: %v\n", err)
			return nil
		}
		return err
	}

	fmt.Printf("[+] updated!\n")
	return nil
}

// update finds a binary by name and updates it to its latest version.
func update(name string, skip bool) error {
	bin, err := internal.BinaryFor(name)
	if err != nil {
		return err
	}

	return runUpdate(internal.NewUpdater(), bin, skip)
}

// updateAll updates all binaries to their latest versions.
func updateAll(skip bool) error {
	bins, err := internal.InstalledBinaries()
	if err != nil {
		return err
	}

	u := internal.NewUpdater()

	for _, bin := range bins {
		if err := runUpdate(u, bin, skip); err != nil {
			return err
		}
	}
	return nil
}
