package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	policyv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
)

const GroupName = "policyview.open-cluster-management.io"
const resource = "policies"

var (
	GroupVersion  = schema.GroupVersion{Group: GroupName, Version: "v1"}
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// Install is a function which adds this version to a scheme
	Install = schemeBuilder.AddToScheme

	// SchemeGroupVersion generated code relies on this name
	// Deprecated
	SchemeGroupVersion = GroupVersion
	// AddToScheme exists solely to keep the old generators creating valid code
	// DEPRECATED
	AddToScheme = schemeBuilder.AddToScheme
)

func Resource() schema.GroupResource {
	return schema.GroupResource{Group: policyv1.GroupVersion.Group, Resource: resource}
}

func GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    policyv1.GroupVersion.Group,
		Version:  policyv1.GroupVersion.Version,
		Resource: resource,
	}
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&policyv1.Policy{},
		&policyv1.PolicyList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
