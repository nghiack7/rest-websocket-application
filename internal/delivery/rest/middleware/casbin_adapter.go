package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// NewCasbinEnforcer creates a new Casbin enforcer with PostgreSQL adapter
func newCasbinEnforcer(cfg *viper.Viper, db *gorm.DB) (*casbin.Enforcer, error) {
	// Load the RBAC model from file
	modelPath := cfg.GetString("casbin.model_path")
	if modelPath == "" {
		modelPath = "config/rbac_model.conf"
	}

	// Load the model
	modelText, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	// Create a new PostgreSQL adapter
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// Create a new enforcer
	enforcer, err := casbin.NewEnforcer(modelText, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// Load the policy
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	// Enable auto-save for policy changes
	enforcer.EnableAutoSave(true)

	return enforcer, nil
}
