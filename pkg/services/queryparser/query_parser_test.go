package queryparser

import (
	"testing"

	"github.com/onsi/gomega"
)

func Test_QueryParser(t *testing.T) {
	tests := []struct {
		name       string
		qry        string
		qryParser  QueryParser
		outQry     string
		outValues  []interface{}
		wantErr    bool
		errMessage string
	}{
		{
			name:      "Testing just `=` sign",
			qry:       "=",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Testing incomplete query",
			qry:       "name=",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Testing incomplete join",
			qry:       "name='test' and ",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Testing escaped quote",
			qry:       `name='test\'123'`,
			qryParser: NewQueryParser(),
			outQry:    "name = ?",
			outValues: []interface{}{"test'123"},
			wantErr:   false,
		},
		{
			name:      "Testing wrong unescaped quote",
			qry:       `name='test'123'`,
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Testing IN",
			qry:       "name IN ('value1', 'value2')",
			qryParser: NewQueryParser(),
			outQry:    "name IN( ? , ?)",
			outValues: []interface{}{"value1", "value2"},
			wantErr:   false,
		},
		{
			name:      "Testing IN with single value",
			qry:       "name IN ('value1')",
			qryParser: NewQueryParser(),
			outQry:    "name IN( ?)",
			outValues: []interface{}{"value1"},
			wantErr:   false,
		},
		{
			name:      "Testing invalid IN (no values)",
			qry:       "name IN ()",
			qryParser: NewQueryParser(),
			outQry:    "",
			outValues: nil,
			wantErr:   true,
		},
		{
			name:      "Testing invalid IN (ends with comma)",
			qry:       "name IN ('value1',)",
			qryParser: NewQueryParser(),
			outQry:    "",
			outValues: nil,
			wantErr:   true,
		},
		{
			name:      "Testing invalid IN (no closed brace)",
			qry:       "name IN ('value1'",
			qryParser: NewQueryParser(),
			outQry:    "",
			outValues: nil,
			wantErr:   true,
		},
		{
			name:      "Testing IN in complex query",
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) or owner in ('owner1', 'owner2', 'owner3')) or owner=c or name=e and region LIKE '%test%' and instance_type=standard",
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?) or owner in( ? , ? , ?)) or owner = ? or name = ? and region LIKE ? and instance_type = ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "owner1", "owner2", "owner3", "c", "e", "%test%", "standard"},
			wantErr:   false,
		},
		{
			name:      "Testing IN in with non quoted and quoted values",
			qry:       "owner in (owner1, 'owner2', owner3)",
			qryParser: NewQueryParser(),
			outQry:    "owner in( ? , ? , ?)",
			outValues: []interface{}{"owner1", "owner2", "owner3"},
			wantErr:   false,
		},
		{
			name:      "Testing IN in quoted value containing a comma",
			qry:       "owner in (owner1, 'owner2,', owner3)",
			qryParser: NewQueryParser(),
			outQry:    "owner in( ? , ? , ?)",
			outValues: []interface{}{"owner1", "owner2,", "owner3"},
			wantErr:   false,
		},
		{
			name:      "Testing negated IN in complex query",
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) or owner not in ('owner1', 'owner2', 'owner3')) or owner=c or name=e and region LIKE '%test%'",
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?) or owner not  in( ? , ? , ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "owner1", "owner2", "owner3", "c", "e", "%test%"},
			wantErr:   false,
		},
		{
			name:      "Complex query with braces",
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		},
		{
			name:      "Complex query with ilike",
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region ILIKE '%TEST%'",
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner <> ? or region = ?)) or owner = ? or name = ? and region ILIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%TEST%"},
			wantErr:   false,
		},
		{
			name:      "Complex query with braces and quoted values with escaped quote",
			qry:       `((cloud_provider = 'Value' and name = 'val\'ue1') and (owner = value2 or region='b' ) ) or owner=c or name=e and region LIKE '%test%'`,
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "val'ue1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		},
		{
			name:      "Complex query with braces and quoted values with spaces",
			qry:       `((cloud_provider = 'Value' and name = 'val ue1') and (owner = ' value2  ' or region='b' ) ) or owner=c or name=e and region LIKE '%test%'`,
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "val ue1", " value2  ", "b", "c", "e", "%test%"},
			wantErr:   false,
		},
		{
			name:      "Complex query with braces and empty quoted values",
			qry:       `((cloud_provider = 'Value' and name = '') and (owner = ' value2  ' or region='' ) ) or owner=c or name=e and region LIKE '%test%'`,
			qryParser: NewQueryParser(),
			outQry:    "((cloud_provider = ? and name = ?) and (owner = ? or region = ?)) or owner = ? or name = ? and region LIKE ?",
			outValues: []interface{}{"Value", "", " value2  ", "", "c", "e", "%test%"},
			wantErr:   false,
		},
		{
			name:      "JSONB query",
			qry:       `manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar'`,
			qryParser: NewQueryParser("manifest"),
			outQry:    "manifest -> 'data' -> 'manifest' -> 'metadata' -> 'labels' ->> 'foo' = ?",
			outValues: []interface{}{"bar"},
			wantErr:   false,
		},
		{
			name:       "Invalid JSONB query",
			qry:        `manifest->'data'->'manifest'->'metadata'->'labels'->'foo' = 'bar'`,
			qryParser:  NewQueryParser("manifest"),
			outQry:     "manifest -> 'data' -> 'manifest' -> 'metadata' -> 'labels' ->> 'foo' = ?",
			outValues:  nil,
			wantErr:    true,
			errMessage: "[59] error parsing the filter: unexpected token `=`",
		},
		{
			name: "Complex JSONB query",
			qry: `manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar' and ` +
				`( manifest->'data'->'manifest' ->> 'foo' in ('value1', 'value2') or ` +
				`manifest->'data'->'manifest'->>'labels' <> 'foo1')`,
			qryParser: NewQueryParser("manifest"),
			outQry: "manifest -> 'data' -> 'manifest' -> 'metadata' -> 'labels' ->> 'foo' = ? and " +
				"(manifest -> 'data' -> 'manifest' ->> 'foo' in( ? , ?) or " +
				"manifest -> 'data' -> 'manifest' ->> 'labels' <> ?)",
			outValues: []interface{}{"bar", "value1", "value2", "foo1"},
			wantErr:   false,
		},

		{
			name: "10 JOINS (maximum allowed)",
			qry: "name = value1 " +
				"and name = value2 " +
				"and name = value3 " +
				"or name = value4 " +
				"and name=value5 " +
				"and name = value6 " +
				"and name = value7 " +
				"and name = value8 " +
				"and name = value9 " +
				"and name = value10 " +
				"or name = value11",
			qryParser: NewQueryParser(),
			wantErr:   false,
		},
		{
			name: "11 JOINS (too many)",
			qry: "name = value1 " +
				"and name = value2 " +
				"and name = value3 " +
				"or name = value4 " +
				"and name=value5 " +
				"and name = value6 " +
				"and name = value7 " +
				"and name = value8 " +
				"and name = value9 " +
				"and name = value10 " +
				"or name = value11 " +
				"and name = value12",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Complex query with unbalanced braces",
			qry:       "((cloud_provider = Value and name = value1) and (owner = value2 or region=b  ) or owner=c or name=e and region LIKE '%test%'",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Bad column name",
			qry:       "badcolumn=test",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Bad column name in complex query",
			qry:       "((cloud_provider = Value and name = value1) and (owner = value2 or region=b  ) or badcolumn=c or name=e and region LIKE '%test%'",
			qryParser: NewQueryParser(),
			wantErr:   true,
		},
		{
			name:      "Parse with column prefix",
			qry:       "((cloud_provider = Value and name = value1) and (owner <> value2 or region=b ) ) or owner=c or name=e and region LIKE '%test%'",
			qryParser: NewQueryParserWithColumnPrefix("prefix"),
			outQry:    "((prefix.cloud_provider = ? and prefix.name = ?) and (prefix.owner <> ? or prefix.region = ?)) or prefix.owner = ? or prefix.name = ? and prefix.region LIKE ?",
			outValues: []interface{}{"Value", "value1", "value2", "b", "c", "e", "%test%"},
			wantErr:   false,
		},
	}

	for _, testcase := range tests {
		tt := testcase

		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			qry, err := tt.qryParser.Parse(tt.qry)

			if err != nil && !tt.wantErr {
				t.Errorf("QueryParser() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && !tt.wantErr {
				t.Logf("qry: %s", tt.qry)
				t.Logf("err: %v", err)
			}
			g.Expect(err != nil).To(gomega.Equal(tt.wantErr))

			if err == nil && tt.outQry != "" {
				if tt.outQry != "" {
					g.Expect(qry.Query).To(gomega.Equal(tt.outQry))
				}
				if tt.outValues != nil {
					g.Expect(qry.Values).To(gomega.Equal(tt.outValues))
				}
			}

			if err != nil && tt.wantErr && tt.errMessage != "" {
				g.Expect(err.Error()).To(gomega.Equal(tt.errMessage))
			}
		})
	}
}
