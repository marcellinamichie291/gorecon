package srctleaks

import (
	"github.com/mr-pmillz/gorecon/localio"
	"github.com/spf13/cobra"
	"reflect"
)

type Options struct {
	Company     string
	GithubToken string
	Domain      interface{}
	Output      string
}

func ConfigureCommand(cmd *cobra.Command) error {
	cmd.PersistentFlags().StringP("company", "c", "", "company name that your testing")
	cmd.PersistentFlags().StringP("github-token", "", "", "github personal access token for github API interaction")
	cmd.PersistentFlags().StringP("domain", "d", "", "domain string or file containing domains ex. domains.txt")
	cmd.PersistentFlags().StringP("output", "o", "", "report output dir")
	return nil
}

func (opts *Options) LoadFromCommand(cmd *cobra.Command) error {
	company, err := localio.ConfigureFlagOpts(cmd, &localio.LoadFromCommandOpts{
		Flag:       "company",
		IsFilePath: false,
		Opts:       opts.Company,
	})
	if err != nil {
		return err
	}
	opts.Company = company.(string)

	githubToken, err := localio.ConfigureFlagOpts(cmd, &localio.LoadFromCommandOpts{
		Flag:       "github-token",
		IsFilePath: false,
		Opts:       opts.GithubToken,
	})
	if err != nil {
		return err
	}
	opts.GithubToken = githubToken.(string)

	domain, err := localio.ConfigureFlagOpts(cmd, &localio.LoadFromCommandOpts{
		Flag:       "domain",
		IsFilePath: true,
		Opts:       opts.Domain,
	})
	if err != nil {
		return err
	}
	rt := reflect.TypeOf(domain)
	switch rt.Kind() {
	case reflect.Slice:
		opts.Domain = domain.([]string)
	case reflect.String:
		opts.Domain = domain.(string)
	}

	output, err := localio.ConfigureFlagOpts(cmd, &localio.LoadFromCommandOpts{
		Flag:       "output",
		IsFilePath: true,
		Opts:       opts.Output,
	})
	if err != nil {
		return err
	}
	opts.Output = output.(string)

	return nil
}
