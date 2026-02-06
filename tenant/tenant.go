package tenant

import (
	"fmt"

	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
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
func GetDb(tenantId string) (*jql.DB, error) {
	tenant, ok := tenants[tenantId]
	if ok {
		return tenant.DB, nil
	}

	result, err := jql.GetDb(tenantId)
	if err != nil {
		return nil, jql.ErrDbNotFound
	}

	tenants[tenantId] = newTenant(result)
	return result, nil
}

/**
* GetModel
* @param tenantId, schema,name string
* @return (*Model, error)
**/
func GetModel(tenantId, schema, name string) (*jql.Model, bool) {
	tenant, ok := tenants[tenantId]
	if !ok {
		return nil, false
	}

	key := name
	key = strs.Append(schema, key, ".")
	result, ok := tenant.Models[key]
	if ok {
		return result, true
	}

	result, err := tenant.DB.GetModel(schema, name)
	if err != nil {
		return nil, false
	}

	tenant.Models[key] = result
	return result, true
}

/**
* NewDb
* @param tenantId string
* @return (*DB, error)
**/
func NewDb(tenantId, host string, port int) (*jql.DB, error) {
	tenant, ok := tenants[tenantId]
	if ok {
		return tenant.DB, nil
	}

	result, err := jql.NewDb(tenantId, host, port)
	if err != nil {
		return nil, err
	}

	tenants[tenantId] = newTenant(result)
	return result, nil
}
