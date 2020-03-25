/**
 * Copyright 2019 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibmblockcsi

import (
	"fmt"

	csiv1 "github.com/IBM/ibm-block-csi-operator/pkg/apis/csi/v1"
	"github.com/IBM/ibm-block-csi-operator/pkg/config"
	csiversion "github.com/IBM/ibm-block-csi-operator/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// IBMBlockCSI is the wrapper for csiv1.IBMBlockCSI type
type IBMBlockCSI struct {
	*csiv1.IBMBlockCSI
	ServerVersion string
}

// New returns a wrapper for csiv1.IBMBlockCSI
func New(c *csiv1.IBMBlockCSI, serverVersion string) *IBMBlockCSI {
	return &IBMBlockCSI{
		IBMBlockCSI:   c,
		ServerVersion: serverVersion,
	}
}

// Unwrap returns the csiv1.IBMBlockCSI object
func (c *IBMBlockCSI) Unwrap() *csiv1.IBMBlockCSI {
	return c.IBMBlockCSI
}

// GetLabels returns all the labels to be set on all resources
func (c *IBMBlockCSI) GetLabels() labels.Set {
	labels := labels.Set{
		"app.kubernetes.io/name":       config.ProductName,
		"app.kubernetes.io/instance":   c.Name,
		"app.kubernetes.io/version":    csiversion.Version,
		"app.kubernetes.io/managed-by": config.Name,
		"csi":                          "ibm",
		"product":                      config.ProductName,
	}

	if c.Labels != nil {
		for k, v := range c.Labels {
			if !labels.Has(k) {
				labels[k] = v
			}
		}
	}

	return labels
}

// GetAnnotations returns all the annotations to be set on all resources
func (c *IBMBlockCSI) GetAnnotations() labels.Set {
	labels := labels.Set{
		"productID":      config.ProductName,
		"productName":    config.ProductName,
		"productVersion": csiversion.Version,
	}

	if c.Annotations != nil {
		for k, v := range c.Annotations {
			if !labels.Has(k) {
				labels[k] = v
			}
		}
	}

	return labels
}

// GetSelectorLabels returns labels used in label selectors
func (c *IBMBlockCSI) GetSelectorLabels(component string) labels.Set {
	return labels.Set{
		"app.kubernetes.io/component": component,
	}
}

func (c *IBMBlockCSI) GetCSIControllerSelectorLabels() labels.Set {
	return c.GetSelectorLabels(config.CSIController.String())
}

func (c *IBMBlockCSI) GetCSINodeSelectorLabels() labels.Set {
	return c.GetSelectorLabels(config.CSINode.String())
}

func (c *IBMBlockCSI) GetCSIControllerPodLabels() labels.Set {
	return labels.Merge(c.GetLabels(), c.GetCSIControllerSelectorLabels())
}

func (c *IBMBlockCSI) GetCSINodePodLabels() labels.Set {
	return labels.Merge(c.GetLabels(), c.GetCSINodeSelectorLabels())
}

func (c *IBMBlockCSI) GetCSIControllerImage() string {
	if c.Spec.Controller.Tag == "" {
		return c.Spec.Controller.Repository
	}
	return c.Spec.Controller.Repository + ":" + c.Spec.Controller.Tag
}

func (c *IBMBlockCSI) GetCSINodeImage() string {
	if c.Spec.Node.Tag == "" {
		return c.Spec.Node.Repository
	}
	return c.Spec.Node.Repository + ":" + c.Spec.Node.Tag
}

func (c *IBMBlockCSI) GetSidecarNames() []string {
	return []string{
		config.CSIAttacher,
		config.CSIProvisioner,
		config.CSINodeDriverRegistrar,
		config.LivenessProbe,
	}
}

func (c *IBMBlockCSI) GetSidecarByName(name string) *csiv1.CSISidecar {
	for _, sidecar := range c.Spec.Sidecars {
		if sidecar.Name == name {
			return &sidecar
		}
	}
	return nil
}

func (c *IBMBlockCSI) GetSidecarPullPolicyByName(name string) corev1.PullPolicy {
	sidecar := c.GetSidecarByName(name)
	if sidecar != nil && sidecar.ImagePullPolicy != "" {
		return sidecar.ImagePullPolicy
	}
	return corev1.PullIfNotPresent
}

func (c *IBMBlockCSI) GetSidecarImageByName(name string) string {
	sidecar := c.GetSidecarByName(name)
	if sidecar != nil {
		return fmt.Sprintf("%s:%s", sidecar.Repository, sidecar.Tag)
	}
	return ""
}

func (c *IBMBlockCSI) GetDefaultImageByPlatformAndName(platform, name string) string {
	switch platform {
	case config.OpenShift:
		return c.GetOpenShiftDefaultImageByName(name)
	case config.Kubernetes:
		return c.GetKubernetesDefaultImageByName(name)
	default:
		return ""
	}
}

func (c *IBMBlockCSI) GetKubernetesDefaultImageByName(name string) string {
	switch name {
	case config.ControllerImage:
		return config.ControllerRepository + ":" + config.ControllerTag
	case config.NodeImage:
		return config.NodeRepository + ":" + config.NodeTag
	case config.CSIProvisioner:
		return config.CSIProvisionerImage
	case config.CSIAttacher:
		return config.CSIAttacherImage
	case config.CSINodeDriverRegistrar:
		return config.NodeDriverRegistrarImage
	case config.LivenessProbe:
		return config.CSILivenessProbeImage
	default:
		return ""
	}
}

func (c *IBMBlockCSI) GetOpenShiftDefaultImageByName(name string) string {
	switch name {
	case config.ControllerImage:
		return config.OpenShiftControllerRepository + ":" + config.ControllerTag
	case config.NodeImage:
		return config.OpenShiftNodeRepository + ":" + config.NodeTag
	case config.CSIProvisioner:
		return config.OpenShiftCSIProvisionerImage
	case config.CSIAttacher:
		return config.OpenShiftCSIAttacherImage
	case config.CSINodeDriverRegistrar:
		return config.OpenShiftNodeDriverRegistrarImage
	case config.LivenessProbe:
		return config.OpenShiftCSILivenessProbeImage
	default:
		return ""
	}
}
