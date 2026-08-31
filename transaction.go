package sanswitch

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
)

// ZoneTransaction coordinates a caller-controlled sequence of Zone mutations.
// It captures the database checksum at creation so Commit can detect changes
// made by other clients or processes.
type ZoneTransaction struct {
	client   *Client
	checksum string
	mutated  bool
	closed   bool
	mu       sync.Mutex
}

// BeginZoneTransaction validates write support and captures the current Zone
// database checksum. Transactions must be committed or aborted explicitly.
func (c *Client) BeginZoneTransaction(ctx context.Context) (*ZoneTransaction, error) {
	if err := c.ensureWriteSupported(); err != nil {
		return nil, err
	}
	checksum, err := c.ZoneChecksum(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin zone transaction: %w", err)
	}
	return &ZoneTransaction{client: c, checksum: checksum}, nil
}

// CreateZone creates a Zone inside tx.
func (tx *ZoneTransaction) CreateZone(ctx context.Context, name string, members, principalMembers []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name required")
	}
	if len(members) == 0 && len(principalMembers) == 0 {
		return errors.New("zone members required")
	}
	return tx.mutate(func() (bool, error) {
		return true, tx.client.CreateZone(ctx, name, members, principalMembers)
	})
}

// ReplaceZone replaces all members of an existing Zone inside tx.
func (tx *ZoneTransaction) ReplaceZone(ctx context.Context, name string, members, principalMembers []string) error {
	if len(members) == 0 {
		return errors.New("zone members required")
	}
	return tx.mutate(func() (bool, error) {
		return true, tx.client.UpdateZone(ctx, name, members, principalMembers)
	})
}

// DeleteZone deletes an existing Zone inside tx.
func (tx *ZoneTransaction) DeleteZone(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name required")
	}
	return tx.mutate(func() (bool, error) {
		return true, tx.client.DeleteZone(ctx, name)
	})
}

// AddZoneToConfig adds zoneName to a defined configuration if absent.
func (tx *ZoneTransaction) AddZoneToConfig(ctx context.Context, configName, zoneName string) error {
	if strings.TrimSpace(configName) == "" || strings.TrimSpace(zoneName) == "" {
		return errors.New("config and zone names required")
	}
	return tx.mutate(func() (bool, error) {
		configs, err := tx.client.DefinedConfigs(ctx)
		if err != nil {
			return false, err
		}
		members, err := memberZonesForConfig(configs, configName)
		if err != nil {
			return false, err
		}
		if slices.Contains(members, zoneName) {
			return false, nil
		}
		members = append(members, zoneName)
		return true, tx.client.UpdateDefinedConfig(ctx, configName, members)
	})
}

// RemoveZoneFromConfig removes zoneName from a defined configuration if
// present.
func (tx *ZoneTransaction) RemoveZoneFromConfig(ctx context.Context, configName, zoneName string) error {
	if strings.TrimSpace(configName) == "" || strings.TrimSpace(zoneName) == "" {
		return errors.New("config and zone names required")
	}
	return tx.mutate(func() (bool, error) {
		configs, err := tx.client.DefinedConfigs(ctx)
		if err != nil {
			return false, err
		}
		members, err := memberZonesForConfig(configs, configName)
		if err != nil {
			return false, err
		}
		members, removed := removeString(members, zoneName)
		if !removed {
			return false, nil
		}
		return true, tx.client.UpdateDefinedConfig(ctx, configName, members)
	})
}

// Commit saves all mutations and activates configName. A checksum conflict is
// returned by the switch and indicates that another writer won the race.
func (tx *ZoneTransaction) Commit(ctx context.Context, configName string) error {
	if tx == nil {
		return ErrTransactionClosed
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if !tx.mutated {
		return fmt.Errorf("%w: transaction has no mutations", ErrInvalidConfig)
	}
	if strings.TrimSpace(configName) == "" {
		return errors.New("config name required")
	}
	if err := tx.client.SaveZoneConfig(ctx, tx.checksum); err != nil {
		return tx.partial(fmt.Errorf("save zone configuration: %w", err))
	}
	checksum, err := tx.client.ZoneChecksum(ctx)
	if err != nil {
		return tx.partial(fmt.Errorf("refresh zone checksum: %w", err))
	}
	if err := tx.client.ActivateZoneConfig(ctx, configName, checksum); err != nil {
		return tx.partial(fmt.Errorf("activate zone configuration: %w", err))
	}
	tx.closed = true
	return nil
}

// Abort asks Fabric OS to discard the outstanding transaction.
func (tx *ZoneTransaction) Abort(ctx context.Context) error {
	if tx == nil {
		return ErrTransactionClosed
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if err := tx.client.AbortZoneTransaction(ctx); err != nil {
		return fmt.Errorf("abort zone transaction: %w", err)
	}
	tx.closed = true
	return nil
}

func (tx *ZoneTransaction) mutate(operation func() (bool, error)) error {
	if tx == nil {
		return ErrTransactionClosed
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	attempted, err := operation()
	tx.mutated = tx.mutated || attempted
	if err != nil {
		if attempted || tx.mutated {
			return tx.partial(err)
		}
		return err
	}
	return nil
}

func (tx *ZoneTransaction) ensureOpen() error {
	if tx == nil || tx.client == nil || tx.closed {
		return ErrTransactionClosed
	}
	return nil
}

func (tx *ZoneTransaction) partial(err error) error {
	return &PartialMutationError{Err: err}
}
