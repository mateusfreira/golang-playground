package main


import (
  //"gopkg.in/mgo.v2"
	"log"
	//"net/http"
  "github.com/xanzy/go-gitlab"
  "encoding/json"
	//"flag"
	"fmt"
	"io/ioutil"
	"os"
//	"strconv"
	//"time"
)

type Config struct {
	Host    string `json:"host"`
	ApiPath string `json:"api_path"`
	Token   string `json:"token"`
}

func main() {
	log.SetFlags(log.Lshortfile)
  file, e := ioutil.ReadFile("../config-gitlab.json")
  if e != nil {
		fmt.Printf("Config file error: %v\n", e)
		os.Exit(1)
	}
	var config Config
	json.Unmarshal(file, &config)
	fmt.Printf("Results: %+v\n", config)
	git := gitlab.NewClient(nil, config.Token)
	git.SetBaseURL(config.Host+config.ApiPath)
	//opt := &gitlab.ListProjectsOptions{gitlab.ListOptions{1, 100}, true, "name", "asc", "query", true}
	for i := 0; i < 10; i++ {
		opt := &gitlab.ListProjectsOptions{gitlab.ListOptions{i, 100}, false, "name", "asc", "query", false}
		projects, _, err := git.Projects.ListProjects(opt)
	  if err != nil {
	    fmt.Println(err.Error())
	    return
	  }

	  for _, project := range projects {
	  	fmt.Printf("Project: %s %d \n", *project.Name)
	  	var optCommit = &gitlab.ListCommitsOptions { gitlab.ListOptions{0, 100}, "" }
	  	commits, _, err := git.Commits.ListCommits(*project.ID, optCommit)
		if err != nil {
			fmt.Printf("Error to try to get the commit! %+v \n", err)
			os.Exit(1)
		}	  	
	  	for _, commit := range commits {
	  		fmt.Printf("Commit: %v \n", *commit)
	  	}
	  }		
	}

}
