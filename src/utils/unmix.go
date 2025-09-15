package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)




func UseDemucs(jobId string) error {
	//execution de la separation
	uploadFolder,err:=filepath.Abs("../storage/uploads")
	if err != nil {
		return fmt.Errorf("Erreur obtention du chemin courant:"+err.Error())
	}
	resultFolder,err:=filepath.Abs("../storage/separated")
	if err != nil {
		return fmt.Errorf("Erreur obtention du chemin courant: " + err.Error())
	}
	SetJobStatus("pending",jobId)
	cmd := exec.Command("docker", "run", "--rm" , "-v", fmt.Sprintf("%s:/app", uploadFolder), "-v", fmt.Sprintf("%s:/app/separated", resultFolder), "demucs-cpu", "demucs", "/app/"+jobId+".mp3")
	// Récupérer la sortie
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if (err != nil){
		return fmt.Errorf("erreur durant la separation:"+string(output))
	}
	//compression de des pistes

	err=CompressToZip(resultFolder+"/htdemucs","../storage/zipped",jobId)
	if err != nil{
		SetJobStatus("error",jobId)
	}
	SetJobStatus("donne",jobId)
	defer os.RemoveAll("../storage/separated/htdemucs/"+jobId)
	defer os.RemoveAll("../storage/uploads/"+jobId+".mp3")
	return nil
}
