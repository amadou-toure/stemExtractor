package main

import (
	"fmt"
	"os"

	"path/filepath"
	utils "stemExtractor/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)



func main() {
	app := fiber.New(
		fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		BodyLimit:     1024 * 1024 * 1024, //1GB
		ServerHeader:  "Fiber"})
	app.Get("/", func(c *fiber.Ctx) error {
		utils.CompressToZip("../storage/separated/htdemucs/Starship Syncopation - Cory Wong","../storage/zipped","starship-syncopation")
		return c.SendString("Hello, World!")
	})
	app.Get("/status/:id",func (c *fiber.Ctx) error {
		jobId:=c.Params("id")
		
		status:=utils.GetJobStatus(jobId)

		return c.Status(200).SendString(status)
	})
	

	app.Get("/download/:id", func(c *fiber.Ctx) error {
	    jobId := c.Params("id")
	    filePath := fmt.Sprintf("../storage/zipped/%s.zip", jobId)

	    f, err := os.Open(filePath)
	    if err != nil {
	        return c.Status(500).SendString("Erreur ouverture fichier: " + err.Error())
	    }

	    fileInfo, err := f.Stat()
	    if err != nil {
	        f.Close()
	        return c.Status(500).SendString("Erreur récupération info fichier: " + err.Error())
	    }

	    c.Set("Content-Type", "application/zip")
	    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", jobId))
	    c.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	    return c.SendStream(f)
	})
	

	app.Post("/unmix",func(c *fiber.Ctx)error{
		jobId:= uuid.New().String()
		//reception et formattage du fichier
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).SendString(fmt.Sprintf("Erreur lecture fichier: %v", err))
		}
		savePath := fmt.Sprintf("../storage/uploads/%s", jobId+filepath.Ext(fileHeader.Filename))
		err = c.SaveFile(fileHeader, savePath)
		if (err!=nil){
			return c.Status(500).SendString("erreur pendant la sauvegarde: "+err.Error())
		}
		
		go utils.UseDemucs(jobId)//lancement de l'unmixing
	return c.Status(200).SendString("job id:"+jobId)
		//sauvegarde dans le repertoire

		

	// return c.SendFile("../storage/zipped/"+stemDirName+".zip")
	})
	app.Listen(":3000")
}
