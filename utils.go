//  Copyright (c) 2019, salesforce.com, inc.
//  All rights reserved.
//  SPDX-License-Identifier: Apache License 2.0
//  For full license text, see the LICENSE file in the repo root

package main

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func processTar(tarBallReader *tar.Reader, cfg Config) (int, []File, []File, error) {
	var identifiedFiles []File
	var verifiedFiles []File
	numInstancesFound := 0
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[ERROR] Header Error", err)
			os.Exit(1)
		}
		switch header.Typeflag {
		case tar.TypeReg:
			bs, _ := ioutil.ReadAll(tarBallReader)
			if filepath.Ext(strings.TrimSpace(header.Name)) == ".tar" {
				r := bytes.NewReader(bs)
				tarBallReader2 := tar.NewReader(r)
				if len(cfg.Files) > 0 {
					tmpNumInstancesFound, tmpIdentifiedFiles, tmpVerifiedFiles, err := inspectTarForFiles(tarBallReader2, header.Name, cfg)
					if err != nil {
						return 0, identifiedFiles, verifiedFiles, err
					}
					numInstancesFound += tmpNumInstancesFound
					identifiedFiles = append(identifiedFiles, tmpIdentifiedFiles...)
					verifiedFiles = append(verifiedFiles, tmpVerifiedFiles...)
				}
				if len(cfg.Hashes) > 0 {
					inspectTarForHashes(tarBallReader2, header.Name, cfg)
				}
			}
		default:
			if cfg.VeryVerbose {
				fmt.Printf("Unable to untar type : '%c' in file %s\n", header.Typeflag, header.Name)
			}
		}
	}
	return numInstancesFound, identifiedFiles, verifiedFiles, nil
}

func inspectTarForFiles(tarBallReader *tar.Reader, tarName string, cfg Config) (int, []File, []File, error) {
	imageLayer := tarName[0:strings.Index(tarName, "/")]
	var identifiedFiles []File
	var verifiedFiles []File
	numInstancesFound := 0
	fmt.Println("[*] Inspecting Layer: ", imageLayer)

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[ERROR] Header Error", err)
			os.Exit(1)
		}
		filename := header.Name
		linkName := header.Linkname
		if cfg.VeryVerbose {
			fmt.Println("[*] Analysing File: ", filename)
		}

		if cfg.VeryVerbose {
			fmt.Println("[*] Linkname File: ", filename, linkName)
		}

		switch header.Typeflag {
		case tar.TypeReg:

			if cfg.VeryVerbose {
				fmt.Println("[*] File is TypeReg: ", filename)
			}

			numFilesToCheck := len(cfg.Files)
			if cfg.VeryVerbose {
				fmt.Println("[*] Number of Files to analyse:", numFilesToCheck)
			}

			for _, file := range cfg.Files {
				tarFilename := filename
				if strings.Count(filename, ".") >= 1 {
					if string(tarFilename[0]) == "." {
						tarFilename = strings.Replace(filename, ".", "", 1)
					}
				}
				if string(tarFilename[0]) != "/" {
					tarFilename = strings.ToLower("/" + tarFilename)
				}

				if strings.ToLower(file.Name) == tarFilename {

					if cfg.VeryVerbose {
						fmt.Println("[*] Found File in Tar:", tarFilename)
					}
					filePath := fmt.Sprintf("%s/%s", tarName[0:strings.Index(tarName, "/")], filename)

					if cfg.Verbose {
						fmt.Printf("\t[!] Identified  instance of '%s' at: %s \n", file.Name, filePath)
					}
					identifiedFiles = append(identifiedFiles, file)

					bs, err := ioutil.ReadAll(tarBallReader)
					if err != nil {
						log.Fatal("[ERROR]", err)
					}
					sum := sha256.Sum256(bs)
					hashVal := fmt.Sprintf("%x", sum)

					if cfg.VeryVerbose {
						fmt.Println("[*] Hash of Found File in Tar:", filename, hashVal)
					}

					for _, hash := range file.Hashes {

						if cfg.VeryVerbose {
							fmt.Println("[*] Checking if file has been verfied already:", file.Name)
						}

						if !checkIfFileInList(file.Name, verifiedFiles) {
							if cfg.Verbose {
								fmt.Println("[*] File has not been verified:", file.Name)
							}
							if strings.ToLower(hash.Hash) == strings.ToLower(hashVal) {
								if cfg.Verbose {
									fmt.Printf("[!] Found matching instance of '%s' at: %s with hash:%s\n", file.Name, filePath, hashVal)
								}
								numInstancesFound++
								verifiedFiles = append(verifiedFiles, file)
							}
						} else {
							if cfg.Verbose {
								fmt.Printf("[!] File has been verified already %s and this hash is probably duplicated:%s\n", file.Name, hash.Hash)
							}
						}
					}
				} else {
					if cfg.VeryVerbose {
						fmt.Println("[*] File not found in tar:", tarFilename)
					}
				}
			}
		case tar.TypeSymlink, tar.TypeLink:
			{
				for _, file := range cfg.Files {
					tarFilename := strings.Replace(filename, ".", "", 1)
					if string(tarFilename[0]) != "/" {
						tarFilename = strings.ToLower("/" + tarFilename)
					}
					linkFilename := strings.Replace(linkName, ".", "", 1)
					if string(linkFilename[0]) != "/" {
						linkFilename = strings.ToLower("/" + linkFilename)
					}
					if strings.ToLower(file.Name) == tarFilename {
						if !checkIfFileInList(linkFilename, cfg.Files) {
							if cfg.Verbose {
								fmt.Printf("[!] This file is a duplicate(%s), please add original component to ensure verification: %s\n", tarFilename, linkFilename)
							}
						} else {
							if cfg.Verbose {
								fmt.Printf("[!] This file is a duplicate(%s), original verification performed on: %s\n", tarFilename, linkFilename)
							}
							numInstancesFound++
							verifiedFiles = append(verifiedFiles, file)
						}
					}
				}
			}
		default:
			if cfg.VeryVerbose {
				fmt.Printf("Unable to untar type : %c in file %s\n", header.Typeflag, filename)
			}
		}
	}
	return numInstancesFound, identifiedFiles, verifiedFiles, nil
}

func inspectTarForHashes(tarBallReader *tar.Reader, tarName string, cfg Config) {
	imageLayer := tarName[0:strings.Index(tarName, "/")]
	fmt.Println("[*] Inspecting Layer: ", imageLayer)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		filename := header.Name
		if cfg.VeryVerbose {
			fmt.Println("[*] Analysing File: ", filename)
		}
		switch header.Typeflag {
		case tar.TypeReg:
			bs, err := ioutil.ReadAll(tarBallReader)
			if err != nil {
				log.Fatal("[ERROR]", err)
			}
			sum := sha256.Sum256(bs)
			hashVal := fmt.Sprintf("%x", sum)
			filePath := fmt.Sprintf("%s/%s", tarName[0:strings.Index(tarName, "/")], filename)
			for _, file := range cfg.Hashes {
				if strings.ToLower(file.Hash) == strings.ToLower(hashVal) {
					fmt.Printf("\t[!] Found file '%s' with hash: %s\n", filePath, hashVal)
				}
			}
		default:
			if cfg.VeryVerbose {
				fmt.Printf("Unable to untar type : %c in file %s\n", header.Typeflag, filename)
			}
		}
	}
}

func getOnDiskFile(sourcefile string) *os.File {
	file, err := os.Open(filepath.Clean(sourcefile))
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}
	return file
}

func inspectContainerForFiles(cfg Config) (int, []File, []File, error) {
	if cfg.VeryVerbose {
		fmt.Println("[+] Inspecting Container for Files")
	}
	var numInstancesFoundContainer = 0
	var identifiedContainerFiles []File
	var verifiedFiles []File
	if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
		return 0, identifiedContainerFiles, verifiedFiles, err
	}
	filepath.Walk(cfg.Path, func(path string, info os.FileInfo, err error) error {
		switch mode := info.Mode(); {

		case mode.IsRegular():

			for _, file := range cfg.Files {
				if cfg.VeryVerbose {
					fmt.Println("[+] Looking At: ", path)
				}

				if strings.ToLower(path) == strings.ToLower(cfg.Path+file.Name) {

					if cfg.VeryVerbose {
						fmt.Println("[+] Found: ", path)
					}
					identifiedContainerFiles = append(identifiedContainerFiles, file)

					fileToCheck := getOnDiskFile(path)
					fileToCheckBytes, err := ioutil.ReadAll(fileToCheck)
					if err != nil {
						fmt.Println("[ERROR]", err)
					}
					sum := sha256.Sum256(fileToCheckBytes)
					hashVal := fmt.Sprintf("%x", sum)

					if cfg.VeryVerbose {
						fmt.Println("[*] Hash of file: ", hashVal)
					}

					for _, hash := range file.Hashes {

						if cfg.VeryVerbose {
							fmt.Println("[*] Checking if file has been verfied already:", file.Name)
						}

						if !checkIfFileInList(file.Name, verifiedFiles) {
							if cfg.Verbose {
								fmt.Println("[*] File has not been verified:", file.Name)
							}
							if strings.ToLower(hash.Hash) == strings.ToLower(hashVal) {
								if cfg.Verbose {
									fmt.Printf("[!] Found matching instance of '%s' at: %s with hash:%s\n", file.Name, path, hashVal)
								}
								numInstancesFoundContainer++
								verifiedFiles = append(verifiedFiles, file)
							}
						} else {
							if cfg.Verbose {
								fmt.Printf("[!] File has been verified already %s and this hash is probably duplicated:%s\n", file.Name, hash.Hash)
							}
						}
					}
				}
			}
		default:
		}
		return nil
	})
	return numInstancesFoundContainer, identifiedContainerFiles, verifiedFiles, nil
}

func inspectContainerForHashes(cfg Config) {
	if cfg.VeryVerbose {
		fmt.Println("[+] Inspecting Container for Hashes")
	}

	filepath.Walk(cfg.Path, func(path string, info os.FileInfo, err error) error {
		switch mode := info.Mode(); {
		case mode.IsRegular():
			for _, hash := range cfg.Hashes {
				if cfg.VeryVerbose {
					fmt.Println("[+] Looking At: ", path)
				}
				fileToCheck := getOnDiskFile(path)
				fileToCheckBytes, err := ioutil.ReadAll(fileToCheck)
				if err != nil {
					fmt.Println("[ERROR]", err)
				}
				sum := sha256.Sum256(fileToCheckBytes)
				hashVal := fmt.Sprintf("%x", sum)

				if cfg.VeryVerbose {
					fmt.Println("[*] Hash of file: ", hashVal)
				}

				if strings.ToLower(hash.Hash) == strings.ToLower(hashVal) {
					if cfg.Verbose {
						fmt.Printf("[!] Found matching instance of '%s' at: %s with hash:%s\n", hash.Hash, path, hashVal)
					}
				}
			}
		default:
		}
		return nil
	})
}

func checkIfFileInList(fileToFind string, fileList []File) bool {
	for _, file := range fileList {
		if strings.ToLower(file.Name) == strings.ToLower(fileToFind) {
			return true
		}
	}
	return false
}
