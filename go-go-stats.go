package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	fmt.Printf("go-go-stats!\n")
	urlPtr := flag.String("url", "http://localhost:8153/go/api/admin/config.xml", "The url of the GoCD config file")
	userPtr := flag.String("user", "user", "The user for authentication to get the config file")
	passwordPtr := flag.String("pwd", "pwd", "The password for the user")

	data := GetConfigFile(*urlPtr, *userPtr, *passwordPtr)
	if data == nil {
		fmt.Println("Couldn't get config file")
		os.Exit(1)
	}

	configPtr := GetParsedXml(data)
	numberOfPipelines := GetNumberOfPipelines(*configPtr)

	fmt.Println("Number of pipelines:", numberOfPipelines)
	fmt.Println("")

	numberOfPipelinesByTemplate := GetNumberOfPipelinesByTemplate(*configPtr)
	for template, count := range numberOfPipelinesByTemplate {
		if template == "" {
			template = "No template"
		}

		fmt.Print(count)
		fmt.Println("\t" + template)
	}
}

func GetConfigFile(url string, user string, password string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

type Cruise struct {
	XMLName        xml.Name        `xml:"cruise"`
	PipelineGroups []PipelineGroup `xml:"pipelines"`
}

type PipelineGroup struct {
	Group     string     `xml:"group,attr"`
	Pipelines []Pipeline `xml:"pipeline"`
}

type Pipeline struct {
	Name     string `xml:"name,attr"`
	Template string `xml:"template,attr"`
}

func GetParsedXml(data []byte) *Cruise {
	c := new(Cruise)
	err := xml.Unmarshal(data, c)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	return c
}

func GetNumberOfPipelineGroups(c Cruise) int {
	return len(c.PipelineGroups)
}

func GetNumberOfPipelines(c Cruise) int {
	count := 0
	for _, ps := range c.PipelineGroups {
		count = count + len(ps.Pipelines)
	}

	return count
}

func GetNumberOfPipelinesByTemplate(c Cruise) map[string]int {
	var m map[string]int
	m = make(map[string]int)

	for _, pg := range c.PipelineGroups {
		for _, p := range pg.Pipelines {
			m[p.Template] = m[p.Template] + 1
		}
	}

	return m
}
