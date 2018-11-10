package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/platform9/nodeadm/pkg/logrus"
	executil "github.com/platform9/nodeadm/utils/exec"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := &apis.JoinConfiguration{}
		configPath := cmd.Flag("cfg").Value.String()
		if len(configPath) != 0 {
			config, err = utils.JoinConfigurationFromFile(configPath)
			if err != nil {
				log.Fatalf("Failed to read configuration from file %q: %v", configPath, err)
			}
		}
		apis.SetJoinDefaults(config)
		if err := apis.SetJoinDynamicDefaults(config); err != nil {
			log.Fatalf("Failed to set dynamic defaults: %v", err)
		}

		if errors := apis.ValidateJoin(config); len(errors) > 0 {
			log.Error("Failed to validate configuration:")
			for i, err := range errors {
				log.Errorf("%v: %v", i, err)
			}
			os.Exit(1)
		}

		nodeConfig, err := yaml.Marshal(config.NodeConfiguration)
		if err != nil {
			log.Fatalf("\nFailed to marshal node config with err %v", err)
		}
		err = ioutil.WriteFile(constants.KubeadmConfig, nodeConfig, constants.Read)
		if err != nil {
			log.Fatalf("\nFailed to write file %q with error %v", constants.KubeadmConfig, err)
		}

		utils.InstallNodeComponents(config)
		kubeadmJoin()
	},
}

func kubeadmJoin() {
	cmd := exec.Command(filepath.Join(constants.BaseInstallDir, "kubeadm"), "join", "--ignore-preflight-errors=all", fmt.Sprintf("--config=%s", constants.KubeadmConfig))
	log.Infof("Running %q", strings.Join(cmd.Args, " "))
	if err := executil.LogRun(cmd); err != nil {
		log.Fatalf("%q failed: %s", strings.Join(cmd.Args, " "), err)
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().String("cfg", "", "Location of configuration file")
}
