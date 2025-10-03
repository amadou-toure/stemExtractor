package main

import (
	"fmt"
	"os"

	"path/filepath"
	utils "stemExtractor/utils"

	//models "stemExtractor/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8081, http://127.0.0.1:8081",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/status/:id", func(c *fiber.Ctx) error {
		jobId := c.Params("id")

		status := utils.GetJobStatus(jobId)

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
		fmt.Println("remote host started a download")
		return c.Status(200).SendStream(f)
		//return c.Status(200).SendFile(filePath)
	})

	app.Post("/unmix", func(c *fiber.Ctx) error {
		jobId := uuid.New().String()
		//reception et formattage du fichier
		fileHeader, err := c.FormFile("file")
		if err != nil {
			fmt.Println("Erreur lecture fichier:", err)
			return c.Status(400).SendString("Erreur lecture fichier:" + err.Error())
		}
		savePath := fmt.Sprintf("../storage/uploads/%s", jobId+filepath.Ext(fileHeader.Filename))
		err = c.SaveFile(fileHeader, savePath)
		if err != nil {
			return c.Status(500).SendString("erreur pendant la sauvegarde: " + err.Error())
		}

		go utils.UseDemucs(jobId) //lancement de l'unmixing
		return c.Status(200).SendString(jobId)
	})
	app.Listen(":3000")
}
