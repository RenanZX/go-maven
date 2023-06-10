package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RenanZX/go-maven/libchecker"
	. "github.com/RenanZX/go-maven/types"
)

type Config struct {
	name                string
	modelversion        string
	groupId             string
	artifactId          string
	version             string
	includeLibsNotFound string
}

type AccList struct {
	packList   []*PackageMaven
	exceptList []string
}

func getPackageList(c Config) AccList {
	entries, err := os.ReadDir("./lib")
	var jarList []*PackageMaven
	var exceptList []string

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d jar packages were found, obtaining dependencies...\n", len(entries))

	for _, e := range entries {
		jarname := strings.ReplaceAll(e.Name(), ".jar", "")
		jarname = strings.ReplaceAll(jarname, " ", "")

		packageExtracted := libchecker.CheckPackage(jarname)
		if packageExtracted != nil {
			jarList = append(jarList, packageExtracted)
		} else {
			exceptList = append(exceptList, jarname)
		}
	}
	fmt.Printf("Obtained %d packages in remote repository and %d were not found!\n", len(jarList), len(exceptList))
	if strings.Contains(c.includeLibsNotFound, "false") {
		return AccList{
			packList:   jarList,
			exceptList: nil,
		}
	}
	return AccList{
		packList:   jarList,
		exceptList: exceptList,
	}
}

func readConfig() Config {
	file, err := os.ReadFile("./param.config")
	if err != nil {
		log.Fatal(err)
	}

	str := string(file)
	params := strings.Split(str, "\n")
	var conf Config
	for _, param := range params {
		inputs := strings.Split(param, "=")
		switch inputs[0] {
		case "name":
			conf.name = inputs[1]
		case "modelVersion":
			conf.modelversion = inputs[1]
		case "groupId":
			conf.groupId = inputs[1]
		case "artifactId":
			conf.artifactId = inputs[1]
		case "version":
			conf.version = inputs[1]
		case "includeLibsNotFound":
			conf.includeLibsNotFound = inputs[1]
		}
	}
	return conf
}

func writeln(m *os.File, str string) {
	m.WriteString(str + "\n")
}

func writePom(c Config, IncludePackages AccList) {
	file, err := os.Create("pom.xml")
	if err != nil {
		return
	}
	defer file.Close()

	writeln(file, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	writeln(file, "<project xmlns=\"http://maven.apache.org/POM/4.0.0\"")
	writeln(file, "  xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"")
	writeln(file, "  xsi:schemaLocation=\"http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd\">")
	writeln(file, "  <modelVersion>"+c.modelversion+"</modelVersion>")
	writeln(file, "  <groupId>"+c.groupId+"</groupId>")
	writeln(file, "  <artifactId>"+c.artifactId+"</artifactId>")
	writeln(file, "  <version>"+c.version+"</version>")
	writeln(file, "  <name>"+c.name+"</name>")
	writeln(file, "  <properties>")
	writeln(file, "    <local.lib>lib</local.lib>")
	writeln(file, "  </properties>\n")
	writeln(file, "  <repositories>")
	writeln(file, "  		<id>local-maven-repository</id>")
	writeln(file, "  		<url>file:///${project.basedir}/maven-repository</url>")
	writeln(file, "  </repositories>\n")
	writeln(file, "  <dependencies>")
	for _, dep := range IncludePackages.packList {
		if dep != nil {
			writeln(file, "   <dependency>")
			writeln(file, "     <groupId>"+libchecker.FilterString(dep.GroupId)+"</groupId>")
			writeln(file, "     <artifactId>"+libchecker.FilterString(dep.ArtifactId)+"</artifactId>")
			writeln(file, "     <version>"+libchecker.FilterString(dep.Version)+"</version>")
			writeln(file, "   </dependency>")
		}
	}
	for _, Pack := range IncludePackages.exceptList {
		writeln(file, "   <dependency>")
		writeln(file, "     <groupId>"+libchecker.FilterString(Pack)+"</groupId>")
		writeln(file, "     <artifactId>"+libchecker.FilterString(Pack)+"</artifactId>")
		writeln(file, "     <version>1.0</version>")
		writeln(file, "   </dependency>")
	}
	writeln(file, "  </dependencies>\n")
	writeln(file, "  <build>")
	writeln(file, "  		<plugins>")
	writeln(file, "  			<plugin>")
	writeln(file, "  				<groupId>org.apache.maven.plugins</groupId>")
	writeln(file, "  				<artifactId>maven-install-plugin</artifactId>")
	writeln(file, "  				<version>2.5.2</version>")
	if len(IncludePackages.exceptList) > 0 {
		writeln(file, "  				<executions>")
		writeln(file, "  					<execution>")
		for index, Pack := range IncludePackages.exceptList {
			writeln(file, "  					  <execution>")
			writeln(file, fmt.Sprintf("  					    <id>%d</id>", index))
			writeln(file, "  					    <phase>initialize</phase>")
			writeln(file, "  					    <goals>")
			writeln(file, "  					      <goal>install-file</goal>")
			writeln(file, "  					    </goals>")
			writeln(file, "  					    <configuration>")
			writeln(file, "  					      <file>${local.lib}/"+libchecker.FilterString(Pack)+".jar")
			writeln(file, "  					      <groupId>"+libchecker.FilterString(Pack)+"</groupId>")
			writeln(file, "  					      <artifactId>"+libchecker.FilterString(Pack)+"</artifactId>")
			writeln(file, "  					      <version>1.0</version>")
			writeln(file, "  					      <packaging>jar</packaging>")
			writeln(file, "  					      <generatePom>true</generatePom>")
			writeln(file, "  					    </configuration>")
			writeln(file, "  					  </execution>")
		}
		writeln(file, "  					</execution>")
		writeln(file, "  				</executions>")
	}
	writeln(file, "  			</plugin>")
	writeln(file, "  		</plugins>")
	writeln(file, "  </build>")
	writeln(file, "</project>")

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

func main() {
	fmt.Println("Loading configuration file...")
	conf := readConfig()
	fmt.Println("Analyzing lib of jar folder...")
	if checkDir("./lib") {
		packagelist := getPackageList(conf)
		fmt.Println("Generating Pom.xml...")
		writePom(conf, packagelist)
		os.RemoveAll("tmp")
		fmt.Println("Pom generated sucessfully!")
	} else {
		fmt.Println("lib folder not found")
	}
}
