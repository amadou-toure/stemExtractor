package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)
 func CompressToZip(sourceDir string,zipPath string, fileName string) error {

	zipFile,err:= os.Create(zipPath+"/"+fileName+".zip")
	if err != nil {
		return fmt.Errorf("Erreur création fichier zip: "+err.Error())
	}
	defer zipFile.Close()
	
	zipWriter:= zip.NewWriter(zipFile)
	defer zipWriter.Close()

	filepath.Walk(sourceDir,func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir(){
			return nil
		}
		relPath, _ := filepath.Rel(sourceDir, path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, f)
		return err
	})
	return nil
 }