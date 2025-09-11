package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	utils "stemExtractor/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(
		fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		BodyLimit:     10 * 1024 * 1024,
		ServerHeader:  "Fiber"})
	app.Get("/", func(c *fiber.Ctx) error {
		utils.CompressToZip("../storage/uploads","../storage/zipped","Baby I'm Yours - Breakbot")
		return c.SendString("Hello, World!")
	})
	app.Post("/unmix",func(c *fiber.Ctx)error{

		//reception et formattage du fichier
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).SendString(fmt.Sprintf("Erreur lecture fichier: %v", err))
		}
		base:= filepath.Base(fileHeader.Filename)
		stemDirName:= base[:len(base)-len(filepath.Ext(base))] //nom du fichier audio sans l'extention pour avoir le nom du repertoire dans lequel est stoquee les stem
	
		//sauvegarde dans le repertoire

		savePath := fmt.Sprintf("./%s", fileHeader.Filename)
		err = c.SaveFile(fileHeader, savePath)
		if (err!=nil){
			return c.Status(500).SendString("erreur pendant la sauvegarde: "+err.Error())
		}
		//execution de la separation
		pwd, err := os.Getwd()
		if err != nil {
			return c.Status(500).SendString("Erreur obtention du chemin courant: " + err.Error())
		}
		cmd := exec.Command("docker", "run", "-v", fmt.Sprintf("%s:/app", pwd), "demucs-cpu", "demucs", "/app/"+fileHeader.Filename)
		// Récupérer la sortie
		output, err := cmd.CombinedOutput()
		if (err != nil){
			fmt.Println(output)
			return c.Status(500).SendString("Erreur de la sortie")
		}
	
		//compression de des pistes

	return c.SendFile("./"+stemDirName+".zip")
	})
	app.Listen(":3000")
}
