package path

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	cmds "github.com/bittorrent/go-btfs-cmds"

	"github.com/dustin/go-humanize"
	logging "github.com/ipfs/go-log"
	"github.com/mitchellh/go-homedir"
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	defaultPath  = "~/.btfs"
	properties   = ".btfs.properties"
	BtfsPathKey  = "BTFS_PATH"
	BtfsPathFlag = "BTFS_PATH_FLAG"
)

var (
	PropertiesFileName string
	srcProperties      string
)

/*
	can be dir of `btfs` or path like `/private/var/folders/q0/lc8cmwd93gv50ygrsy3bwfyc0000gn/T`,

depends on how `btfs` is called
*/
func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Error("err", err)
		return
	}
	exPath := filepath.Dir(ex)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("err", err)
		return
	}
	srcProperties = filepath.Join(home, properties)
	PropertiesFileName = filepath.Join(exPath, properties)
	// .btfs.properties migration
	if !CheckExist(PropertiesFileName) && CheckExist(srcProperties) {
		if err := copyFile(srcProperties, PropertiesFileName); err != nil {
			log.Errorf("error occurred when copy .btfs.properties", err)
			return
		}
		err := os.Remove(srcProperties)
		if err != nil {
			log.Errorf("error occurred when remove %s", srcProperties)
		}
	}
	SetEnvVariables()
}

var Executable = func() string {
	if ex, err := os.Executable(); err == nil {
		return ex
	}
	return "btfs"
}()

var log = logging.Logger("core/commands/path")

var (
	btfsPath   string
	StorePath  string
	OriginPath string
	lock       Mutex
)

const mutexLocked = 1 << iota

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}

var PathCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Modify the Host storage folder path for BTFS client.",
		ShortDescription: `
The default local repository path is located at ~/.btfs folder, in order to
improve the hard disk space usage, provide the function to change the original 
storage location, a specified path as a parameter need to be passed.
`,
	},
	Subcommands: map[string]*cmds.Command{
		"status":   PathStatusCmd,
		"capacity": PathCapacityCmd,
		"migrate":  PathMigrateCmd,
		"list":     PathListCmd,
		"mkdir":    PathMkdirCmd,
		"volumes":  PathVolumesCmd,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("path-name", true, false,
			"New BTFS Path.Should be absolute path."),
		cmds.StringArg("storage-size", true, false, "Storage Commitment Size"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		locked := lock.TryLock()
		if locked {
			defer lock.Unlock()
		} else {
			return errors.New("Cannot set path concurrently.")
		}
		StorePath = strings.Trim(req.Arguments[0], " ")

		if StorePath == "" {
			return fmt.Errorf("path is not defined")
		}
		var err error
		if StorePath, err = homedir.Expand(StorePath); err != nil {
			return err
		}
		if !filepath.IsAbs(StorePath) {
			StorePath, err = filepath.Abs(StorePath)
			if err != nil {
				return err
			}
		}
		if btfsPath != "" {
			if btfsPath != StorePath {
				OriginPath = btfsPath
			} else {
				return fmt.Errorf("specifed path is same with current path")
			}
		} else if envBtfsPath := os.Getenv(BtfsPathKey); envBtfsPath != "" {
			OriginPath = envBtfsPath
		} else if home, err := homedir.Expand(defaultPath); err == nil && home != "" {
			OriginPath = home
		} else {
			return fmt.Errorf("can not find the original stored path")
		}

		if err := validatePath(OriginPath, StorePath); err != nil {
			return err
		}

		usage, err := disk.UsageWithContext(req.Context, filepath.Dir(StorePath))
		if err != nil {
			return err
		}
		promisedStorageSize, err := humanize.ParseBytes(req.Arguments[1])
		if err != nil {
			return err
		}
		if usage.Free < promisedStorageSize {
			return fmt.Errorf("Not enough disk space, expect: ge %v bytes, actual: %v bytes",
				promisedStorageSize, usage.Free)
		}
		go func() {
			if e := os.Unsetenv(BtfsPathFlag); e != nil {
				log.Error(e)
			}
			DoRestart(true)
		}()
		return nil
	},
}

func DoRestart(p bool) {
	time.Sleep(2 * time.Second)
	restartCmd := exec.Command(Executable, "restart", fmt.Sprintf("-p=%v", p))
	if err := restartCmd.Run(); err != nil {
		log.Errorf("restart error, %v", err)
	}
}

var PathStatusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Get status of resetting path.",
		ShortDescription: "Get status of resetting path.",
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		tryLock := lock.TryLock()
		if tryLock {
			lock.Unlock()
		}
		return cmds.EmitOnce(res, PathStatus{
			Resetting: !tryLock,
			Path:      StorePath,
		})
	},
	Type: PathStatus{},
}

type PathStatus struct {
	Resetting bool
	Path      string
}

var PathCapacityCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Get free space of passed path.",
		ShortDescription: "Get free space of passed path.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("path-name", true, true,
			"New BTFS Path. Should be absolute path."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		path := strings.Trim(req.Arguments[0], " ")
		if path == "" {
			return fmt.Errorf("path is not defined")
		}
		var err error
		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
		}
		if btfsPath != "" {
			if btfsPath != StorePath {
				OriginPath = btfsPath
			} else {
				return fmt.Errorf("specifed path is same with current path")
			}
		} else if envBtfsPath := os.Getenv(BtfsPathKey); envBtfsPath != "" {
			OriginPath = envBtfsPath
		} else if home, err := homedir.Expand(defaultPath); err == nil && home != "" {
			OriginPath = home
		} else {
			return fmt.Errorf("can not find the original stored path")
		}
		if err := validatePath(OriginPath, path); err != nil {
			return err
		}
		if !CheckDirEmpty(path) {
			return fmt.Errorf("path %s is not empty", path)
		}
		valid := true
		usage, err := disk.UsageWithContext(req.Context, filepath.Dir(path))
		if err != nil {
			return err
		}
		humanizedFreeSpace := humanize.Bytes(usage.Free)
		return cmds.EmitOnce(res, &PathCapacity{
			FreeSpace:          usage.Free,
			Valid:              valid,
			HumanizedFreeSpace: humanizedFreeSpace,
		})
	},
	Type: &PathCapacity{},
}

var PathMigrateCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "path migrate. e.x.: btfs storage path migrate /Users/tron/.btfs.new",
		ShortDescription: "path migrate.",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("btfs-dir", true, true,
			"Current BTFS Path. Should be absolute path."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		if _, k := os.LookupEnv(BtfsPathKey); k || CheckExist(srcProperties) || CheckExist(PropertiesFileName) {
			return errors.New("no need to migrate")
		}
		fmt.Printf("Writing \"%s\" to %s\n", req.Arguments[0], PropertiesFileName)
		return ioutil.WriteFile(PropertiesFileName, []byte(req.Arguments[0]), os.ModePerm)
	},
}

var PathListCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "list directories",
		ShortDescription: "list directories",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("parent", true, true,
			"parent path, should be absolute path."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		list, err := list(req.Arguments[0])
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, stringList{Strings: list})
	},
	Type: stringList{},
}

var PathVolumesCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "list disk volumes",
		ShortDescription: "list disk volumes",
	},
	Arguments: []cmds.Argument{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		ps, err := volumes()
		if err != nil {
			return err
		}
		return cmds.EmitOnce(res, ps)
	},
	Type: []volume{},
}

type volume struct {
	Name       string `json:"name"`
	MountPoint string `json:"mount_point"`
}

type stringList struct {
	Strings []string
}

var PathMkdirCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "add folder",
		ShortDescription: "add folder",
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("parent", true, false,
			"parent path, should be absolute path."),
		cmds.StringArg("name", true, false,
			"path name"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return add(req.Arguments[0], req.Arguments[1])
	},
}

func validatePath(src string, dest string) error {
	log.Debug("src", src, "dest", dest)
	// clean: /abc/ => /abc
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)
	if src == dest || strings.HasPrefix(dest, src+string(filepath.Separator)) {
		return errors.New("invalid path")
	}
	dir := filepath.Dir(src)
	if !CheckExist(src) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("mkdir: %s", err)
		}
	}
	return nil
}

type PathCapacity struct {
	FreeSpace          uint64
	Valid              bool
	HumanizedFreeSpace string
}

func WriteProperties() error {
	if !CheckExist(PropertiesFileName) {
		newFile, err := os.Create(PropertiesFileName)
		defer newFile.Close()
		if err != nil {
			return err
		}
	}
	data := []byte(StorePath)
	err := ioutil.WriteFile(PropertiesFileName, data, 0666)
	if err == nil {
		fmt.Printf("Storage location was reset in %v\n", StorePath)
	}
	return err
}

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

func MoveFolder() error {
	err := os.Rename(OriginPath, StorePath)
	// src and dest dir are not in the same partition
	if err != nil {
		err := move(OriginPath, StorePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func move(src string, dst string) error {
	if err := copyDir(src, dst); err != nil {
		return err
	}
	if err := os.RemoveAll(src); err != nil {
		fmt.Println("remove old path failed: %w", err)
	}
	return nil
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func ReadProperties(filePath string) string {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Read properties fail: [%v]\n", err)
	}
	return string(f)
}

func CheckDirEmpty(dirname string) bool {
	dir, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Debug("Read directory fail: [%v]\n", err)
	}
	return len(dir) == 0
}

func SetEnvVariables() {
	if CheckExist(PropertiesFileName) {
		btfsPath = ReadProperties(PropertiesFileName)
		btfsPath = strings.Trim(btfsPath, " \n\r")
		if btfsPath != "" {
			newPath := btfsPath
			_, b := os.LookupEnv(BtfsPathFlag)
			if !b {
				err := os.Setenv(BtfsPathKey, newPath)
				if err != nil {
					log.Errorf("cannot set env variable of BTFS_PATH: [%v] \n", err)
				} else {
					if err := os.Setenv(BtfsPathFlag, "T"); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

func CheckExist(pathName string) bool {
	_, err := os.Stat(pathName)
	return !os.IsNotExist(err)
}
