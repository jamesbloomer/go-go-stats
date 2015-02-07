package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConfigfileSetsBasicAuth(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "CONFIG")
		u, p, o := (*r).BasicAuth()

		assert.Equal(u, "user")
		assert.Equal(p, "password")
		assert.True(o)
	}))

	defer server.Close()

	GetConfigFile(server.URL, "user", "password")
}

func TestGetConfigFileReturnsCorrectBody(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "CONFIG")
	}))

	defer server.Close()

	f := GetConfigFile(server.URL, "user", "password")
	assert.Equal("CONFIG\n", f)
}

func TestGetParsedXmlReturnsCorrectObjectsForCompleteXml(t *testing.T) {
	assert := assert.New(t)

	data := "<cruise><pipelines group=\"MyGroup\"><pipeline name=\"PName\" template=\"TName\"></pipeline></pipelines></cruise>"
	dom := GetParsedXml([]byte(data))
	assert.Equal(dom.PipelineGroups[0].Group, "MyGroup")
	assert.Equal(len(dom.PipelineGroups[0].Pipelines), 1)
	assert.Equal(dom.PipelineGroups[0].Pipelines[0].Name, "PName")
	assert.Equal(dom.PipelineGroups[0].Pipelines[0].Template, "TName")
}

func TestGetParsedXmlReturnsCorrectObjectsForIncompleteXml(t *testing.T) {
	assert := assert.New(t)

	data := "<cruise><server serverId=\"id\" /><pipelines group=\"MyGroup\"><pipeline name=\"PName\" template=\"TName\"></pipeline></pipelines></cruise>"
	dom := GetParsedXml([]byte(data))
	assert.Equal(len(dom.PipelineGroups), 1)
	assert.Equal(dom.PipelineGroups[0].Group, "MyGroup")
	assert.Equal(len(dom.PipelineGroups[0].Pipelines), 1)
	assert.Equal(dom.PipelineGroups[0].Pipelines[0].Name, "PName")
	assert.Equal(dom.PipelineGroups[0].Pipelines[0].Template, "TName")
}

func TestGetNumberOfPipelineGroupsForOnePipelineGroup(t *testing.T) {
	assert := assert.New(t)

	p := Pipeline{Name: "P1", Template: "T1"}
	pg := PipelineGroup{Pipelines: []Pipeline{p}}
	c := Cruise{PipelineGroups: []PipelineGroup{pg}}

	assert.Equal(1, GetNumberOfPipelineGroups(c))
}

func TestGetNumberOfPipelineGroupsForTwoPipelineGroups(t *testing.T) {
	assert := assert.New(t)

	p1 := Pipeline{Name: "P1", Template: "T1"}
	p2 := Pipeline{Name: "P2", Template: "T2"}
	pg1 := PipelineGroup{Pipelines: []Pipeline{p1}, Group: "G1"}
	pg2 := PipelineGroup{Pipelines: []Pipeline{p2}, Group: "G2"}
	c := Cruise{PipelineGroups: []PipelineGroup{pg1, pg2}}

	assert.Equal(2, GetNumberOfPipelineGroups(c))
}

func TestGetNumberOfPipelinesForOnePipeline(t *testing.T) {
	assert := assert.New(t)

	p := Pipeline{Name: "P1", Template: "T1"}
	pg := PipelineGroup{Pipelines: []Pipeline{p}}
	c := Cruise{PipelineGroups: []PipelineGroup{pg}}

	assert.Equal(1, GetNumberOfPipelines(c))
}

func TestGetNumberOfPipelinesForTwoPipelines(t *testing.T) {
	assert := assert.New(t)

	p1 := Pipeline{Name: "P1", Template: "T1"}
	p2 := Pipeline{Name: "P2", Template: "T2"}
	pg := PipelineGroup{Pipelines: []Pipeline{p1, p2}}
	c := Cruise{PipelineGroups: []PipelineGroup{pg}}

	assert.Equal(2, len(c.PipelineGroups[0].Pipelines))
	assert.Equal(2, GetNumberOfPipelines(c))
}

func TestGetNumberOfPipelinesByTemplate(t *testing.T) {
	assert := assert.New(t)

	p1 := Pipeline{Name: "P1", Template: "T1"}
	p2 := Pipeline{Name: "P2", Template: "T2"}
	p3 := Pipeline{Name: "P3", Template: "T2"}
	pg := PipelineGroup{Pipelines: []Pipeline{p1, p2, p3}}
	c := Cruise{PipelineGroups: []PipelineGroup{pg}}

	assert.Equal(1, GetNumberOfPipelinesByTemplate(c)["T1"])
	assert.Equal(2, GetNumberOfPipelinesByTemplate(c)["T2"])
}
