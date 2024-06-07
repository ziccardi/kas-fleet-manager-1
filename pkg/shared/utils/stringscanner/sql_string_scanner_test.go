package stringscanner

import (
	"testing"

	"github.com/onsi/gomega"
)

func Test_SQLStringScanner(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectedTokens []Token
	}{{
		name:  "Test simple select",
		value: "SELECT * FROM ADDRESS_BOOK",
		expectedTokens: []Token{
			{TokenType: LITERAL, Value: "SELECT", Position: 0},
			{TokenType: LITERAL, Value: "*", Position: 7},
			{TokenType: LITERAL, Value: "FROM", Position: 9},
			{TokenType: LITERAL, Value: "ADDRESS_BOOK", Position: 14},
		},
	},
		{
			name:  "Test quoted string",
			value: "SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with spaces'",
			expectedTokens: []Token{
				{TokenType: LITERAL, Value: "SELECT", Position: 0},
				{TokenType: LITERAL, Value: "*", Position: 7},
				{TokenType: LITERAL, Value: "FROM", Position: 9},
				{TokenType: LITERAL, Value: "ADDRESS_BOOK", Position: 14},
				{TokenType: LITERAL, Value: "WHERE", Position: 27},
				{TokenType: LITERAL, Value: "SURNAME", Position: 33},
				{TokenType: OP, Value: "=", Position: 41},
				{TokenType: LITERAL, Value: "'surname with spaces'", Position: 43},
			},
		},
		{
			name:  "Test quoted string with comma",
			value: "SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with , comma'",
			expectedTokens: []Token{
				{TokenType: LITERAL, Value: "SELECT", Position: 0},
				{TokenType: LITERAL, Value: "*", Position: 7},
				{TokenType: LITERAL, Value: "FROM", Position: 9},
				{TokenType: LITERAL, Value: "ADDRESS_BOOK", Position: 14},
				{TokenType: LITERAL, Value: "WHERE", Position: 27},
				{TokenType: LITERAL, Value: "SURNAME", Position: 33},
				{TokenType: OP, Value: "=", Position: 41},
				{TokenType: LITERAL, Value: "'surname with , comma'", Position: 43},
			},
		},
		{
			name:  "Test escaped chars",
			value: `SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'surname with spaces and \'quote\''`,
			expectedTokens: []Token{
				{TokenType: LITERAL, Value: "SELECT", Position: 0},
				{TokenType: LITERAL, Value: "*", Position: 7},
				{TokenType: LITERAL, Value: "FROM", Position: 9},
				{TokenType: LITERAL, Value: "ADDRESS_BOOK", Position: 14},
				{TokenType: LITERAL, Value: "WHERE", Position: 27},
				{TokenType: LITERAL, Value: "SURNAME", Position: 33},
				{TokenType: OP, Value: "=", Position: 41},
				{TokenType: LITERAL, Value: `'surname with spaces and \'quote\''`, Position: 43},
			},
		},
		{
			name:  "Test SQL with operators",
			value: `SELECT * FROM ADDRESS_BOOK WHERE SURNAME = 'Mouse' AND AGE > 3`,
			expectedTokens: []Token{
				{TokenType: LITERAL, Value: "SELECT", Position: 0},
				{TokenType: LITERAL, Value: "*", Position: 7},
				{TokenType: LITERAL, Value: "FROM", Position: 9},
				{TokenType: LITERAL, Value: "ADDRESS_BOOK", Position: 14},
				{TokenType: LITERAL, Value: "WHERE", Position: 27},
				{TokenType: LITERAL, Value: "SURNAME", Position: 33},
				{TokenType: OP, Value: "=", Position: 41},
				{TokenType: LITERAL, Value: "'Mouse'", Position: 43},
				{TokenType: LITERAL, Value: "AND", Position: 51},
				{TokenType: LITERAL, Value: "AGE", Position: 55},
				{TokenType: OP, Value: ">", Position: 59},
				{TokenType: LITERAL, Value: "3", Position: 61},
			},
		},
		{
			name:  "Test SQL with JSONB",
			value: `select * from table where manifest->'data'->'manifest'->'metadata'->'labels'->>'foo' = 'bar'`,
			expectedTokens: []Token{
				{TokenType: LITERAL, Value: "select", Position: 0},
				{TokenType: LITERAL, Value: "*", Position: 7},
				{TokenType: LITERAL, Value: "from", Position: 9},
				{TokenType: LITERAL, Value: "table", Position: 14},
				{TokenType: LITERAL, Value: "where", Position: 20},
				{TokenType: LITERAL, Value: "manifest", Position: 26},
				{TokenType: OP, Value: "->", Position: 34},
				{TokenType: LITERAL, Value: "'data'", Position: 36},
				{TokenType: OP, Value: "->", Position: 42},
				{TokenType: LITERAL, Value: "'manifest'", Position: 44},
				{TokenType: OP, Value: "->", Position: 54},
				{TokenType: LITERAL, Value: "'metadata'", Position: 56},
				{TokenType: OP, Value: "->", Position: 66},
				{TokenType: LITERAL, Value: "'labels'", Position: 68},
				{TokenType: OP, Value: "->>", Position: 76},
				{TokenType: LITERAL, Value: "'foo'", Position: 79},
				{TokenType: OP, Value: "=", Position: 85},
				{TokenType: LITERAL, Value: "'bar'", Position: 87},
			},
		},
	}
	for _, testcase := range tests {
		tt := testcase

		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)

			scanner := NewSQLScanner()
			scanner.Init(tt.value)
			var allTokens []Token
			for scanner.Next() {
				allTokens = append(allTokens, *scanner.Token())
			}
			g.Expect(allTokens).To(gomega.Equal(tt.expectedTokens))
		})
	}
}

func Test_scanner_Peek(t *testing.T) {
	tests := []struct {
		name    string
		s       *scanner
		want    bool
		wantVal *Token
	}{
		{
			name: "return true and token if pos < length of tokens -1",
			s: &scanner{
				pos: 1,
				tokens: []Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
					{Value: "testToken3"},
				},
			},
			want: true,
			wantVal: &Token{
				Value: "testToken3",
			},
		},
		{
			name: "return false and nil if pos < length of tokens -1",
			s: &scanner{
				pos: 2,
				tokens: []Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
				},
			},
			want:    false,
			wantVal: nil,
		},
		{
			name: "return false and nil if pos == length of tokens -1",
			s: &scanner{
				pos: 1,
				tokens: []Token{
					{Value: "testToken1"},
					{Value: "testToken2"},
				},
			},
			want:    false,
			wantVal: nil,
		},
	}

	for _, testcase := range tests {
		tt := testcase

		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			got, gotVal := tt.s.Peek()
			g.Expect(got).To(gomega.Equal(tt.want))
			g.Expect(gotVal).To(gomega.Equal(tt.wantVal))
		})
	}
}
