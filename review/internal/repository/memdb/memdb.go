// Memdb provides an in-memory implementation of the Repository interface.
package memdb

import (
	"coupon_service/internal/service/entity"
	"fmt"
	"sync"
)

type Repository struct {
	mu      sync.RWMutex
	entries map[string]entity.Coupon
}

// New creates a new in-memory repository.
func New() *Repository {
	return &Repository{entries: make(map[string]entity.Coupon, 0)}
}

// FindByCode finds a coupon by its code.
func (r *Repository) FindByCode(code string) (entity.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	coupon, ok := r.entries[code]
	if !ok {
		return entity.Coupon{}, fmt.Errorf("coupon not found")
	}
	return coupon, nil
}

// Save saves a coupon in the repository.
func (r *Repository) Save(coupon entity.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries[coupon.Code] = coupon
	return nil
}
