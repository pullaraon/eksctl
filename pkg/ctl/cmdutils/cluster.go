package cmdutils

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// ApplyFilter applies nodegroup filters and returns a log function
func ApplyFilter(clusterConfig *api.ClusterConfig, ngFilter *filter.NodeGroupFilter) func() {
	var (
		filteredNodeGroups        []*api.NodeGroup
		filteredManagedNodeGroups []*api.ManagedNodeGroup
	)

	for _, ng := range clusterConfig.NodeGroups {
		if ngFilter.Match(ng.NameString()) {
			filteredNodeGroups = append(filteredNodeGroups, ng)
		}
	}

	for _, ng := range clusterConfig.ManagedNodeGroups {
		if ngFilter.Match(ng.NameString()) {
			filteredManagedNodeGroups = append(filteredManagedNodeGroups, ng)
		}
	}

	clusterConfig.NodeGroups, clusterConfig.ManagedNodeGroups = filteredNodeGroups, filteredManagedNodeGroups

	return func() {
		ngFilter.LogInfo(clusterConfig)
	}
}

// ToKubeNodeGroups combines managed and unmanaged nodegroups and returns a slice of eks.KubeNodeGroup containing
// both types of nodegroups
func ToKubeNodeGroups(clusterConfig *api.ClusterConfig) []eks.KubeNodeGroup {
	var kubeNodeGroups []eks.KubeNodeGroup
	for _, ng := range clusterConfig.NodeGroups {
		kubeNodeGroups = append(kubeNodeGroups, ng)
	}
	for _, ng := range clusterConfig.ManagedNodeGroups {
		kubeNodeGroups = append(kubeNodeGroups, ng)
	}
	return kubeNodeGroups
}

// ToNodePools combines managed and self-managed nodegroups and returns a slice of api.NodePool
func ToNodePools(clusterConfig *api.ClusterConfig) []api.NodePool {
	var nodePools []api.NodePool
	for _, ng := range clusterConfig.NodeGroups {
		nodePools = append(nodePools, ng)
	}
	for _, ng := range clusterConfig.ManagedNodeGroups {
		nodePools = append(nodePools, ng)
	}
	return nodePools
}
