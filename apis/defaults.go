package apis

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
	corev1 "k8s.io/api/core/v1"
	kubeadmv1alpha2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha2"
)

// SetInitDefaults sets defaults on the configuration used by init
func SetInitDefaults(config *InitConfiguration) {
	kubeadmv1alpha2.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
	config.MasterConfiguration.Kind = "MasterConfiguration"
	config.MasterConfiguration.APIVersion = "kubeadm.k8s.io/v1alpha2"
	config.MasterConfiguration.KubernetesVersion = constants.KubernetesVersion
	config.MasterConfiguration.NodeRegistration.Taints = []corev1.Taint{} // empty slice denotes no taints
	addOrAppend(&config.MasterConfiguration.APIServerExtraArgs, "feature-gates", constants.FeatureGates)
	addOrAppend(&config.MasterConfiguration.ControllerManagerExtraArgs, "feature-gates", constants.FeatureGates)
	addOrAppend(&config.MasterConfiguration.SchedulerExtraArgs, "feature-gates", constants.FeatureGates)
}

// SetInitDynamicDefaults sets defaults derived at runtime
func SetInitDynamicDefaults(config *InitConfiguration) error {
	nodeName, err := constants.GetHostnameOverride()
	if err != nil {
		return fmt.Errorf("unable to dervice hostname override: %v", err)
	}
	config.MasterConfiguration.NodeRegistration.Name = nodeName
	return nil
}

// SetJoinDefaults sets defaults on the configuration used by join
func SetJoinDefaults(config *JoinConfiguration) {
}

// SetJoinDynamicDefaults sets defaults derived at runtime
func SetJoinDynamicDefaults(config *JoinConfiguration) error {
	nodeName, err := constants.GetHostnameOverride()
	if err != nil {
		return fmt.Errorf("unable to dervice hostname override: %v", err)
	}
	config.NodeConfiguration.NodeRegistration.Name = nodeName
	return nil
}

func addOrAppend(extraArgs *map[string]string, key string, value string) {
	// Create a new map if it doesn't exist.
	if *extraArgs == nil {
		*extraArgs = make(map[string]string)
	}
	// Add the key with the value if it doesn't exist. Otherwise, append the value
	// to the pre-existing values.
	prevFeatureGates := (*extraArgs)[key]
	if prevFeatureGates == "" {
		(*extraArgs)[key] = value
	} else {
		featureGates := prevFeatureGates + "," + value
		(*extraArgs)[key] = featureGates
	}
}
