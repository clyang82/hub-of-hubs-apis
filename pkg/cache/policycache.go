// package cache

// import (
// 	"time"

// 	v1 "k8s.io/api/rbac/v1"
// 	"k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/labels"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/apimachinery/pkg/util/sets"
// 	utilwait "k8s.io/apimachinery/pkg/util/wait"
// 	"k8s.io/apiserver/pkg/authentication/user"
// 	rbacv1informers "k8s.io/client-go/informers/rbac/v1"
// 	clusterinformerv1 "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1"
// 	clusterv1lister "open-cluster-management.io/api/client/cluster/listers/cluster/v1"
// 	policyv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
// )

// // PolicyLister enforces ability to enumerate cluster based on role
// type PolicyLister interface {
// 	// List returns the list of PolicyList items that the user can access
// 	List(user user.Info, selector labels.Selector) (*policyv1.PolicyList, error)
// }

// type PolicyCache struct {
// 	policyLister clusterv1lister.ManagedPolicyLister
// }

// func NewPolicyCache(clusterInformer clusterinformerv1.ManagedClusterInformer,
// 	clusterRoleInformer rbacv1informers.ClusterRoleInformer,
// 	clusterRolebindingInformer rbacv1informers.ClusterRoleBindingInformer,
// 	getResourceNamesFromClusterRole func(*v1.ClusterRole, string, string) (sets.String, bool),
// ) *PolicyCache {
// 	dynInformer := dynamicinformer.NewDynamicSharedInformerFactory(client, 0)
// 	informer := dynInformer.ForResource(monboDBResource).Informer()

// 	clusterCache := &PolicyCache{
// 		policyLister: clusterInformer.Lister(),
// 	}

// 	return clusterCache
// }

// func (c *PolicyCache) ListResources() (sets.String, error) {
// 	allClusters := sets.String{}
// 	clusters, err := c.policyLister.List(labels.Everything())
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, cluster := range clusters {
// 		allClusters.Insert(cluster.Name)
// 	}
// 	return allClusters, nil
// }

// func (c *PolicyCache) List(userInfo user.Info, selector labels.Selector) (*clusterv1.ManagedClusterList, error) {
// 	names := c.cache.listNames(userInfo)

// 	clusterList := &clusterv1.ManagedClusterList{}
// 	for key := range names {
// 		cluster, err := c.policyLister.Get(key)
// 		if errors.IsNotFound(err) {
// 			continue
// 		}
// 		if err != nil {
// 			return nil, err
// 		}

// 		if !selector.Matches(labels.Set(cluster.Labels)) {
// 			continue
// 		}
// 		clusterList.Items = append(clusterList.Items, *cluster)
// 	}
// 	return clusterList, nil
// }

// func (c *PolicyCache) ListObjects(userInfo user.Info) (runtime.Object, error) {
// 	return c.List(userInfo, labels.Everything())
// }

// func (c *PolicyCache) Get(name string) (runtime.Object, error) {
// 	return c.policyLister.Get(name)
// }

// func (c *PolicyCache) ConvertResource(name string) runtime.Object {
// 	cluster, err := c.policyLister.Get(name)
// 	if err != nil {
// 		cluster = &clusterv1.ManagedCluster{ObjectMeta: metav1.ObjectMeta{Name: name}}
// 	}

// 	return cluster
// }

// func (c *PolicyCache) RemoveWatcher(w CacheWatcher) {
// 	c.cache.RemoveWatcher(w)
// }

// func (c *PolicyCache) AddWatcher(w CacheWatcher) {
// 	c.cache.AddWatcher(w)
// }

// // Run begins watching and synchronizing the cache
// func (c *PolicyCache) Run(period time.Duration) {
// 	go utilwait.Forever(func() { c.cache.synchronize() }, period)
// }
