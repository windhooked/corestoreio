// Code generated by codegen. DO NOT EDIT.
// +build csall db

package storage

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

func TestNewTablesNonDB(t *testing.T) {
	ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de", MaxFloatDecimals: 6})
	_ = ps
}

func TestNewTablesDB(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewTables(ctx, ddl.WithConnPool(db))
	assert.NoError(t, err)
	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"core_configuration"}, tblNames)
	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de", MaxFloatDecimals: 6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
	)
	t.Run("CoreConfiguration_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCoreConfiguration)
		entSELECT := tbl.SelectByPK("*")
		// WithDBR generates the cached SQL string with empty key "".
		entSELECTStmtA := entSELECT.WithDBR().ExpandPlaceHolders()
		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
			dml.Column("id").LessOrEqual().Int(10),
		).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewCoreConfigurationCollection()
		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithDBR(ctx)
		for i := 0; i < 9; i++ {
			entIn := new(CoreConfiguration)
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}
			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfiguration_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()
			entOut := new(CoreConfiguration)
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.ID, entOut.ID, "IDX%d: ID should match", lID)
			assert.ExactlyLength(t, 8, &entIn.Scope, &entOut.Scope, "IDX%d: Scope should match", lID)
			assert.Exactly(t, entIn.ScopeID, entOut.ScopeID, "IDX%d: ScopeID should match", lID)
			assert.Exactly(t, entIn.Expires, entOut.Expires, "IDX%d: Expires should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Path, &entOut.Path, "IDX%d: Path should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.Value, &entOut.Value, "IDX%d: Value should match", lID)
			// ignoring: version_ts
			// ignoring: version_te
		}
		dmltest.Close(t, entINSERTStmtA)
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)
		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithDBR(ctx)
		lID := dmltest.CheckLastInsertID(t, "Error:  CoreConfigurationCollection ")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
}
