package logtohtml

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	_ "embed"
)

//go:embed log.tmpl
var tmpl string

func Write(in string, out string, initial string, chars []string) error {
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer file.Close()
	var logs strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs.WriteString(scanner.Text())
		logs.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file.Close()

	t, err := template.New("log").Parse(tmpl)
	if err != nil {
		return err
	}

	var tConfig struct {
		Active string
		Team   string
		Log    string
	}
	tConfig.Active = initial
	tConfig.Team = fmt.Sprint(strings.Join(strings.Split(fmt.Sprintf("%+q", chars), " "), ", "))
	tConfig.Log = logs.String()

	os.Remove(out)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	err = t.Execute(f, tConfig)
	if err != nil {
		return err
	}
	return nil
}

func WriteString(in string, out string, initial string, chars []string) error {
	t, err := template.New("log").Parse(tmpl)
	if err != nil {
		return err
	}

	var tConfig struct {
		Active string
		Team   string
		Log    string
	}
	tConfig.Active = initial
	var sb strings.Builder
	for _, v := range chars {
		sb.WriteString(v)
		sb.WriteString(",")
	}
	tConfig.Team = strings.Trim(sb.String(), ",")
	tConfig.Log = in

	os.Remove(out)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	err = t.Execute(f, tConfig)
	if err != nil {
		return err
	}
	return nil
}
