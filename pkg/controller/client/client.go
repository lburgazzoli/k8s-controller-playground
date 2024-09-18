package client

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	ctrl.Client
}

func (c *Client) Apply(ctx context.Context, obj ctrl.Object, opts ...ctrl.PatchOption) error {

	err := c.Client.Patch(ctx, obj, ctrl.Apply, opts...)
	if err != nil {
		return fmt.Errorf("unable to pactch object %s: %w", obj, err)
	}

	return nil
}

func (c *Client) ApplyStatus(ctx context.Context, obj ctrl.Object, opts ...ctrl.SubResourcePatchOption) error {
	err := c.Client.Status().Patch(ctx, obj, ctrl.Apply, opts...)
	if err != nil {
		return fmt.Errorf("unable to pactch object %s: %w", obj, err)
	}

	return nil
}
