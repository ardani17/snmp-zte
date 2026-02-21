package service

import (
	"context"
	"fmt"

	"github.com/ardani/snmp-zte/internal/cache"
	"github.com/ardani/snmp-zte/internal/config"
	"github.com/ardani/snmp-zte/internal/driver"
	"github.com/ardani/snmp-zte/internal/driver/c320"
	"github.com/ardani/snmp-zte/internal/model"
	"github.com/redis/go-redis/v9"
)

// ONUService provides ONU operations
type ONUService struct {
	cfg      *config.Config
	cache    cache.Cache
	drivers  map[string]driver.Driver
}

// NewONUService creates a new ONU service
func NewONUService(cfg *config.Config, redisClient *redis.Client) *ONUService {
	var c cache.Cache
	if redisClient != nil {
		c = cache.NewRedisCache(redisClient, cache.DefaultTTL)
	} else {
		c = cache.NewNoOpCache()
	}

	s := &ONUService{
		cfg:     cfg,
		cache:   c,
		drivers: make(map[string]driver.Driver),
	}

	// Initialize drivers for configured OLTs
	for _, oltCfg := range cfg.OLTs {
		d := s.createDriver(oltCfg)
		if d != nil {
			s.drivers[oltCfg.ID] = d
		}
	}

	return s
}

// createDriver creates a driver for the given OLT config
func (s *ONUService) createDriver(cfg config.OLTConfig) driver.Driver {
	switch cfg.Model {
	case "C320", "c320":
		return c320.New(cfg.IPAddress, uint16(cfg.Port), cfg.Community)
	// C300 and C600 will be added later
	// case "C300", "c300":
	// 	return c300.New(cfg.IPAddress, uint16(cfg.Port), cfg.Community)
	// case "C600", "c600":
	// 	return c600.New(cfg.IPAddress, uint16(cfg.Port), cfg.Community)
	default:
		return nil
	}
}

// getDriver returns driver for the given OLT ID
func (s *ONUService) getDriver(oltID string) (driver.Driver, error) {
	d, ok := s.drivers[oltID]
	if !ok {
		return nil, fmt.Errorf("OLT not found or unsupported model: %s", oltID)
	}
	return d, nil
}

// GetONUList returns list of ONUs for a Board/PON
func (s *ONUService) GetONUList(ctx context.Context, oltID string, boardID, ponID int) ([]model.ONUInfo, error) {
	d, err := s.getDriver(oltID)
	if err != nil {
		return nil, err
	}

	// Validate
	if !d.ValidateBoardID(boardID) {
		return nil, fmt.Errorf("invalid board ID: %d", boardID)
	}
	if !d.ValidatePonID(ponID) {
		return nil, fmt.Errorf("invalid PON ID: %d", ponID)
	}

	// Try cache first
	cacheKey := cache.ONUListKey(oltID, boardID, ponID)
	var cached []model.ONUInfo
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Fetch from driver
	onuList, err := d.GetONUList(ctx, boardID, ponID)
	if err != nil {
		return nil, err
	}

	// Add OLT ID to each ONU
	for i := range onuList {
		onuList[i].OLTID = oltID
	}

	// Cache result
	s.cache.Set(ctx, cacheKey, onuList, cache.DefaultTTL)

	return onuList, nil
}

// GetONUDetail returns detailed information for a single ONU
func (s *ONUService) GetONUDetail(ctx context.Context, oltID string, boardID, ponID, onuID int) (*model.ONUDetail, error) {
	d, err := s.getDriver(oltID)
	if err != nil {
		return nil, err
	}

	// Validate
	if !d.ValidateBoardID(boardID) {
		return nil, fmt.Errorf("invalid board ID: %d", boardID)
	}
	if !d.ValidatePonID(ponID) {
		return nil, fmt.Errorf("invalid PON ID: %d", ponID)
	}
	if !d.ValidateOnuID(onuID) {
		return nil, fmt.Errorf("invalid ONU ID: %d", onuID)
	}

	// Try cache first
	cacheKey := cache.ONUDetailKey(oltID, boardID, ponID, onuID)
	var cached model.ONUDetail
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Fetch from driver
	detail, err := d.GetONUDetail(ctx, boardID, ponID, onuID)
	if err != nil {
		return nil, err
	}

	// Add OLT ID
	detail.OLTID = oltID

	// Cache result
	s.cache.Set(ctx, cacheKey, detail, cache.DefaultTTL)

	return detail, nil
}

// GetEmptySlots returns available ONU slots
func (s *ONUService) GetEmptySlots(ctx context.Context, oltID string, boardID, ponID int) ([]model.ONUSlot, error) {
	d, err := s.getDriver(oltID)
	if err != nil {
		return nil, err
	}

	// Validate
	if !d.ValidateBoardID(boardID) {
		return nil, fmt.Errorf("invalid board ID: %d", boardID)
	}
	if !d.ValidatePonID(ponID) {
		return nil, fmt.Errorf("invalid PON ID: %d", ponID)
	}

	// Try cache first
	cacheKey := cache.EmptySlotsKey(oltID, boardID, ponID)
	var cached []model.ONUSlot
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	// Fetch from driver
	slots, err := d.GetEmptySlots(ctx, boardID, ponID)
	if err != nil {
		return nil, err
	}

	// Cache result
	s.cache.Set(ctx, cacheKey, slots, cache.DefaultTTL)

	return slots, nil
}

// ClearCache clears cache for a Board/PON
func (s *ONUService) ClearCache(ctx context.Context, oltID string, boardID, ponID int) error {
	// Clear all related cache keys
	keys := []string{
		cache.ONUListKey(oltID, boardID, ponID),
		cache.EmptySlotsKey(oltID, boardID, ponID),
	}

	for _, key := range keys {
		s.cache.Delete(ctx, key)
	}

	// Also clear individual ONU caches (1-128)
	for i := 1; i <= 128; i++ {
		key := cache.ONUDetailKey(oltID, boardID, ponID, i)
		s.cache.Delete(ctx, key)
	}

	return nil
}
