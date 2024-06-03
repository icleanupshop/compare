package main

import (
	"console/internal/hcl"
	"fmt"
	"html/template"
	"os"
)

type Report struct {
	ReportName      string
	IsModule        bool
	Environments    []Environment
	ReportLineItems []ReportLineItem
}

type Environment struct {
	Name   string
	TFVars *tfvars.Tfvars
}

type ReportLineItem struct {
	KeyName string
	Values  []ReportLineItemValue
}

type ReportLineItemValue struct {
	Environment string
	Value       string
}

func compare(m *tfvars.Tfvars, t *tfvars.Tfvars, r *Report) {
	//compare sample keys to target keys
	fmt.Println("Mode,sample,target")
	for _, v := range m.Keys() {
		rl := ReportLineItem{KeyName: v}
		//Get the model key from the target environments
		for _, targetEnv := range r.Environments {
			//Pull the model key from the target
			tv := targetEnv.TFVars.Get(v)

			//Add the target value to the report
			rl.Values = append(rl.Values, ReportLineItemValue{Environment: targetEnv.Name, Value: tv})
		}
		r.ReportLineItems = append(r.ReportLineItems, rl)
	}
}

func createReport(r *Report) {
	funcMap := template.FuncMap{
		"dec": func(i int) int { return i - 1 },
	}

	var tmplFile = "report-html.tmpl"
	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("pets.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, r)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	//Discover the environments as target for comparision with the model
	model, err := tfvars.New("/home/mj/projects/console/test/hcl/sample.tfvars")
	if err != nil {
		panic(err)
	}
	sample, err := tfvars.New("/home/mj/projects/console/test/hcl/sample.tfvars")
	if err != nil {
		panic(err)
	}

	target, err := tfvars.New("/home/mj/projects/console/test/hcl/target.tfvars")
	if err != nil {
		panic(err)
	}

	se := Environment{Name: "sample", TFVars: sample}
	te := Environment{Name: "target", TFVars: target}

	tfr := Report{}
	tfr.Environments = []Environment{se, te}
	tfr.ReportName = "GSS Terraform Config Report"

	compare(model, target, &tfr)

	createReport(&tfr)

}
