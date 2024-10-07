package run

import (
	"crypto/tls"
	"github.com/spf13/cobra"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	playgroundv1alpha1 "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	playgroundcontroller "github.com/lburgazzoli/k8s-controller-playground/internal/controller/playground"
	"github.com/lburgazzoli/k8s-controller-playground/pkg/controller/client"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(playgroundv1alpha1.AddToScheme(scheme))
}

const (
	cmdName = "run"
)

func NewCmd() *cobra.Command {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var tlsOpts []func(*tls.Config)

	// Outputs openAPI schema JSON containing the schema definitions in zz_generated.openapi.go.
	cmd := cobra.Command{
		Use:   cmdName,
		Short: cmdName,
		RunE: func(cmd *cobra.Command, args []string) error {
			disableHTTP2 := func(c *tls.Config) {
				setupLog.Info("disabling http/2")
				c.NextProtos = []string{"http/1.1"}
			}

			if !enableHTTP2 {
				tlsOpts = append(tlsOpts, disableHTTP2)
			}

			webhookServer := webhook.NewServer(webhook.Options{
				TLSOpts: tlsOpts,
			})

			metricsServerOptions := metricsserver.Options{
				BindAddress:   metricsAddr,
				SecureServing: secureMetrics,
				TLSOpts:       tlsOpts,
			}

			if secureMetrics {
				metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
			}

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:                        scheme,
				Metrics:                       metricsServerOptions,
				WebhookServer:                 webhookServer,
				HealthProbeBindAddress:        probeAddr,
				LeaderElection:                enableLeaderElection,
				LeaderElectionID:              "05b57063.lburgazzoli.github.io",
				LeaderElectionReleaseOnCancel: true,
			})

			if err != nil {
				setupLog.Error(err, "unable to cxreate manager")
				os.Exit(1)
			}

			c, err := client.New(cmd.Context(), mgr)
			if err != nil {
				setupLog.Error(err, "unable to create client")
				os.Exit(1)
			}

			if err = (&playgroundcontroller.ComponentReconciler{
				Client: c,
				Scheme: mgr.GetScheme(),
			}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "Component")
				os.Exit(1)
			}
			// +kubebuilder:scaffold:builder

			if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
				setupLog.Error(err, "unable to set up health check")
				os.Exit(1)
			}
			if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
				setupLog.Error(err, "unable to set up ready check")
				os.Exit(1)
			}

			setupLog.Info("starting manager")
			if err := mgr.Start(cmd.Context()); err != nil {
				setupLog.Error(err, "problem running manager")
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	cmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	cmd.Flags().BoolVar(&secureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	cmd.Flags().BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")

	return &cmd
}
