package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vmindtech/vke-agent/models"
	"github.com/vmindtech/vke-agent/utils"
)

var config models.Config

var rootCmd = &cobra.Command{
	Use:   "vke-agent",
	Short: "Simple command line tool for setting up Kubernetes clusters.",
	Long: `vke-agent is a simple command line tool for setting up Kubernetes clusters.
With this tool, you can quickly provision both master and worker nodes.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := cmd.Flags().Parse(args); err != nil {
			logrus.Error("Parsing flags error:", err)
			return
		}

		if err := utils.UpdateSystem(); err != nil {
			logrus.Error("System update error:", err)
			return
		}

		if err := utils.CreateDirectory("/etc/rancher/rke2"); err != nil {
			logrus.Error("Indexing error:", err)
			return
		}
		if err := utils.RKE2Config(config.Initialize, config.ServerAddress, config.RKE2AgentType, config.RKE2Token, config.TLSSan, config.RKE2NodeLabel); err != nil {
			logrus.Error("Config creation error:", err)
			return
		}

		if err := utils.RKE2Install(config.Kubeversion, config.RKE2AgentType); err != nil {
			logrus.Error("RKE2 installation error:", err)
			return
		}

		if err := utils.RKE2ServiceEnable(config.RKE2AgentType); err != nil {
			logrus.Error("Service enabled error:", err)
			return
		}
		if err := utils.RKE2ServiceStart(config.RKE2AgentType); err != nil {
			logrus.Error("Service initialization error:", err)
			return
		}
		if config.Initialize {
			err := utils.PushRKE2Config(config.Initialize, config.RKE2AgentType, config.ServerAddress, config.RKE2ClusterName, config.RKE2ClusterUUID, config.RKE2AgentVKEAPIEndpoint, config.RKE2AgentVKEAPIAuthToken)
			if err != nil {
				logrus.Error("RKE2 config push error:", err)
				return
			}
			logrus.Info("RKE2 config pushed.")
		} else {
			logrus.Info("RKE2 config not pushed.")
		}
		logrus.Info("Process completed.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.ServerAddress, "serverAddress", "", "Server Address (required)")
	rootCmd.PersistentFlags().StringVar(&config.Kubeversion, "kubeversion", "", "Kube Version (required)")
	rootCmd.PersistentFlags().StringVar(&config.TLSSan, "tlsSan", "", "TLS San (required)")
	rootCmd.PersistentFlags().BoolVar(&config.Initialize, "initialize", false, "Initialize (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2Token, "rke2Token", "", "RKE2 Token (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2AgentType, "rke2AgentType", "", "Type (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2NodeLabel, "rke2NodeLabel", "", "Node Label (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2ClusterName, "rke2ClusterName", "", "Cluster Name (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2ClusterUUID, "rke2ClusterUUID", "", "Cluster UUID (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2AgentVKEAPIEndpoint, "rke2AgentVKEAPIEndpoint", "", "VKE API Endpoint (required)")
	rootCmd.PersistentFlags().StringVar(&config.RKE2AgentVKEAPIAuthToken, "rke2AgentVKEAPIAuthToken", "", "VKE API Auth Token (required)")

	rootCmd.SetHelpCommand(&cobra.Command{Use: "no-help-flag"})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
