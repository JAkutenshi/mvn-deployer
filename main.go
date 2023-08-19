package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
)

type JarInfo struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Path       string `json:"path"`
}

func loadJarsInfoList(jarsJsonFilePath string) *[]JarInfo {
	log.Default().Println("JARs to maven repo uploader starts...")
	jarsJsonFile, err := os.Open(jarsJsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Printf("Successfully load \"%s\"", jarsJsonFilePath)
	log.Default().Printf("Unmarshalling json file...")
	jarsInfoList := unmarshallJarsInfo(jarsJsonFile)
	defer jarsJsonFile.Close()
	log.Default().Printf("Unmarshalling json file done! The file is closed")
	return jarsInfoList
}

func unmarshallJarsInfo(jarsJsonFile *os.File) *[]JarInfo {
	jsonFileBytes, err := io.ReadAll(jarsJsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var jarsInfoList []JarInfo
	json.Unmarshal(jsonFileBytes, &jarsInfoList)
	return &jarsInfoList
}

func deployJar(jarInfo JarInfo, url, mavenServerID string) {
	mavenJarInfo := jarInfo.GroupID + ":" + jarInfo.ArtifactID + ":" + jarInfo.Version
	log.Default().Printf("Uploading artifact \"%s\" started...", mavenJarInfo)
	err := exec.Command(
		"mvn",
		"deploy:deploy-file",
		"-Durl="+url,
		"-DrepositoryId="+mavenServerID,
		"-Dfile="+jarInfo.Path,
		"-DgroupId="+jarInfo.GroupID,
		"-DartifactId="+jarInfo.ArtifactID,
		"-Dversion="+jarInfo.Version,
		"-Dpackaging=jar",
		"-DgeneratePom=true").Run()
	if err != nil {
		log.Default().Printf("Error while uploading, \"%s\" is skipped\n%s", mavenJarInfo, err)
	} else {
		log.Default().Printf("Uploading artifact \"%s\" done!", mavenJarInfo)
	}
}

func readArgFlags() (jarsJsonFilePath, url, mavenServerID string) {
	flag.StringVar(
		&jarsJsonFilePath,
		"file",
		"",
		"Required: JSON-file path with JARs' descriptions")
	var host string
	flag.StringVar(
		&host,
		"host",
		"",
		"Required: GitLab instance's hostname e.g. \"gitlab.com\"")
	var projectID string
	flag.StringVar(
		&projectID,
		"proj",
		"",
		"Required: GitLab's project ID: number or unique name")
	flag.StringVar(
		&mavenServerID,
		"serv",
		"",
		"Required: Desired maven server ID in .m2/settings.xml with gitlab-token provided")
	flag.Parse()
	checkFlagsNotEmpty(jarsJsonFilePath, host, projectID, mavenServerID)
	url = "https://" + host + "/api/v4/projects/" + projectID + "/packages/maven"
	return
}

func checkFlagsNotEmpty(jarsJsonFilePath, host, projectID, mavenServerID string) {
	if jarsJsonFilePath == "" {
		log.Fatal("There no JARs' json file provided, exit")
	}
	if host == "" {
		log.Fatal("There no GitLab's hostname provided, exit")
	}
	if projectID == "" {
		log.Fatal("There no GitLab repository ID provided, exit")
	}
	if mavenServerID == "" {
		log.Fatal("There no maven server ID provided, exit")
	}
}

func main() {
	jarsJsonFilePath, url, mavenServerID := readArgFlags()
	jarsInfoList := *loadJarsInfoList(jarsJsonFilePath)
	log.Default().Println("JARs to maven's repo uploading starts...")
	for i := 0; i < len(jarsInfoList); i++ {
		deployJar(jarsInfoList[i], url, mavenServerID)
	}
	log.Default().Println("JARs to maven's repo uploading done!")
}
