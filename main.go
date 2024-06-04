package main

import (
	"console/internal/hcl"
	"flag"
	"github.com/go-git/go-git/v5"
	//"github.com/src-d/go-git"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	modelHelpMessage                 = "The full path to the TFVars file to use as the model."
	targetTFVarsDirectoryHelpMessage = "The full path to the directory where the TFVars for comparison reside."
	gitSourceConfigRepoHelpMessage   = "This is the HTTP path to the git repo holding the config for comparison."
	gitSourceConfigBranchHelpMessage = "This is the branch to use for the config comparison."
	gitSourceConfigRepoTmpDir        = "/tmp/console"
)

var (
	tfvarsDir             = "test/hcl"
	modelTFVarsFile       = "./test/hcl/sample.tfvars"
	gitSourceConfigRepo   = "https://github.com/icleanupshop/compare.git"
	gitSourceConfigBranch = "main"
)

func getTFVarsFiles(root string) []string {

	files, err := filepath.Glob(root + "**/**/*.tfvars")

	if err != nil {
		log.Panic(err)

	}
	return files
}

func getConfigRepo(r string) {
	// Clone the given repository to the given directory

	_, err := git.PlainClone(gitSourceConfigRepoTmpDir, false, &git.CloneOptions{
		URL:      r,
		Progress: os.Stdout,
	})
	if err != nil && err.Error() != "repository already exists" {
		log.Panic(err)

	}

}

func main() {
	//Get the date and time
	now := time.Now()
	log.Println("Current date and time:", now)

	//Discover the environments as target for comparision with the model

	flag.StringVar(&modelTFVarsFile, "model", modelTFVarsFile, modelHelpMessage)
	flag.StringVar(&tfvarsDir, "tfvarsDir", tfvarsDir, targetTFVarsDirectoryHelpMessage)
	flag.StringVar(&gitSourceConfigRepo, "sourceRepo", gitSourceConfigRepo, gitSourceConfigRepoHelpMessage)
	flag.StringVar(&gitSourceConfigBranch, "sourceRepoBranch", gitSourceConfigBranch, gitSourceConfigBranchHelpMessage)
	flag.Parse()

	//Get the module TFVars file
	model, err := tfvars.New(modelTFVarsFile)
	if err != nil {
		panic(err)
	}

	//Collect a list of target TFVars files
	files := getTFVarsFiles(gitSourceConfigRepoTmpDir + "/" + tfvarsDir)
	if len(files) == 0 {
		log.Println("Exiting, no TFVars files found in", gitSourceConfigRepoTmpDir+"/"+tfvarsDir)
		os.Exit(0)

	}

	defer os.Remove(gitSourceConfigRepoTmpDir)
	getConfigRepo(gitSourceConfigRepo)

	tfr := Report{}
	tfr.Environments = []Environment{}
	tfr.ReportName = "GSS Terraform Config Report - " + now.String()
	tfr.ModelFile = modelTFVarsFile

	/*files, e := OSReadDir(tfvarsDir)
	if e != nil {
		log.Panic(e)

	}*/

	//loop through list of TFVar files and extract meta data
	for _, f := range files {
		log.Println("Processing file ", f)

		p := strings.Split(f, "/")
		client := p[len(p)-2]
		e := p[len(p)-1]
		sample, err := tfvars.New(f)
		se := Environment{Name: strings.Split(e, ".")[0] + "-" + client, TFVars: sample, PathToTFVarsFile: f}
		tfr.Environments = append(tfr.Environments, se)

		if err != nil {
			panic(err)
		}

	}
	Compare(model, &tfr)

	CreateReport(&tfr)

	log.Println("Done.")

}
