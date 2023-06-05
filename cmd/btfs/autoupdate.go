// This auto update project linux, mac OS or windows operating system.
package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	btfs_version "github.com/bittorrent/go-btfs"
	"github.com/bittorrent/go-btfs-api"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	LatestConfigFile  = "config_%s_%s.yaml"
	CurrentConfigFile = "config.yaml"
)

var (
	Binaries = map[string]string{
		UpdateBinaryKey:           "update-%s-%s%s",
		BtfsBinaryKey:             "btfs-%s-%s%s",
		FsRepoMigrationsBinaryKey: "fs-repo-migrations-%s-%s%s",
	}
	UpdateBinaryKey           = "update"
	BtfsBinaryKey             = "btfs"
	FsRepoMigrationsBinaryKey = "fs-repo-migrations"
)

type Config struct {
	Version          string `yaml:"version"`
	Md5Check         string `yaml:"md5"`
	AutoupdateFlg    bool   `yaml:"autoupdateFlg"`
	SleepTimeSeconds int    `yaml:"sleepTimeSeconds"`
	BeginNumber      int    `yaml:"beginNumber"`
	EndNumber        int    `yaml:"endNumber"`
}

const (
	UPGRADE_FLAG_LATEST = iota
	UPGRADE_FLAG_OUT_OF_DATE
	UPGRADE_FLAG_SKIP
	UPGRADE_FLAG_ERR
)

type Repo struct {
	url        string
	compressed bool
}

// Auto update function.
func update(url, hval string) {
	// Get current program execution path.
	defaultBtfsPath, err := getCurrentPath()
	if err != nil {
		log.Errorf("Get current program execution path error, reasons: [%v]", err)
		return
	}

	configRepo := Repo{
		url:        "https://dist.btfs.io/release/",
		compressed: true,
	}

	var latestConfigFile, latestConfigPath, currentConfigPath string
	latestBinaryFiles := map[string]map[string]string{}

	// Select binary files based on operating system.
	if (runtime.GOOS == "darwin" || runtime.GOOS == "linux" || runtime.GOOS == "windows") &&
		(runtime.GOARCH == "amd64" || runtime.GOARCH == "386" ||
			runtime.GOARCH == "arm64" || runtime.GOARCH == "arm") {
		ext := ""
		sep := "/"
		compressedExt := ".tar.gz"
		if runtime.GOOS == "windows" {
			ext = ".exe"
			sep = "\\"
			compressedExt = ".zip"
		}

		latestConfigFile = fmt.Sprintf(LatestConfigFile, runtime.GOOS, runtime.GOARCH)
		latestConfigPath = fmt.Sprint(defaultBtfsPath, sep, latestConfigFile)
		currentConfigPath = fmt.Sprint(defaultBtfsPath, sep, CurrentConfigFile)
		for key, binFmt := range Binaries {
			bin := fmt.Sprintf(binFmt, runtime.GOOS, runtime.GOARCH, ext)
			binCompressed := fmt.Sprintf(binFmt, runtime.GOOS, runtime.GOARCH, compressedExt)
			path := fmt.Sprint(defaultBtfsPath, sep, bin)
			pathCompressed := fmt.Sprint(defaultBtfsPath, sep, binCompressed)
			pathMap := map[string]string{
				"bin":            bin,
				"binCompressed":  binCompressed,
				"path":           path,
				"pathCompressed": pathCompressed,
			}
			latestBinaryFiles[key] = pathMap
		}
	} else {
		log.Errorf("Operating system [%s], arch [%s] does not support automatic updates", runtime.GOOS, runtime.GOARCH)
		return
	}

	sleepTimeSeconds := 20
	version := "0.0.0"
	autoupdateFlg := true
	routePath := fmt.Sprint(runtime.GOOS, "/", runtime.GOARCH, "/")

updateLoop:
	for {
		time.Sleep(time.Second * time.Duration(sleepTimeSeconds))

		// Get current config file.
		currentConfig, err := getConfigure(currentConfigPath)
		if err == nil {
			autoupdateFlg = currentConfig.AutoupdateFlg
			version = currentConfig.Version
			sleepTimeSeconds = currentConfig.SleepTimeSeconds
		} else {
			version = btfs_version.CurrentVersionNumber
		}

		if !autoupdateFlg {
			continue
		}

		if pathExists(latestConfigPath) {
			// Delete the latest btfs config file.
			err = os.Remove(latestConfigPath)
			if err != nil {
				log.Errorf("Remove latest btfs config file error, reasons: [%v]", err)
				continue
			}
		}

		// Get latest btfs config file.
		err = download(latestConfigPath, fmt.Sprint(configRepo.url, routePath, latestConfigFile))
		if err != nil {
			log.Errorf("Download latest btfs config file error, reasons: [%v]", err)
			continue
		}

		// Parse latest config file.
		latestConfig, err := getConfigure(latestConfigPath)
		if err != nil {
			log.Errorf("Parse latest config file error, reasons: [%v]", err)
			continue
		}
		time.Sleep(time.Second * 5)

		// Where your local node is running on localhost:5001
		sh := shell.NewShell(url)

		// Get btfs id.
		idOutput, err := sh.ID()
		if err != nil {
			log.Errorf("Get btfs id error, reasons: [%v]", err)
			return
		}

		// beginNumber   endNumber         range
		//      0            0       no nodes updated
		//      0            1       [0, 1) 1% updated
		//      0           100      [0, 100)100% updated
		//     100          100      no nodes updated
		if convertStringToInt(idOutput.ID)%100 < latestConfig.BeginNumber || convertStringToInt(idOutput.ID)%100 >= latestConfig.EndNumber {
			fmt.Println("This node is not in the scope of this automatic update.")
			continue
		}

		// Compare version.
		upgradeFlg, err := versionCompare(latestConfig.Version, version)
		if err != nil {
			log.Errorf("Version compare error, reasons: [%v]", err)
			continue
		}

		if upgradeFlg == UPGRADE_FLAG_SKIP {
			fmt.Println("BTFS will not automatically update.")
			sleepTimeSeconds = latestConfig.SleepTimeSeconds
			continue
		}

		if upgradeFlg == UPGRADE_FLAG_LATEST {
			fmt.Println("BTFS is up-to-date.")
			sleepTimeSeconds = latestConfig.SleepTimeSeconds
			continue
		}

		fmt.Println("BTFS auto update begin.")

		for key, pathMap := range latestBinaryFiles {
			// Determine if the latest binary file exists.
			if pathExists(pathMap["path"]) {
				// Delete the btfs latest binary file.
				err = os.Remove(pathMap["path"])
				if err != nil {
					log.Errorf("Remove %s latest binary file error, reasons: [%v]", key, err)
					continue updateLoop
				}
			}

			if configRepo.compressed {
				// Delete the newly downloaded config file since the compressed file contains the config file
				// This is only included as part of btfs binary compressed
				if key == BtfsBinaryKey && pathExists(latestConfigPath) {
					// Delete the latest btfs config file.
					err = os.Remove(latestConfigPath)
					if err != nil {
						log.Errorf("Remove latest btfs config file error, reasons: [%v]", err)
						continue updateLoop
					}
				}

				// Get the latest compressed file.
				err = download(pathMap["pathCompressed"], fmt.Sprint(configRepo.url, routePath,
					pathMap["binCompressed"]))
				if err != nil {
					log.Errorf("Download %s latest compressed file error, reasons: [%v]", key, err)
					continue updateLoop
				}

				fmt.Printf("Download %s latest compressed file success!\n", key)

				// Unarchive the tar.gz or .zip binary file
				err = archiver.Unarchive(pathMap["pathCompressed"], defaultBtfsPath)
				if err != nil {
					log.Errorf("Unarchive of %s latest binary file error, reasons: [%v]", key, err)
					continue updateLoop
				}

				// Delete the archive file.
				err = removeCompressedFile(pathMap["pathCompressed"])
				if err != nil {
					log.Errorf("Remove %s latest compressed file error, reasons: [%v]", key, err)
					continue updateLoop
				}

				// Verify config file is present after unarchiving compressed file
				// This is only included as part of btfs binary compressed
				if key == BtfsBinaryKey && !pathExists(latestConfigPath) {
					// Re-download the config file
					err = download(latestConfigPath, fmt.Sprint(configRepo.url, routePath, latestConfigFile))
					if err != nil {
						log.Errorf("Re-download latest btfs config file error, reasons: [%v]", err)
						continue updateLoop
					}
				}

			} else {
				err = download(pathMap["path"], fmt.Sprint(configRepo.url, routePath, pathMap["bin"]))
				if err != nil {
					log.Errorf("Download %s latest binary file error, reasons: [%v]", key, err)
					continue updateLoop
				}
				fmt.Println("Download btfs binary file success!")
			}

			// MD5 check only availlable for btfs binary
			// Md5 encode file.
			if key == BtfsBinaryKey {
				latestMd5Hash, err := md5Encode(pathMap["path"])
				if err != nil {
					log.Errorf("MD5 encode %s latest binary file error, reasons: [%v].", key, err)
					continue updateLoop
				}

				if latestMd5Hash != latestConfig.Md5Check {
					log.Errorf("MD5 verify %s latest binary file failed.", key)
					continue updateLoop
				}

				fmt.Printf("MD5 check %s binary file success!\n", key)
			}
		}

		// Add executable permissions to update binary file.
		err = os.Chmod(latestBinaryFiles[UpdateBinaryKey]["path"], 0775)
		if err != nil {
			log.Errorf("Chmod latest update binary file 775 error, reasons: [%v]", err)
			continue
		}

		// Start the btfs-updater binary process.
		cmd := exec.Command(latestBinaryFiles[UpdateBinaryKey]["path"],
			"-url", url,
			"-project", defaultBtfsPath+string(os.PathSeparator),
			"-download", defaultBtfsPath+string(os.PathSeparator),
			"-hval", hval,
		)
		err = cmd.Start()
		if err != nil {
			log.Errorf("Update start failed: [%v]", err)
			continue
		}

		fmt.Println("Process will exit now and restart after the update completes.")
		os.Exit(0)
	}
}

// Determine if the path file exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Get config struct from yaml file.
func getConfigure(fileName string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	conf := new(Config)
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Compare version.
func versionCompare(version1, version2 string) (int, error) {
	// Split string of version1.
	s1 := strings.Split(version1, ".")
	if s1 == nil || len(s1) != 3 {
		return UPGRADE_FLAG_ERR, errors.New("String for major version has wrong format")
	}

	// Split string of version2.
	s2 := strings.Split(version2, ".")
	if s2 == nil || len(s2) != 3 {
		return UPGRADE_FLAG_ERR, errors.New("String for minor version has wrong format")
	}

	// If the current config.yaml contains a dash in the last section
	// then do not automatic update.
	if strings.Contains(s2[2], "-") {
		fmt.Println("BTFS config.yaml shows a dev version.")
		return UPGRADE_FLAG_SKIP, nil
	}
	// If the newly downloaded config.yaml contains a dash in the last section
	// then do not automatic update.
	if strings.Contains(s1[2], "-") {
		fmt.Println("BTFS upgrade config.yaml shows a dev version.")
		return UPGRADE_FLAG_SKIP, nil
	}

	for i := 0; i < 3; i++ {
		// Convert version1 from string to int.
		int1, err := strconv.Atoi(s1[i])
		if err != nil {
			log.Errorf("Convert version1 from string to int error, reasons: [%v]", err)
			return UPGRADE_FLAG_ERR, err
		}

		// Convert version2 from string to int.
		int2, err := strconv.Atoi(s2[i])
		if err != nil {
			log.Errorf("Convert version2 from string to int error, reasons: [%v]", err)
			return UPGRADE_FLAG_ERR, err
		}

		if int1 > int2 {
			return UPGRADE_FLAG_OUT_OF_DATE, nil
		} else if int1 < int2 {
			return UPGRADE_FLAG_LATEST, nil
		}
	}
	return UPGRADE_FLAG_LATEST, nil
}

// Get current program execution path.
func getCurrentPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// http get download function.
func download(downloadPath, url string) error {
	// http get.
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("Http get error, reasons: [%v]", err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Download failed with %d, message: %s", res.StatusCode, res.Status)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	// Create file on local.
	f, err := os.Create(downloadPath)
	if err != nil {
		log.Errorf("Create file error, reasons: [%v]", err)
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	// Copy file from response body to local file.
	written, err := io.Copy(f, res.Body)
	if err != nil {
		log.Errorf("Copy file error, reasons: [%v]", err)
		return err
	}

	log.Infof("Download success, file size: [%f]MiB", float32(written)/(1024*1024))
	return nil
}

// Removes the archive file from disk.
func removeCompressedFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		log.Errorf("Remove gz file error, reasons: [%v]", err)
		return err
	}
	return nil
}

// Md5 encode file by file path.
func md5Encode(name string) (string, error) {
	// Open file.
	file, err := os.Open(name)
	if err != nil {

		log.Errorf("Open file failed, reasons: [%v]", err)
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	// New reader of file.
	buffer := bufio.NewReader(file)

	// New md5.
	md5Hash := md5.New()

	// Copy file stream to md5 hash.
	_, err = io.Copy(md5Hash, buffer)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

// Convert string to int.
func convertStringToInt(s string) int {
	sum := 0
	for _, v := range []byte(s) {
		sum += int(v)
	}
	return sum
}
