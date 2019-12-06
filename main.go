//  Copyright (c) 2019, salesforce.com, inc.
//  All rights reserved.
//  SPDX-License-Identifier: Apache License 2.0
//  For full license text, see the LICENSE file in the repo root

package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Hash struct {
	Hash string
}

type File struct {
	Name   string
	Hashes []Hash
}

type Config struct {
	Path        string
	Image       string
	Mode        string
	Verbose     bool
	VeryVerbose bool
	Overwrite   bool
	Files       []File
	Hashes      []Hash
}

var cfgPtr *string

func main() {
	cfgPtr = flag.String("cfg", "cfg.yml", "Load config from provided yaml file")
	flag.Parse()

	_, err := os.Stat(*cfgPtr)

	var cfg Config
	if err != nil {
		fmt.Println("[+] Could not load config from:", *cfgPtr)
		os.Exit(1)
	} else {
		fmt.Println("[+] Loading config: ", *cfgPtr)
		cfg, err = loadConfig()
		if err != nil {
			fmt.Println("[ERROR] ", err)
			os.Exit(1)
		}

		switch cfg.Mode {
		case "image":
			fmt.Println("[+] Analysing Image")
			if len(cfg.Files) > 0 || len(cfg.Hashes) > 0 {
				doImageAnalysis(cfg)
			} else {
				fmt.Println("[!] Please provide list of files to analyse")
			}
		case "container":
			fmt.Println("[+] Analysing Container")
			if len(cfg.Files) > 0 || len(cfg.Hashes) > 0 {
				doContainerAnalysis(cfg)
			} else {
				fmt.Println("[!] Please provide list of files to analyse")
				os.Exit(1)
			}
		default:
			fmt.Println("[!] Please specify a mode")
		}
	}
}

func doContainerAnalysis(cfg Config) {
	instancesFound, filesFound, err := inspectContainerForFiles(cfg)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		os.Exit(1)
	}

	if len(cfg.Files) > 0 {
		if instancesFound < len(cfg.Files) {
			fmt.Printf("[!] Not all files were identified: (%d/%d)\n", instancesFound, len(cfg.Files))
			for _, file := range cfg.Files {
				if !checkIfFileInList(file.Name, filesFound) {
					fmt.Println("[!] File not verified: ", file)
				}
			}
			os.Exit(1)
		} else {
			fmt.Printf("[!] All files were identified and verified: (%d/%d)\n", instancesFound, len(cfg.Files))
		}
	}

}

func startContainerAnalysis(cfg Config) (int, []File, []File, error) {
	fmt.Println("[+] Docker Image Source: ", cfg.Image)
	file := getOnDiskFile(cfg.Image)
	defer file.Close()
	var fileReader io.ReadCloser = file
	tarBallReader := tar.NewReader(fileReader)
	return processTar(tarBallReader, cfg)
}

func doImageAnalysis(cfg Config) {

	numInstancesFound := 0
	var identifiedFiles []File
	var verifiedFiles []File
	var err error
	numInstancesFound, identifiedFiles, verifiedFiles, err = startContainerAnalysis(cfg)
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}

	if len(cfg.Files) > 0 {

		if len(cfg.Files) != len(identifiedFiles) {
			numNotIdentifiedFiles := len(cfg.Files) - len(identifiedFiles)
			fmt.Printf("[!] Not all components were identifed: (%d/%d)\n", numNotIdentifiedFiles, len(cfg.Files))
			for _, file := range cfg.Files {
				if !checkIfFileInList(file.Name, identifiedFiles) {
					fmt.Println("[!] Component not identified: ", file.Name)
				}
			}
			os.Exit(1)
		} else {
			fmt.Printf("[!] All components were identified: (%d/%d)\n", len(identifiedFiles), len(cfg.Files))
		}

		if len(cfg.Files) != len(verifiedFiles) {
			fmt.Printf("[!] Not all components were verified: (%d/%d)\n", len(verifiedFiles), len(cfg.Files))
			for _, file := range identifiedFiles {
				if !checkIfFileInList(file.Name, verifiedFiles) {
					fmt.Println("[!] Component not verified: ", file.Name)
				}
			}
			os.Exit(1)
		} else {
			fmt.Printf("[!] All components were identified and verified: (%d/%d)\n", numInstancesFound, len(cfg.Files))
		}
	}
}

func loadConfig() (Config, error) {
	var config Config
	source, err := ioutil.ReadFile("cfg.yml")
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
