package scheduler

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/sdk"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/snowflake"
)

func (scheduler *Scheduler) ImageSync() {
	scheduler.ScheduleAtFixedRate("dataSync", func() {
		log.Println("Data sync starting...")
		rootDir := os.Getenv("IMAGE_DIR")

		if rootDir == "" {
			log.Println("We dont have rootDir")
			return
		}

		files := []string{}
		err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error accessing %s: %v\n", path, err)
				return nil // continue traversing other files
			}

			// Get file information
			info, err := d.Info()
			if err != nil {
				log.Printf("Failed to get file information: %v\n", err)
				return nil
			}

			if d.IsDir() {
				log.Printf("📁 Directory: %s\n", path)
			} else {
				filename := d.Name()
				log.Printf("📄 File: %s (Size: %d bytes)\n", filename, info.Size())
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			log.Printf("Failed to traverse directory: %v\n", err)
		}

		if len(files) > 0 {
			var node, err = snowflake.NewNode(1)
			if err != nil {
				log.Fatal(err)
			}

			uploadfiles := []string{}
			images := make([]*db.Image, len(files))
			for i, f := range files {
				id := node.Generate()
				p := filepath.Dir(f)
				ext := filepath.Ext(f)
				newname := p + "/" + id.String() + ext
				if ext != "" {
					os.Rename(f, newname)
				}
				uploadfiles = append(uploadfiles, newname)
				image := &db.Image{
					Filename:       newname,
					OriginFilename: f,
					FileExtension:  ext,
				}
				images[i] = image
			}

			uploadrets := sdk.UploadMultipleFiles(uploadfiles)
			uploadmap := make(map[string]sdk.UploadResult)
			imagemap := make(map[string]*db.Image)

			for _, r := range uploadrets {
				uploadmap[r.Filename] = r
			}

			for _, i := range images {
				imagemap[i.Filename] = i
			}

			for _, i := range images {
				uploadResult := uploadmap[i.Filename]
				if uploadResult.Error != nil {
					i.Uploaded = false
				} else {
					i.Uploaded = true
				}
				scheduler.mediaService.CreateImage(i.Filename, i.OriginFilename, i.FileExtension, uploadResult.S3Bucket, uploadResult.S3Key, i.Uploaded)
			}

			// Clean directory
			for _, f := range uploadfiles {
				os.Remove(f)
			}
		}

		log.Println("Data sync finished...")
	}, 10*time.Second)
}
