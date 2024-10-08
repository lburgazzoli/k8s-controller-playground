package client

import (
	"context"
	"fmt"
	playgroundClient "github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset/versioned"
	ctrlRt "sigs.k8s.io/controller-runtime"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
)

func New(_ context.Context, manager ctrlRt.Manager) (*Client, error) {
	playgroundCl, err := playgroundClient.NewForConfig(manager.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("unable to construct a playground client: %w", err)
	}

	k8sCli, err := kubernetes.NewForConfig(manager.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("unable to construct a playground client: %w", err)
	}

	return &Client{
		Client: manager.GetClient(),
		P:      playgroundCl,
		K:      k8sCli,
	}, nil
}

type Client struct {
	ctrlCli.Client
	P playgroundClient.Interface
	K kubernetes.Interface
}

func (c *Client) Apply(ctx context.Context, obj ctrlCli.Object, opts ...ctrlCli.PatchOption) error {

	err := c.Client.Patch(ctx, obj, ctrlCli.Apply, opts...)
	if err != nil {
		return fmt.Errorf("unable to pactch object %s: %w", obj, err)
	}

	return nil
}

func (c *Client) ApplyStatus(ctx context.Context, obj ctrlCli.Object, opts ...ctrlCli.SubResourcePatchOption) error {
	err := c.Client.Status().Patch(ctx, obj, ctrlCli.Apply, opts...)
	if err != nil {
		return fmt.Errorf("unable to pactch object %s: %w", obj, err)
	}

	return nil
}
