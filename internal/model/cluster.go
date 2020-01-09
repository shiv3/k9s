package model

import (
	"github.com/derailed/k9s/internal/client"
	v1 "k8s.io/api/core/v1"
	mv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type (
	// MetricsServer gather metrics information from pods and nodes.
	MetricsServer interface {
		MetricsService

		ClusterLoad(*v1.NodeList, *mv1beta1.NodeMetricsList, *client.ClusterMetrics) error
		NodesMetrics(*v1.NodeList, *mv1beta1.NodeMetricsList, client.NodesMetrics)
		PodsMetrics(*mv1beta1.PodMetricsList, client.PodsMetrics)
	}

	// MetricsService calls the metrics server for metrics info.
	MetricsService interface {
		HasMetrics() bool
		FetchNodesMetrics() (*mv1beta1.NodeMetricsList, error)
		FetchPodsMetrics(ns string) (*mv1beta1.PodMetricsList, error)
	}

	// Cluster represents a kubernetes resource.
	Cluster struct {
		client client.Connection
		mx     MetricsServer
	}
)

// NewCluster returns a new cluster info resource.
func NewCluster(c client.Connection, mx MetricsServer) *Cluster {
	return NewClusterWithArgs(c, mx)
}

// NewClusterWithArgs for tests only!
func NewClusterWithArgs(c client.Connection, mx MetricsServer) *Cluster {
	return &Cluster{client: c, mx: mx}
}

// Version returns the current K8s cluster version.
func (c *Cluster) Version() string {
	info, err := c.client.ServerVersion()
	if err != nil {
		return "n/a"
	}

	return info.GitVersion
}

// ContextName returns the context name.
func (c *Cluster) ContextName() string {
	n, err := c.client.Config().CurrentContextName()
	if err != nil {
		return "n/a"
	}
	return n
}

// ClusterName returns the cluster name.
func (c *Cluster) ClusterName() string {
	n, err := c.client.Config().CurrentClusterName()
	if err != nil {
		return "n/a"
	}
	return n
}

// UserName returns the user name.
func (c *Cluster) UserName() string {
	n, err := c.client.Config().CurrentUserName()
	if err != nil {
		return "n/a"
	}
	return n
}

// Metrics gathers node level metrics and compute utilization percentages.
func (c *Cluster) Metrics(nn *v1.NodeList, nmx *mv1beta1.NodeMetricsList, mx *client.ClusterMetrics) error {
	return c.mx.ClusterLoad(nn, nmx, mx)
}
