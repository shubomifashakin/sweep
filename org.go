package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SweepInfo struct {
	TotalFiles int    `json:"total_files"`
	Date       string `json:"date"`
	Directory string `json:"directory"`
	Success int `json:"success"`
	SuccessfulMoves []MoveInfo `json:"successfulMoves"`
	Failures int `json:"failures"`
	FailedMoves []MoveInfo `json:"failedMoves"`
}

type MoveInfo struct {
	Source string `json:"source"`
	Destination string `json:"destination"`
}

func main(){
	var verbose bool
	var dir string
	
	flag.StringVar(&dir, "dir", "", "Directory to sweep")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	usersHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	storageDir:=filepath.Join(usersHomeDir,".sweep")

	err= os.MkdirAll(storageDir,0755)
	if err != nil {
		panic(err)
	}

	config:=zap.NewProductionConfig()

	if verbose {
		config.Level=zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	config.EncoderConfig.TimeKey="timestamp"
	config.EncoderConfig.EncodeTime=zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	if dir == "" {
		logger.Fatal("Please enter a directory")
	}

	dir=filepath.Clean(dir)

	 var infos SweepInfo

	 infos.Directory = dir

	 dirs:=[]string{"images", "videos", "documents", "others"}
	 logger.Debug("Creating directories", zap.Strings("dirs", dirs))
	 
	 for _, newDir := range dirs {
		created:=filepath.Join(dir,newDir)
		err = os.MkdirAll(created, 0755)

		if err != nil {
			logger.Fatal("Failed to create directory", zap.String("dir", created), zap.Error(err))
		}

		logger.Debug("Created directory", zap.String("dir", created))
	 }

	err=filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			logger.Debug("Skipping directory", zap.String("dir", d.Name()))
			return nil
		}

		if strings.Contains(path, filepath.Join(dir,"images"))  || 
		   strings.Contains(path, filepath.Join(dir,"videos"))  || 
		   strings.Contains(path, filepath.Join(dir,"documents"))  || 
		   strings.Contains(path, filepath.Join(dir,"others")) {
			logger.Debug("Skipping subdirectory", zap.String("dir", path))
			return nil
		}
		
		fileExt:= filepath.Ext(d.Name())

		if fileExt == "" {
			logger.Debug("Skipping file with no extension", zap.String("file", d.Name()))
			return nil
		}

		if fileExt == ".ini" {
			logger.Debug("Skipping ini file", zap.String("file", d.Name()))
			return nil
		}

		infos.TotalFiles++

		if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".gif" {

				handleMove(logger,dir,d,&infos,"images")

				return nil
		}

		if fileExt == ".mp4" || fileExt == ".avi" || fileExt == ".mkv" {
			

				handleMove(logger,dir,d,&infos,"videos")

				return nil
		}

		if fileExt == ".pdf" || fileExt == ".doc" || fileExt == ".docx" {
			

			handleMove(logger,dir,d,&infos,"documents")

			return nil
		}

		handleMove(logger,dir,d,&infos,"others")
		
		return nil
	})

	if err != nil {
		logger.Fatal("Failed to walk directory", zap.Error(err))
	} 
	
	infos.Date=time.Now().Format("2006-01-02 15:04:05")

	jsonData, err := json.Marshal(infos)
	if err != nil {
		logger.Fatal("Failed to marshal JSON", zap.Error(err))
	}

	reportPath := filepath.Join(storageDir, fmt.Sprintf("sweep-%s.json", time.Now().Format("2006-01-02-15-04-05")))
	err=os.WriteFile(reportPath, jsonData, 0644)
	if err != nil {
		logger.Fatal("Failed to write report",zap.Error(err))
	}else {
		logger.Info("Report saved", zap.String("path", reportPath))
	}
}

func handleMove(logger *zap.Logger,dir string, d fs.DirEntry, infos *SweepInfo, location string){
	logger.Debug("Moving file", zap.String("file", d.Name()))
		
				
	err := os.Rename(
			filepath.Join(dir, d.Name()),           
			filepath.Join(dir, location, d.Name()), 
				)

	if err != nil {
			logger.Error("Failed to move file", zap.String("file", d.Name()), zap.Error(err))

			infos.Failures++
			infos.FailedMoves = append(infos.FailedMoves, MoveInfo{
			Source: d.Name(),
			Destination: location + string(filepath.Separator) + d.Name(),
			})

			return 
	}	

			infos.Success++
			infos.SuccessfulMoves = append(infos.SuccessfulMoves, MoveInfo{
				Source: d.Name(),
				Destination: location + string(filepath.Separator) + d.Name(),
			})
		
}