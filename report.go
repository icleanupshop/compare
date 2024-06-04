package main

import (
	"console/internal/hcl"
	"html/template"
	"os"
)

const (
	terraformHTMLTemplate = "terraform-report-html.tmpl"
)

type Report struct {
	ReportName      string
	ModelFile       string
	Environments    []Environment
	ReportLineItems []ReportLineItem
}

type Environment struct {
	Name             string
	PathToTFVarsFile string
	TFVars           *tfvars.Tfvars
}

type ReportLineItem struct {
	KeyName string
	Values  []ReportLineItemValue
}

type ReportLineItemValue struct {
	Environment string
	Value       string
}

func Compare(m *tfvars.Tfvars, r *Report) {
	//compare sample keys to target keys
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

func CreateReport(r *Report) {
	funcMap := template.FuncMap{
		"dec": func(i int) int { return i - 1 },
	}

	//var tmplFile = "terraform-report-html.tmpl"
	tmpl, err := template.New(terraformHTMLTemplate).Funcs(funcMap).ParseFiles(terraformHTMLTemplate)
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("terraform_config_report.html")
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
