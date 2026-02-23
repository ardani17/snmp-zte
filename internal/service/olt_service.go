package service

import (
	"sync"

	"github.com/ardani/snmp-zte/internal/config"
	"github.com/ardani/snmp-zte/internal/model"
)

// OLTService menyediakan operasi pengelolaan OLT.
type OLTService struct {
	cfg    *config.Config
	mu     sync.RWMutex
}

// NewOLTService membuat instance OLT service baru.
func NewOLTService(cfg *config.Config) *OLTService {
	return &OLTService{
		cfg: cfg,
	}
}

// List mengembalikan semua OLT yang terkonfigurasi
func (s *OLTService) List() []model.OLT {
	s.mu.RLock()
	defer s.mu.RUnlock()

	olts := make([]model.OLT, len(s.cfg.OLTs))
	for i, cfg := range s.cfg.OLTs {
		olts[i] = model.OLT{
			ID:          cfg.ID,
			Name:        cfg.Name,
			Model:       cfg.Model,
			IPAddress:   cfg.IPAddress,
			Port:        cfg.Port,
			Community:   "***", // Don't expose community
			BoardCount:  cfg.BoardCount,
			PonPerBoard: cfg.PonPerBoard,
		}
	}
	return olts
}

// Get mengembalikan data OLT berdasarkan ID
func (s *OLTService) Get(id string) (*model.OLT, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, cfg := range s.cfg.OLTs {
		if cfg.ID == id {
			return &model.OLT{
				ID:          cfg.ID,
				Name:        cfg.Name,
				Model:       cfg.Model,
				IPAddress:   cfg.IPAddress,
				Port:        cfg.Port,
				Community:   "***",
				BoardCount:  cfg.BoardCount,
				PonPerBoard: cfg.PonPerBoard,
			}, nil
		}
	}

	return nil, ErrOLTNotFound
}

// Create menambah OLT baru dan menyimpannya ke konfigurasi.
func (s *OLTService) Create(olt model.OLT) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := config.OLTConfig{
		ID:          olt.ID,
		Name:        olt.Name,
		Model:       olt.Model,
		IPAddress:   olt.IPAddress,
		Port:        olt.Port,
		Community:   olt.Community,
		BoardCount:  olt.BoardCount,
		PonPerBoard: olt.PonPerBoard,
	}

	if err := s.cfg.AddOLT(cfg); err != nil {
		return err
	}

	// Simpan perubahan ke file config/olts.json
	return config.Save(s.cfg)
}

// Update memperbarui OLT yang sudah ada
func (s *OLTService) Update(id string, olt model.OLT) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := config.OLTConfig{
		ID:          id,
		Name:        olt.Name,
		Model:       olt.Model,
		IPAddress:   olt.IPAddress,
		Port:        olt.Port,
		Community:   olt.Community,
		BoardCount:  olt.BoardCount,
		PonPerBoard: olt.PonPerBoard,
	}

	if err := s.cfg.UpdateOLT(id, cfg); err != nil {
		return err
	}

	return config.Save(s.cfg)
}

// Delete menghapus OLT
func (s *OLTService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.cfg.DeleteOLT(id); err != nil {
		return err
	}

	return config.Save(s.cfg)
}

// ErrOLTNotFound dikembalikan saat OLT tidak ditemukan
var ErrOLTNotFound = &ServiceError{Message: "OLT not found"}

// ServiceError merepresentasikan error service
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
