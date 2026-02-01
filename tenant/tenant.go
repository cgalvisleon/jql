package tenant

import (
	"fmt"

	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jql/jql"
)

const (
	MSG_TENANT_NOT_FOUND = "tenant not found"
)

var (
	ErrTenantNotFound = fmt.Errorf(MSG_TENANT_NOT_FOUND)
)

type Tenant struct {
	DB     *jql.DB               `json:"db"`
	Models map[string]*jql.Model `json:"models"`
}

/**
* newTenant
* @param db *jql.DB
* @return *Tenant
**/
func newTenant(db *jql.DB) *Tenant {
	return &Tenant{
		DB:     db,
		Models: make(map[string]*jql.Model),
	}
}

/**
* ToJson
* @return et.Json
**/
func (s *Tenant) ToJson() et.Json {
	result := et.Json{
		"database": s.DB.Name,
		"models":   s.Models,
	}

	return result
}

var tenants map[string]*Tenant

func init() {
	tenants = make(map[string]*Tenant)
}

/**
* save
* @param tenantId string
**/
func save(tenantId string) {
	tenant, ok := tenants[tenantId]
	if !ok {
		return
	}

	scr := tenant.ToJson()
	cache.ObjetSet("tenant", tenantId, scr)
}

/**
* Delete
* @param tenantId string
**/
func Delete(tenantId string) {
	cache.ObjetDelete("tenant", tenantId)
}

/**
* GetDb
* @param tenantId string
* @return (*DB, error)
**/
func GetDb(tenantId string) (*jql.DB, bool) {
	if _, ok := tenants[tenantId]; !ok {
		return nil, false
	}

	return tenants[tenantId].DB, true
}

/**
* GetModel
* @param tenantId string, name string
* @return (*Model, error)
**/
func GetModel(tenantId string, name string) (*jql.Model, bool) {
	tenant, ok := tenants[tenantId]
	if !ok {
		return nil, false
	}

	if _, ok := tenant.Models[name]; !ok {
		return nil, false
	}

	return tenant.Models[name], true
}

/**
* LoadDb
* @param tenantId string, name string, params et.Json
* @return (*DB, error)
**/
func LoadDb(tenantId, name string, params et.Json) (*jql.DB, error) {
	tenant, ok := tenants[tenantId]
	if ok {
		return tenant.DB, nil
	}

	params.Set("database", name)
	db, err := jql.ConnectTo(tenantId, params)
	if err != nil {
		return nil, err
	}

	tenants[tenantId] = newTenant(db)
	save(tenantId)
	return db, nil
}

/**
* LoadModel
* @param tenantId string, model *Model
* @return (*Model, error)
**/
func LoadModel(tenantId string, model *jql.Model) (*jql.Model, error) {
	tenant, ok := tenants[tenantId]
	if !ok {
		return nil, ErrTenantNotFound
	}

	tenant.Models[model.Name] = model

	save(tenantId)
	return model, nil
}
