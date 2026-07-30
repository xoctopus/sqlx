package models

import "github.com/xoctopus/sqlx/pkg/builder"

// MetaCatalog registers SQL meta tables used by migrator document generation.
var MetaCatalog = builder.NewCatalog()
