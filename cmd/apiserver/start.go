package main

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"

	"github.com/munnerz/k8s-api-pager-demo/pkg/apis/pager/v1alpha1"
	"github.com/munnerz/k8s-api-pager-demo/pkg/apiserver"
	clientset "github.com/munnerz/k8s-api-pager-demo/pkg/client/clientset/internalversion"
	informers "github.com/munnerz/k8s-api-pager-demo/pkg/client/informers/internalversion"
)

const defaultEtcdPathPrefix = "/registry/pager.k8s.co"

type PagerServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	Admission          *genericoptions.AdmissionOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewPagerServerOptions(out, errOut io.Writer) *PagerServerOptions {
	o := &PagerServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, apiserver.Scheme, apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)),
		Admission:          genericoptions.NewAdmissionOptions(),

		StdOut: out,
		StdErr: errOut,
	}

	return o
}

// NewCommandStartPagerServer provides a CLI handler for starting the apiserver
func NewCommandStartPagerServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := NewPagerServerOptions(out, errOut)

	cmd := &cobra.Command{
		Short: "Launch a Pager API server",
		Long:  "Launch a Pager API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunPagerServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	o.Admission.AddFlags(flags)

	return cmd
}

func (o PagerServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.Admission.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *PagerServerOptions) Complete() error {
	return nil
}

func (o PagerServerOptions) Config() (*apiserver.Config, error) {
	// register admission plugins
	// banflunder.Register(o.Admission.Plugins)

	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(serverConfig.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}
	informerFactory := informers.NewSharedInformerFactory(client, serverConfig.LoopbackClientConfig.Timeout)

	config := &apiserver.Config{
		GenericConfig:         serverConfig,
		SharedInformerFactory: informerFactory,
	}
	return config, nil
}

func (o PagerServerOptions) RunPagerServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHook("start-pager-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.SharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
