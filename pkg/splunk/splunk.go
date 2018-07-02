package splunk

//go:generate go-bindata -pkg=$GOPACKAGE -o=resources/bindata.go resources

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/plaid/go-envvar/envvar"
	splunkresources "github.com/skatteetaten/radish/pkg/splunk/resources"
	"github.com/skatteetaten/radish/pkg/util"
)

//Data : Struct for the required elements in the configuration json
type Data struct {
	SplunkIndex  string `envvar:"SPLUNK_INDEX" default:""`
	PodNamespace string `envvar:"POD_NAMESPACE" default:""`
	AppName      string `envvar:"APP_NAME" default:""`
	HostName     string `envvar:"HOST_NAME" default:""`
}

//GenerateStanzas :
func GenerateStanzas(templateFilePath string, splunkIndexFlag string,
	podNamespaceFlag string, appNameFlag string, hostNameFlag string, outputFilePath string) (bool, error) {
	vars := Data{}
	if err := envvar.Parse(&vars); err != nil {
		logrus.Fatal(err)
	}

	if splunkIndexFlag != "" {
		vars.SplunkIndex = splunkIndexFlag
	} else if vars.SplunkIndex == "" {
		return false, errors.New("No SplunkIndex present as flag or environment variable")
	}
	if podNamespaceFlag != "" {
		vars.PodNamespace = podNamespaceFlag
	} else if vars.PodNamespace == "" {
		return false, errors.New("No PodNamespace present as flag or environment variable")
	}
	if appNameFlag != "" {
		vars.AppName = appNameFlag
	} else if vars.AppName == "" {
		return false, errors.New("No AppName present as flag or environment variable")
	}
	if hostNameFlag != "" {
		vars.HostName = hostNameFlag
	} else if vars.HostName == "" {
		return false, errors.New("No HostName present as flag or environment variable")
	}

	stanzatemplate, err := readStanzasTemplate(templateFilePath)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to read template file from ", templateFilePath)
	}

	fileWriter := util.NewFileWriter(outputFilePath)

	if err := fileWriter(newSplunkStanzas(string(stanzatemplate), vars), outputFilePath); err != nil {
		return false, errors.Wrap(err, "Failed to write Splunk stanzas")
	}

	return true, nil
}

func readStanzasTemplate(templateFilePath string) ([]byte, error) {

	if len(templateFilePath) > 0 {
		stanzatemplate, err := ioutil.ReadFile(templateFilePath)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to read template file from ", templateFilePath)
		}
		return stanzatemplate, nil
	}

	stanzatemplate, err := splunkresources.Asset("resources/default_stanzas_template")
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read default template file from bindata")
	}
	return stanzatemplate, nil

}

func newSplunkStanzas(template string, data Data) util.WriterFunc {
	return util.NewTemplateWriter(
		data,
		"generatedSplunkStanzas",
		template)
}