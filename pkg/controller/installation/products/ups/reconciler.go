package ups

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	upsv1alpha1 "github.com/aerogear/unifiedpush-operator/pkg/apis/push/v1alpha1"

	"github.com/integr8ly/cloud-resource-operator/pkg/apis/integreatly/v1alpha1/types"
	croUtil "github.com/integr8ly/cloud-resource-operator/pkg/resources"

	integreatlyv1alpha1 "github.com/integr8ly/integreatly-operator/pkg/apis/integreatly/v1alpha1"
	monitoringv1alpha1 "github.com/integr8ly/integreatly-operator/pkg/apis/monitoring/v1alpha1"
	"github.com/integr8ly/integreatly-operator/pkg/controller/installation/marketplace"
	"github.com/integr8ly/integreatly-operator/pkg/controller/installation/products/config"
	"github.com/integr8ly/integreatly-operator/pkg/controller/installation/products/monitoring"
	"github.com/integr8ly/integreatly-operator/pkg/resources"

	routev1 "github.com/openshift/api/route/v1"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/lib/ownerutil"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	pkgclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	defaultInstallationNamespace = "ups"
	defaultUpsName               = "ups"
	defaultSubscriptionName      = "integreatly-unifiedpush"
	defaultRoutename             = defaultUpsName + "-unifiedpush-proxy"
	manifestPackage              = "integreatly-unifiedpush"
	tier                         = "production"
)

type Reconciler struct {
	Config        *config.Ups
	ConfigManager config.ConfigReadWriter
	mpm           marketplace.MarketplaceInterface
	logger        *logrus.Entry
	*resources.Reconciler
}

func NewReconciler(configManager config.ConfigReadWriter, installation *integreatlyv1alpha1.Installation, mpm marketplace.MarketplaceInterface) (*Reconciler, error) {
	upsConfig, err := configManager.ReadUps()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve ups config: %w", err)
	}

	if upsConfig.GetNamespace() == "" {
		upsConfig.SetNamespace(installation.Spec.NamespacePrefix + defaultInstallationNamespace)
	}

	upsConfig.SetBlackboxTargetPath("/rest/auth/config/")

	logger := logrus.NewEntry(logrus.StandardLogger())

	return &Reconciler{
		ConfigManager: configManager,
		Config:        upsConfig,
		mpm:           mpm,
		logger:        logger,
		Reconciler:    resources.NewReconciler(mpm),
	}, nil
}

func (r *Reconciler) GetPreflightObject(ns string) runtime.Object {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unifiedpush-operator",
			Namespace: ns,
		},
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, installation *integreatlyv1alpha1.Installation, product *integreatlyv1alpha1.InstallationProductStatus, serverClient pkgclient.Client) (integreatlyv1alpha1.StatusPhase, error) {
	logrus.Infof("Reconciling %s", defaultUpsName)

	phase, err := r.ReconcileFinalizer(ctx, serverClient, installation, string(r.Config.GetProductName()), func() (integreatlyv1alpha1.StatusPhase, error) {
		phase, err := resources.RemoveNamespace(ctx, installation, serverClient, r.Config.GetNamespace())
		if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
			return phase, err
		}
		return integreatlyv1alpha1.PhaseCompleted, nil
	})
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	ns := r.Config.GetNamespace()

	phase, err = r.ReconcileNamespace(ctx, ns, installation, serverClient)
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	namespace, err := resources.GetNS(ctx, ns, serverClient)
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, err
	}

	phase, err = r.ReconcileSubscription(ctx, namespace, marketplace.Target{Pkg: defaultSubscriptionName, Namespace: ns, Channel: marketplace.IntegreatlyChannel, ManifestPackage: manifestPackage}, ns, serverClient)
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	phase, err = r.reconcileComponents(ctx, installation, serverClient)
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	phase, err = r.reconcileHost(ctx, serverClient)
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	phase, err = r.reconcileBlackboxTargets(ctx, installation, serverClient)
	if err != nil || phase != integreatlyv1alpha1.PhaseCompleted {
		return phase, err
	}

	product.Host = r.Config.GetHost()
	product.Version = r.Config.GetProductVersion()
	product.OperatorVersion = r.Config.GetOperatorVersion()

	logrus.Infof("%s is successfully reconciled", defaultUpsName)

	return integreatlyv1alpha1.PhaseCompleted, nil
}

func (r *Reconciler) reconcileComponents(ctx context.Context, installation *integreatlyv1alpha1.Installation, client pkgclient.Client) (integreatlyv1alpha1.StatusPhase, error) {
	logrus.Info("Reconciling external postgres")
	ns := installation.Namespace

	// setup postgres custom resource
	// this will be used by the cloud resources operator to provision a postgres instance
	postgresName := fmt.Sprintf("ups-postgres-%s", installation.Name)
	postgres, err := croUtil.ReconcilePostgres(ctx, client, installation.Spec.Type, tier, postgresName, ns, postgresName, ns, func(cr metav1.Object) error {
		ownerutil.EnsureOwner(cr, installation)
		return nil
	})
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("failed to reconcile postgres request: %w", err)
	}

	// wait for the postgres instance to reconcile
	if postgres.Status.Phase != types.PhaseComplete {
		return integreatlyv1alpha1.PhaseAwaitingComponents, nil
	}

	// get the secret created by the cloud resources operator
	// containing postgres connection details
	connSec := &corev1.Secret{}
	err = client.Get(ctx, pkgclient.ObjectKey{Name: postgres.Status.SecretRef.Name, Namespace: postgres.Status.SecretRef.Namespace}, connSec)
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("failed to get postgres credential secret: %w", err)
	}

	// Reconcile ups custom resource
	logrus.Info("Reconciling unified push server cr")
	cr := &upsv1alpha1.UnifiedPushServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaultUpsName,
			Namespace: r.Config.GetNamespace(),
		},
		Spec: upsv1alpha1.UnifiedPushServerSpec{
			ExternalDB: true,
			Database: upsv1alpha1.UnifiedPushServerDatabase{
				Name:     string(connSec.Data["database"]),
				Password: string(connSec.Data["password"]),
				User:     string(connSec.Data["username"]),
				Host:     string(connSec.Data["host"]),
				Port:     intstr.FromString(string(connSec.Data["port"])),
			},
		},
	}

	err = client.Get(ctx, pkgclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, cr)
	if err != nil {
		// If the error is not an isNotFound error
		if !k8serr.IsNotFound(err) {
			return integreatlyv1alpha1.PhaseFailed, err
		}

		// Otherwise create the cr
		if err := client.Create(ctx, cr); err != nil {
			return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("failed to create unifiedpush server custom resource during reconcile: %w", err)
		}
	}

	// Wait till the ups cr status is complete
	if cr.Status.Phase != upsv1alpha1.PhaseReconciling {
		logrus.Info("Waiting for unified push server cr phase to complete")
		return integreatlyv1alpha1.PhaseInProgress, nil
	}

	logrus.Info("Successfully reconciled unified push server custom resource")

	return integreatlyv1alpha1.PhaseCompleted, nil
}

func (r *Reconciler) reconcileHost(ctx context.Context, serverClient pkgclient.Client) (integreatlyv1alpha1.StatusPhase, error) {
	// Setting host on config to exposed route
	logrus.Info("Setting unified push server config host")
	upsRoute := &routev1.Route{}
	err := serverClient.Get(ctx, pkgclient.ObjectKey{Name: defaultRoutename, Namespace: r.Config.GetNamespace()}, upsRoute)
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("failed to get route for unified push server: %w", err)
	}

	r.Config.SetHost("https://" + upsRoute.Spec.Host)
	err = r.ConfigManager.WriteConfig(r.Config)
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("could not update unified push server config: %w", err)
	}

	logrus.Info("Successfully set unified push server host")

	return integreatlyv1alpha1.PhaseCompleted, nil
}

func (r *Reconciler) reconcileBlackboxTargets(ctx context.Context, installation *integreatlyv1alpha1.Installation, client pkgclient.Client) (integreatlyv1alpha1.StatusPhase, error) {
	cfg, err := r.ConfigManager.ReadMonitoring()
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("error reading monitoring config: %w", err)
	}

	err = monitoring.CreateBlackboxTarget("integreatly-ups", monitoringv1alpha1.BlackboxtargetData{
		Url:     r.Config.GetHost() + "/" + r.Config.GetBlackboxTargetPath(),
		Service: "ups-ui",
	}, ctx, cfg, installation, client)
	if err != nil {
		return integreatlyv1alpha1.PhaseFailed, fmt.Errorf("error creating ups blackbox target: %w", err)
	}

	return integreatlyv1alpha1.PhaseCompleted, nil
}
