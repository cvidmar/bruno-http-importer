package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/bruno-importer/importer"
)

var VERSION = "v0.0.1-alpha"

var Usage = func() {
	fmt.Printf("\nSYNOPSIS\n")
	fmt.Printf("     go run main.go -i INPUTDIR -o OUTPUTDIR \n\n")
	flag.PrintDefaults()
	fmt.Printf("\nCurrent limitations:\n")
	fmt.Println("- Only migrates POST and GET methods")
	fmt.Println("- POST requests migrations only support JSON body")
	fmt.Println()
}

func usageError(comment string) {
	fmt.Println("ERROR:", comment)
	Usage()
	os.Exit(1)
}

func fatal(infos ...interface{}) {
	fmt.Println(infos...)
	os.Exit(1)
}

func main() {
	// Get parameters from cmd line flags
	flagInputDir := flag.String("i", "", "Full path to input directory")
	flagOutputDir := flag.String("o", "", "Full path to output directory")
	flagVersion := flag.Bool("version", false, "Print version and exit")
	flagHelp := flag.Bool("help", false, "Print help and exit")
	flag.Parse()

	if *flagVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *flagHelp {
		Usage()
		os.Exit(0)
	}

	if *flagInputDir == "" {
		usageError("Missing input directory")
	}

	if *flagOutputDir == "" {
		usageError("Missing output directory")
	}

	err := importer.WalkDir(*flagInputDir, *flagOutputDir)
	if err != nil && err != io.EOF {
		fatal(err)
	}
	fmt.Println("All done!")
}
