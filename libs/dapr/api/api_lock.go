package api

import (
	"context"

	"github.com/dapr/go-sdk/client"
	"github.com/hdget/sdk/common/namespace"
	"github.com/pkg/errors"
)

// Lock 锁
func (a daprApiImpl) Lock(ctx context.Context, lockStore, lockOwner, resource string, expiryInSeconds int) error {
	c, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	resp, err := c.TryLockAlpha1(ctx, namespace.Encapsulate(lockStore), &client.LockRequest{
		LockOwner:       lockOwner,
		ResourceID:      resource,
		ExpiryInSeconds: int32(expiryInSeconds),
	})
	if err != nil {
		return errors.Wrap(err, "try lock")
	}

	if !resp.Success {
		return errors.New("lock failed")
	}

	return nil
}

// Unlock 取消锁
func (a daprApiImpl) Unlock(ctx context.Context, lockStore, lockOwner, resource string) error {
	c, err := client.NewClient()
	if err != nil {
		return errors.Wrap(err, "new dapr client")
	}
	if c == nil {
		return errors.New("dapr client is null, name resolution service may not started, please check it")
	}

	resp, err := c.UnlockAlpha1(ctx, namespace.Encapsulate(lockStore), &client.UnlockRequest{
		LockOwner:  lockOwner,
		ResourceID: resource,
	})
	if err != nil {
		return errors.Wrap(err, "try lock")
	}

	if resp.StatusCode != 0 {
		return errors.New(resp.Status)
	}

	return nil
}
