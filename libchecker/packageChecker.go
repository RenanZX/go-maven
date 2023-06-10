package libchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	. "github.com/RenanZX/go-maven/types"
	"github.com/RenanZX/go-maven/ziper"
)

func HttpGetJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func SolSearch(Pacote string) *PackageMaven {
	r := new(ResponseServer)
	query := fmt.Sprintf("https://search.maven.org/solrsearch/select?q=%s&rows=1&wt=json", Pacote)
	HttpGetJSON(query, r)
	docs := r.Response.Docs
	if len(docs) > 0 {
		//fmt.Println(docs[0])
		pacote := PackageMaven{
			GroupId:    docs[0].G,
			ArtifactId: docs[0].A,
			Version:    docs[0].LatestVersion,
		}
		// fmt.Println(pacote)
		return &pacote
	}
	return nil
}

func checkOutMainFestPack(Pacote string) *PackageMaven {
	file, _ := os.ReadFile("tmp/" + Pacote + "/META-INF/MANIFEST.MF")
	str := string(file)

	str = strings.ReplaceAll(str, "\r", "")
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Name") {
			param := line[strings.Index(line, ":")+2:]
			if strings.Contains(param, ".") || strings.Contains(param, "/") {
				param = strings.ReplaceAll(param, "/", ".")
				paramslist := strings.Split(param, ".")
				if len(paramslist) > 3 {
					parameter := strings.Join(paramslist[:3], ".")
					return SolSearch(parameter)
				}

				return SolSearch(param)
			}
		}
	}
	return nil
}

func FilterString(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")
	return str
}

func checkoutMaven(PacoteName string) *PackageMaven {
	baseFolder := "tmp/" + PacoteName + "/META-INF/maven"
	var Pacote *PackageMaven = &PackageMaven{GroupId: "", ArtifactId: "", Version: ""}
	entries, _ := os.ReadDir(baseFolder)
	entries2, _ := os.ReadDir(baseFolder + "/" + entries[0].Name())
	file, _ := os.ReadFile(baseFolder + "/" + entries[0].Name() + "/" + entries2[0].Name() + "/pom.properties")

	str := string(file)

	lines := strings.Split(str, "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			argument := strings.Split(line, "=")
			switch argument[0] {
			case "version":
				Pacote.Version = argument[1]
				break
			case "groupId":
				Pacote.GroupId = argument[1]
				break
			case "artifactId":
				Pacote.ArtifactId = argument[1]
				break
			}
		}
	}

	if IsEmpty(*Pacote) {
		return nil
	}

	return Pacote
}

func checkDir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
			// File or directory does not exist
		} else {
			return false
			// Some other error such as missing permissions
		}
	}
	return true
}

func CheckPackage(Pacote string) *PackageMaven {
	err := ziper.UnzipSource("lib/"+Pacote+".jar", "tmp/"+Pacote)
	var pacoteExt *PackageMaven = nil

	if err != nil {
		fmt.Println("Error", err)
	}
	if checkDir("tmp/" + Pacote + "/META-INF/maven") {
		pacoteExt = checkoutMaven(Pacote)
	} else if checkDir("tmp/" + Pacote + "/META-INF/MANIFEST.MF") {
		pacoteExt = checkOutMainFestPack(Pacote)
	}

	return pacoteExt
}
