package commands

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	commands "github.com/bittorrent/go-btfs/commands"
	fsrepo "github.com/bittorrent/go-btfs/repo/fsrepo"
)

const (
	outputFileOption = "o"
	compressOption   = "a"
	backupPathOption = "r"
	excludeOption    = "exclude"
)

var BackupCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Back up BTFS's data",
		LongDescription: `
This command will create a backup of the data from the current BTFS node.
`,
	},
	Arguments: []cmds.Argument{
		cmds.FileArg("file", true, false, "data to encode").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(outputFileOption, "backup output file path"),
		cmds.StringOption(compressOption, "gz or zip").WithDefault("gz"),
		cmds.StringsOption(excludeOption, "exclude backup output file path"),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		r, err := fsrepo.Open(env.(*commands.Context).ConfigRoot)
		if err != nil {
			return err
		}
		defer r.Close()

		var fileName = fmt.Sprintf("btfs_backup_%d", time.Now().Unix())

		outputName, ok := req.Options[outputFileOption].(string)
		if ok {
			fileName = outputName
		}
		btfsPath, err := fsrepo.BestKnownPath()
		if err != nil {
			return err
		}

		excludePath, _ := req.Options[excludeOption].([]string)
		for _, v := range excludePath {
			// TODO
			if v != "config" && v != "statestore" && v != "datastore" {
				return errors.New("-exclude only support config, statestore or datastore")
			}
		}
		// exclude the repo.lock to avoid dead lock
		excludePath = append(excludePath, "repo.lock")
		compressWay, _ := req.Options[compressOption].(string)
		// TODO
		if compressWay != "gz" && compressWay != "zip" {
			return errors.New("-a only support zip or gz, gz is default")
		}
		absPath, err := filepath.Abs(fileName)
		if err != nil {
			return err
		}
		if compressWay == "zip" {
			absPath += ".zip"
			err = Zip(btfsPath, absPath, excludePath)
		} else {
			absPath += ".tar.gz"
			err = Tar(btfsPath, absPath, excludePath)
		}
		if err != nil {
			return err
		}
		fmt.Printf("Backup successful! The backup path is %s\n", absPath)
		return nil
	},
}

var RecoveryCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:         "Recover BTFS's data from a archived file of backup",
		LongDescription: `This command will recover data from a previously created backup file`,
	},
	Options: []cmds.Option{
		cmds.StringOption(backupPathOption, "backup output file path"),
	},
	Run: func(req *cmds.Request, resp cmds.ResponseEmitter, env cmds.Environment) error {
		backupPath, ok := req.Options[backupPathOption].(string)
		if !ok {
			return errors.New("you need to specify -r to indicate the path you want to recover")
		}
		btfsPath := env.(*commands.Context).ConfigRoot
		dstPath := filepath.Dir(btfsPath)
		if fsrepo.IsInitialized(btfsPath) {
			newPath := filepath.Join(dstPath, fmt.Sprintf(".btfs_backup_%d", time.Now().Unix()))
			// newPath := filepath.Join(filepath.Dir(btfsPath), backup)
			err := os.Rename(btfsPath, newPath)
			if err != nil {
				return err
			}
			fmt.Println("btfs configuration file already exists!")
			fmt.Println("We have renamed it to ", newPath)
		}

		if err := UnTar(backupPath, dstPath); err != nil {
			err = UnZip(backupPath, dstPath)
			if err != nil {
				return errors.New("your file is not exists or your file format is not tar.gz or zip, please check again")
			}
		}
		fmt.Println("Recovery successful!")
		return nil
	},
}

func Tar(src, dst string, excludePath []string) (err error) {
	fw, err := os.Create(dst)
	if err != nil {
		return
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	basePath := filepath.Dir(src)
	filepath.Walk(src, func(fileAbsPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, v := range excludePath {
			excludeAbsPath := filepath.Join(src, v)
			if v != "" && strings.HasPrefix(fileAbsPath, excludeAbsPath) {
				return nil
			}
		}
		rel, err := filepath.Rel(basePath, fileAbsPath)
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		hdr.Name = rel

		// 写入文件信息
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		fr, err := os.Open(fileAbsPath)
		if err != nil {
			return err
		}
		defer fr.Close()

		// copy 文件数据到 tw
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		return nil
	})
	tw.Flush()
	gw.Flush()
	return
}

func UnTar(src, dst string) (err error) {
	fr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	// tar read
	tr := tar.NewReader(gr)
	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.FileInfo().IsDir() {
			err = os.MkdirAll(filepath.Join(dst, h.Name), h.FileInfo().Mode())
			if err != nil {
				return err
			}
			continue
		}

		fw, err := os.OpenFile(filepath.Join(dst, h.Name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(h.Mode))
		if err != nil {
			return err
		}
		defer fw.Close()
		// 写文件
		_, err = io.Copy(fw, tr)
		if err != nil {
			return err
		}
	}
	return
}

func Zip(src, dst string, excludePath []string) (err error) {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(fw)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	basePath := filepath.Dir(src)
	filepath.Walk(src, func(fileAbsPath string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}
		for _, v := range excludePath {
			excludeAbsPath := filepath.Join(src, v)
			if v != "" && strings.HasPrefix(fileAbsPath, excludeAbsPath) {
				return nil
			}
		}
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		rel, err := filepath.Rel(basePath, fileAbsPath)
		if err != nil {
			return err
		}
		fh.Name = rel

		if fi.IsDir() {
			fh.Name += "/"
		}

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(fileAbsPath)
		if err != nil {
			return
		}
		defer fr.Close()

		_, err = io.Copy(w, fr)
		if err != nil {
			return
		}
		return nil
	})
	zw.Flush()
	return
}

func UnZip(src, dst string) (err error) {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer zr.Close()

	for _, file := range zr.File {
		err = persistZipFile(dst, file)
		if err != nil {
			return
		}
	}
	return nil
}

func persistZipFile(dst string, file *zip.File) (err error) {
	path := filepath.Join(dst, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, file.Mode())
	}

	fr, err := file.Open()
	if err != nil {
		return err
	}
	defer fr.Close()

	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	if err != nil {
		return err
	}
	return
}
