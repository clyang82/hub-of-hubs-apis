package policy

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/tools/cache"

	policyviewv1 "github.com/clyang82/hub-of-hubs-apis/pkg/server/apis/policyview/v1"
	policyv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
)

type REST struct {
	// lister can enumerate policy lists that enforce policy
	lister         cache.GenericLister
	tableConverter rest.TableConvertor
}

// NewREST returns a RESTStorage object that will work against ManagedCluster resources
func NewREST(
	lister cache.GenericLister,
) *REST {
	return &REST{
		lister: lister,
		//tableConverter: storage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(internalversion.AddHandlers)},
	}
}

// New returns a new policy
func (s *REST) New() runtime.Object {
	return &policyv1.Policy{}
}

func (s *REST) NamespaceScoped() bool {
	return true
}

// NewList returns a new policy list
func (*REST) NewList() runtime.Object {
	return &policyv1.PolicyList{}
}

var _ = rest.Lister(&REST{})

// List retrieves a list of policy that match label.
func (s *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	_, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.NewForbidden(policyviewv1.Resource(), "", fmt.Errorf("unable to list policy without a user on the context"))
	}

	runtimePolicyList, err := s.lister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	policyList := &policyv1.PolicyList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyv1.GroupVersion.String(),
			Kind:       policyv1.Kind,
		},
		Items: []policyv1.Policy{},
	}
	for _, runtimePolicy := range runtimePolicyList {
		policyList.Items = append(policyList.Items, *(runtimePolicy.(*policyv1.Policy)))
	}

	return policyList, nil
}

func (c *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return c.tableConverter.ConvertToTable(ctx, object, tableOptions)
}

var _ = rest.Watcher(&REST{})

func (s *REST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	if ctx == nil {
		return nil, fmt.Errorf("Context is nil")
	}
	_, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.NewForbidden(policyviewv1.Resource(), "", fmt.Errorf("unable to list policy without a user on the context"))
	}

	// includeAllExistingClusters := (options != nil) && options.ResourceVersion == "0"
	// // watcher := cache.NewCacheWatcher(user, s.clusterCache, includeAllExistingClusters)
	// // s.clusterCache.AddWatcher(watcher)
	// s.lister.

	// go watcher.Watch()
	return nil, nil
}

var _ = rest.Getter(&REST{})

// Get retrieves a policy by name
func (s *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	_, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.NewForbidden(policyviewv1.Resource(), "", fmt.Errorf("unable to get policy without a user on the context"))
	}

	policyList, err := s.lister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, policy := range policyList {
		// if name == policy['Name'] {
		return policy, nil
		// }
	}

	return nil, errors.NewForbidden(policyviewv1.Resource(), "", fmt.Errorf("the user cannot get the policy %v", name))
}
