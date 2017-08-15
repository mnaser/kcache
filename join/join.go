package join

import (
	"context"

	logutil "github.com/boz/go-logutil"
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/types/daemonset"
	"github.com/boz/kcache/types/deployment"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/replicaset"
	"github.com/boz/kcache/types/replicationcontroller"
	"github.com/boz/kcache/types/service"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func ServicePods(ctx context.Context, srcbase service.Controller, dstbase pod.Controller) (pod.Controller, error) {

	log := logutil.FromContextOrDefault(ctx)

	dst, err := dstbase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}

	update := func(_ *v1.Service) {
		objs, err := srcbase.Cache().List()
		if err == nil {
			log.Err(err, "join(service,pods): cache list")
			return
		}
		dst.Refilter(service.PodsFilter(objs...))
	}

	handler := service.BuildHandler().
		OnInitialize(func(objs []*v1.Service) { dst.Refilter(service.PodsFilter(objs...)) }).
		OnCreate(update).
		OnUpdate(update).
		OnDelete(update).
		Create()

	monitor, err := service.NewMonitor(srcbase, handler)
	if err != nil {
		dst.Close()
		return nil, log.Err(err, "join(service,pods): monitor")
	}

	go func() {
		<-dst.Done()
		monitor.Close()
	}()

	return dst, nil
}

func RCPods(ctx context.Context, srcbase replicationcontroller.Controller, dstbase pod.Controller) (pod.Controller, error) {

	log := logutil.FromContextOrDefault(ctx)

	dst, err := dstbase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}

	update := func(_ *v1.ReplicationController) {
		objs, err := srcbase.Cache().List()
		if err == nil {
			log.Err(err, "join(replicationcontroller,pods): cache list")
			return
		}
		dst.Refilter(replicationcontroller.PodsFilter(objs...))
	}

	handler := replicationcontroller.BuildHandler().
		OnInitialize(func(objs []*v1.ReplicationController) { dst.Refilter(replicationcontroller.PodsFilter(objs...)) }).
		OnCreate(update).
		OnUpdate(update).
		OnDelete(update).
		Create()

	monitor, err := replicationcontroller.NewMonitor(srcbase, handler)
	if err != nil {
		dst.Close()
		return nil, log.Err(err, "join(replicationcontroller,pods): monitor")
	}

	go func() {
		<-dst.Done()
		monitor.Close()
	}()

	return dst, nil
}

func RSPods(ctx context.Context, srcbase replicaset.Controller, dstbase pod.Controller) (pod.Controller, error) {

	log := logutil.FromContextOrDefault(ctx)

	dst, err := dstbase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}

	update := func(_ *v1beta1.ReplicaSet) {
		objs, err := srcbase.Cache().List()
		if err == nil {
			log.Err(err, "join(replicaset,pods): cache list")
			return
		}
		dst.Refilter(replicaset.PodsFilter(objs...))
	}

	handler := replicaset.BuildHandler().
		OnInitialize(func(objs []*v1beta1.ReplicaSet) { dst.Refilter(replicaset.PodsFilter(objs...)) }).
		OnCreate(update).
		OnUpdate(update).
		OnDelete(update).
		Create()

	monitor, err := replicaset.NewMonitor(srcbase, handler)
	if err != nil {
		dst.Close()
		return nil, log.Err(err, "join(replicaset,pods): monitor")
	}

	go func() {
		<-dst.Done()
		monitor.Close()
	}()

	return dst, nil
}

func DeploymentPods(ctx context.Context, srcbase deployment.Controller, dstbase pod.Controller) (pod.Controller, error) {

	log := logutil.FromContextOrDefault(ctx)

	dst, err := dstbase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}

	update := func(_ *v1beta1.Deployment) {
		objs, err := srcbase.Cache().List()
		if err == nil {
			log.Err(err, "join(deployment,pods): cache list")
			return
		}
		dst.Refilter(deployment.PodsFilter(objs...))
	}

	handler := deployment.BuildHandler().
		OnInitialize(func(objs []*v1beta1.Deployment) { dst.Refilter(deployment.PodsFilter(objs...)) }).
		OnCreate(update).
		OnUpdate(update).
		OnDelete(update).
		Create()

	monitor, err := deployment.NewMonitor(srcbase, handler)
	if err != nil {
		dst.Close()
		return nil, log.Err(err, "join(deployment,pods): monitor")
	}

	go func() {
		<-dst.Done()
		monitor.Close()
	}()

	return dst, nil
}

func DaemonSetPods(ctx context.Context, srcbase daemonset.Controller, dstbase pod.Controller) (pod.Controller, error) {

	log := logutil.FromContextOrDefault(ctx)

	dst, err := dstbase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}

	update := func(_ *v1beta1.DaemonSet) {
		objs, err := srcbase.Cache().List()
		if err == nil {
			log.Err(err, "join(daemonset,pods): cache list")
			return
		}
		dst.Refilter(daemonset.PodsFilter(objs...))
	}

	handler := daemonset.BuildHandler().
		OnInitialize(func(objs []*v1beta1.DaemonSet) { dst.Refilter(daemonset.PodsFilter(objs...)) }).
		OnCreate(update).
		OnUpdate(update).
		OnDelete(update).
		Create()

	monitor, err := daemonset.NewMonitor(srcbase, handler)
	if err != nil {
		dst.Close()
		return nil, log.Err(err, "join(daemonset,pods): monitor")
	}

	go func() {
		<-dst.Done()
		monitor.Close()
	}()

	return dst, nil
}