package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	policyv1 "open-cluster-management.io/governance-policy-propagator/api/v1"

	policyview "github.com/clyang82/hub-of-hubs-apis/pkg/server/apis/policyview"
	policyviewv1 "github.com/clyang82/hub-of-hubs-apis/pkg/server/apis/policyview/v1"
	"github.com/clyang82/hub-of-hubs-apis/pkg/server/rest/policy"
)

var (
	// Scheme contains the types needed by the resource metrics API.
	Scheme = runtime.NewScheme()
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
	// Codecs is a codec factory for serving the resource metrics API.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	// we need to add the options to empty v1
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	policyv1.AddToScheme(Scheme)
	policyview.Install(Scheme)
}

func Install(server *genericapiserver.GenericAPIServer,
	client dynamic.Interface, informerFactory dynamicinformer.DynamicSharedInformerFactory) error {
	if err := installPolicyViewGroup(server, client, informerFactory); err != nil {
		return err
	}
	return nil
}

func installPolicyViewGroup(server *genericapiserver.GenericAPIServer,
	client dynamic.Interface, informerFactory dynamicinformer.DynamicSharedInformerFactory) error {

	v1storage := map[string]rest.Storage{
		"policies": policy.NewREST(
			informerFactory.ForResource(policyviewv1.GroupVersionResource()).Lister(),
			client.Resource(policyviewv1.GroupVersionResource()),
		),
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(policyviewv1.GroupName, Scheme, ParameterCodec, Codecs)

	apiGroupInfo.VersionedResourcesStorageMap[policyviewv1.SchemeGroupVersion.Version] = v1storage

	return server.InstallAPIGroup(&apiGroupInfo)
}
