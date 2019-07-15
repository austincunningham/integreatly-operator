package fuse

import (
	"context"
	"fmt"
	"github.com/integr8ly/integreatly-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/integreatly-operator/pkg/controller/installation/marketplace"
	"github.com/integr8ly/integreatly-operator/pkg/controller/installation/products/config"
	"github.com/integr8ly/integreatly-operator/pkg/resources"
	v1alpha12 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	syn "github.com/syndesisio/syndesis/install/operator/pkg/apis/syndesis/v1alpha1"
	v1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	pkgclient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	defaultInstallationNamespace = "fuse"
	defaultSubscriptionName      = "syndesis"
)

type Reconciler struct {
	client        pkgclient.Client
	coreClient    *kubernetes.Clientset
	Config        *config.Fuse
	ConfigManager config.ConfigReadWriter
	mpm           marketplace.MarketplaceInterface
	logger        *logrus.Entry
	nsReconciler  *resources.NamespaceReconciler
}

func NewReconciler(client pkgclient.Client, rc *rest.Config, coreClient *kubernetes.Clientset, configManager config.ConfigReadWriter, instance *v1alpha1.Installation, logger *logrus.Entry) (*Reconciler, error) {
	fuseConfig, err := configManager.ReadFuse()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve keycloak codeReadyConfig")
	}

	if fuseConfig.GetNamespace() == "" {
		fuseConfig.SetNamespace(instance.Spec.NamespacePrefix + defaultInstallationNamespace)
	}
	if err = fuseConfig.Validate(); err != nil {
		return nil, errors.Wrap(err, "keycloak codeReadyConfig is not valid")
	}

	mpm := marketplace.NewManager(client, rc)
	nsr := resources.NewNamespaceReconciler(client, logger)
	return &Reconciler{client: client,
		coreClient:    coreClient,
		ConfigManager: configManager,
		Config:        fuseConfig,
		mpm:           mpm,
		logger:        logger,
		nsReconciler:  nsr,
	}, nil
}

func (r *Reconciler) Reconcile(inst *v1alpha1.Installation) (v1alpha1.StatusPhase, error) {
	ctx := context.TODO()

	reconciledPhase, err := r.reconcileNamespace(ctx, inst)
	if err != nil {
		return v1alpha1.PhaseFailed, errors.Wrap(err, " failed to reconcile namespace for fuse ")
	}

	reconciledPhase, err = r.reconcileSubscription(ctx)
	if err != nil {
		return v1alpha1.PhaseFailed, errors.Wrap(err, " failed to reconcile subscription for fuse ")
	}

	reconciledPhase, err = r.handleAwaitingOperator(ctx)
	if err != nil {
		return v1alpha1.PhaseFailed, errors.Wrap(err, " failed to reconcile subscription for fuse ")
	}
	reconciledPhase, err = r.reconcileCustomResource(ctx, inst)
	if err != nil {
		return v1alpha1.PhaseFailed, errors.Wrap(err, " failed to reconcile fuse custom resource ")
	}

	logrus.Info("End of reconcile Phase : ", reconciledPhase)

	return v1alpha1.StatusPhase(v1alpha1.PhaseCompleted), nil
}

func (r *Reconciler) reconcileNamespace(ctx context.Context, inst *v1alpha1.Installation) (v1alpha1.StatusPhase, error) {
	if !inst.Spec.CreateNamespaces {
		r.logger.Infof("skipping creating namespace: %s", r.Config.GetNamespace())
		return v1alpha1.PhaseCreatingSubscription, nil
	}
	ns := &v1.Namespace{
		ObjectMeta: v12.ObjectMeta{
			Name: r.Config.GetNamespace(),
		},
	}
	ns, err := r.nsReconciler.Reconcile(ctx, ns, inst)
	if err != nil {
		return v1alpha1.PhaseFailed, errors.Wrap(err, "failed to reconcile fuse namespace "+r.Config.GetNamespace())
	}
	if ns.Status.Phase == v1.NamespaceTerminating {
		r.logger.Debugf("namespace %s is terminating, maintaining phase to try again on next reconcile", r.Config.GetNamespace())
		return v1alpha1.PhaseAwaitingNS, nil
	}
	if ns.Status.Phase != v1.NamespaceActive {
		return v1alpha1.PhaseAwaitingNS, nil
	}
	// all good return no status if already
	return v1alpha1.PhaseNone, nil
}

func (r *Reconciler) reconcileSubscription(ctx context.Context) (v1alpha1.StatusPhase, error) {
	// NEED to ensure a subscription is created if not exists
	// need to make sure there is only one operator source
	r.logger.Infof("reconciling subscription %s from channel %s in namespace: %s", defaultSubscriptionName, "integreatly", r.Config.GetNamespace())
	err := r.mpm.CreateSubscription(
		marketplace.GetOperatorSources().Integreatly,
		r.Config.GetNamespace(),
		defaultSubscriptionName,
		marketplace.IntegreatlyChannel,
		[]string{r.Config.GetNamespace()},
		v1alpha12.ApprovalAutomatic)
	if err != nil && !errors2.IsAlreadyExists(err) {
		return v1alpha1.PhaseFailed, errors.Wrap(err, fmt.Sprintf("could not create subscription in namespace: %s", r.Config.GetNamespace()))
	}
	return r.handleAwaitingOperator(ctx)
}

func (r *Reconciler) reconcileCustomResource(ctx context.Context, install *v1alpha1.Installation) (v1alpha1.StatusPhase, error) {
	ref := v12.NewControllerRef(install, v1alpha1.SchemaGroupVersionKind)
	intLimit := -1
	cr := syn.Syndesis{
		ObjectMeta: v12.ObjectMeta{
			OwnerReferences: []v12.OwnerReference{
				*ref,
			},
			Namespace: r.Config.GetNamespace(),
			Name:      "integreatly",
		},
		TypeMeta: v12.TypeMeta{},
		Spec: syn.SyndesisSpec{
			Integration: syn.IntegrationSpec{
				Limit: &intLimit,
			},
			Components: syn.ComponentsSpec{
				Server: syn.ServerConfiguration{
					Features: syn.ServerFeatures{
						ExposeVia3Scale: true,
					},
				},
			},
		},
	}
	if err := r.client.Get(ctx, pkgclient.ObjectKey{Name: cr.Name}, &cr); err != nil {
		if errors2.IsNotFound(err) {
			if err := r.client.Create(ctx, &cr); err != nil && !errors2.IsAlreadyExists(err) {
				return v1alpha1.PhaseFailed, errors.Wrap(err, "failed to create a syndesis cr when reconciling cutom resource")
			}
			return v1alpha1.PhaseNone, nil
		}
		return v1alpha1.PhaseFailed, errors.Wrap(err, "failed to get a syndesis cr when reconciling cutom resource")
	}
	if cr.Status.Phase != syn.SyndesisPhaseInstalled {
		return v1alpha1.PhaseInProgress, nil
	}
	if cr.Status.Phase == syn.SyndesisPhaseStartupFailed {
		return v1alpha1.PhaseFailed, errors.New("syndesis has failed to install " + string(cr.Status.Reason))
	}
	return v1alpha1.PhaseNone, nil
}

func (r *Reconciler) handleAwaitingOperator(ctx context.Context) (v1alpha1.StatusPhase, error) {
	r.logger.Infof("checking installplan is created for subscription %s in namespace: %s", defaultSubscriptionName, r.Config.GetNamespace())
	ip, sub, err := r.mpm.GetSubscriptionInstallPlan(defaultSubscriptionName, r.Config.GetNamespace())
	if err != nil {
		if errors2.IsNotFound(err) {
			if sub != nil {
				logrus.Infof("time since created %v", time.Now().Sub(sub.CreationTimestamp.Time).Seconds())
			}
			if sub != nil && time.Now().Sub(sub.CreationTimestamp.Time) > time.Minute*3 {
				// delete subscription so it is recreated
				logrus.Info("removing subscription as no install plan ready yet will recreate")
				if err := r.client.Delete(ctx, sub, func(options *pkgclient.DeleteOptions) {
					gp := int64(0)
					options.GracePeriodSeconds = &gp

				}); err != nil {
					// not going to fail here will retry
					r.logger.Error("failed to delete sub after install plan was not available for more than 20 seconds . Ignoring will retry ", err)
				}
			}
			r.logger.Debugf(fmt.Sprintf("installplan resource is not found in namespace: %s", r.Config.GetNamespace()))
			return v1alpha1.PhaseAwaitingOperator, nil
		}
		return v1alpha1.PhaseFailed, errors.Wrap(err, fmt.Sprintf("could not retrieve installplan in namespace: %s", r.Config.GetNamespace()))
	}

	r.logger.Infof("installplan phase is %s", ip.Status.Phase)
	if ip.Status.Phase != v1alpha12.InstallPlanPhaseComplete {
		r.logger.Infof("fuse online install plan is not complete yet")
		return v1alpha1.PhaseAwaitingOperator, nil
	}
	r.logger.Infof("fuse online install plan is complete. Installation ready ")
	return v1alpha1.PhaseNone, nil
}
