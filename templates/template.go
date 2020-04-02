package templates

import (
	"fmt"
	"madnessBot/common/logger"
	"madnessBot/state/notify"
	"strings"
	"text/template"
)

const templatesPattern = "./templates/**/*.gotpl"

var tpl *template.Template

var templateFuncNamespaceMap = template.FuncMap{
	"NotifyGenerateString": notify.GenerateString,
}

func init() {
	var err error
	tpl, err = template.New("root").Funcs(templateFuncNamespaceMap).ParseGlob(templatesPattern)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to load templates")
	}
}

//ExecuteTemplate executes a template and returns the result
func ExecuteTemplate(name string, data interface{}) string {
	var buf strings.Builder

	err := tpl.ExecuteTemplate(&buf, fmt.Sprintf("%s.gotpl", name), data)

	if err != nil {
		logger.Log.Error().Err(err).Str("tpl_name", name).Msg("Failed to execute the template")
		return ""
	}
	return buf.String()
}
