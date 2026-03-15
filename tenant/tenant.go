package tenant

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jql/jdb"
)

const (
	MSG_TENANT_NOT_FOUND = "tenant not found"
)

var (
	ErrTenantNotFound = fmt.Errorf(MSG_TENANT_NOT_FOUND)
)

type Tenant struct {
	ID     string                `json:"id"`
	DB     *jdb.DB               `json:"db"`
	Models map[string]*jdb.Model `json:"models"`
}

/**
* newTenant
* @param id string, db *jdb.DB
* @return *Tenant
**/
func newTenant(id string, db *jdb.DB) *Tenant {
	return &Tenant{
		ID:     id,
		DB:     db,
		Models: make(map[string]*jdb.Model),
	}
}

/**
* ToJson
* @return et.Json
**/
func (s *Tenant) ToJson() (et.Json, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

var tenants map[string]*Tenant

func init() {
	tenants = make(map[string]*Tenant)
}

/**
* GetDb
* @param id string
* @return (*DB, error)
**/
func GetDb(id string) (*jdb.DB, error) {
	tenant, ok := tenants[id]
	if ok {
		return tenant.DB, nil
	}

	result, err := jdb.GetDb(id)
	if err != nil {
		return nil, jdb.ErrDbNotFound
	}

	tenants[id] = newTenant(id, result)
	return result, nil
}

/**
* GetModel
* @param tenantId, name string
* @return (*Model, error)
**/
func GetModel(tenantId, name string) (*jdb.Model, error) {
	tenant, ok := tenants[tenantId]
	if !ok {
		return nil, ErrTenantNotFound
	}

	result, ok := tenant.Models[name]
	if ok {
		return result, nil
	}

	result, err := tenant.DB.GetModel(name)
	if err != nil {
		return nil, err
	}

	tenant.Models[name] = result
	return result, nil
}

/**
* NewDb
* @param tenantId string, params et.Json
* @return *jdb.DB, error
**/
func NewDb(tenantId string, params et.Json) (*jdb.DB, error) {
	tenant, ok := tenants[tenantId]
	if ok {
		return tenant.DB, nil
	}

	result, err := jdb.Connect(tenantId, params)
	if err != nil {
		return nil, err
	}

	tenants[tenantId] = newTenant(tenantId, result)
	return result, nil
}
