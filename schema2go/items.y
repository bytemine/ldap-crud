%{
package main

// Documenting yacc code isn't easy. This mostly implements parsing of schema
// definitions like described in http://www.rfc-editor.org/rfc/rfc2252.txt

import (
	"strings"
)

type attribute struct {
	name	int
	val	interface{}
}

var attributetypedefs = make(map[string]*attributetype)
var objectclassdefs = make(map[string]*objectclass,0)

%}

%union {
	val	string
	oids	[]string
	strings	[]string
	attrs	*attribute
	objs	*objectclass
	aattr	attribute
	aattrs	map[int]interface{}
	oattr	attribute
	oattrs 	map[int]interface{}
}

/* tokens without values */
%token ATTRIBUTETYPE
%token OBJECTCLASS
%token NAME
%token DESC
%token OBSOLETE
%token SUP
%token EQUALITY
%token ORDERING
%token SUBSTR
%token SYNTAX
%token SINGLEVALUE
%token COLLECTIVE
%token NOUSERMODIFICATION
%token ABSTRACT
%token STRUCTURAL
%token AUXILIARY
%token MUST
%token MAY
%token DOLLAR
%token LBRACE
%token RBRACE
%token USAGE
%token EOF

/* tokens with values */
%token STRING
%token OID
%token NUMERICOID

%type <val> STRING OID NUMERICOID
%type <strings> strings
%type <oids> oids
%type <aattrs> aattrs
%type <oattrs> oattrs
%type <aattr> aattr
%type <oattr> oattr

%%
start: 	defs EOF
	;

defs:		defs def
	|	def
	;

def:		adef
	|	odef
	;

adef: 	ATTRIBUTETYPE LBRACE NUMERICOID aattrs RBRACE {
//		attributetypedefs = append(attributetypedefs, newAttributetype($4))
		a := newAttributetype($4)
		for _, v := range $4[NAME].([]string) {
			attributetypedefs[strings.ToLower(v)] = a
		}
		//attributetypedefs[$4[NAME].(string)] = newAttributetype($4)
	}
	;

odef:	OBJECTCLASS LBRACE NUMERICOID oattrs RBRACE {
		objectclassdefs[strings.ToLower($4[NAME].(string))] = newObjectclass($4)
//		objectclassdefs = append(objectclassdefs, newObjectclass($4))
	}
	;

aattrs:		aattrs aattr {
			$$[$2.name] = $2.val
		}
	|	aattr {
			$$ = make(map[int]interface{})
			$$[$1.name] = $1.val
		}
	;

aattr:		NAME STRING {
			$$ = attribute{NAME, []string{$2}}
		}
	|	NAME LBRACE strings RBRACE {
			$$ = attribute{NAME, $3}
		}
	|	DESC STRING {
			$$ = attribute{DESC, $2}
		}
	|	OBSOLETE {
			$$ = attribute{OBSOLETE, true}
		}
	|	SUP OID {
			$$ = attribute{SUP, $2}
		}
	|	EQUALITY OID {
			$$ = attribute{EQUALITY, $2}
		}
	|	ORDERING OID {
			$$ = attribute{ORDERING, $2}
		}
	|	SUBSTR OID {
			$$ = attribute{SUBSTR, $2}
		}
	|	SYNTAX NUMERICOID {
			$$ = attribute{SYNTAX, $2}
		}
	|	SINGLEVALUE {
			$$ = attribute{SINGLEVALUE, true}
		}
	|	COLLECTIVE {
			$$ = attribute{COLLECTIVE, true}
		}
	|	NOUSERMODIFICATION {
			$$ = attribute{NOUSERMODIFICATION, true}
		}
	|	USAGE {
		}
	;


oattrs:		oattrs oattr {
			$$[$2.name] = $2.val
		}
	|	oattr {
			$$ = make(map[int]interface{})
			$$[$1.name] = $1.val
		}
	;

oattr:		NAME STRING {
			$$ = attribute{NAME, $2}
		}
	|
		NAME LBRACE strings RBRACE {
			$$ = attribute{NAME, $3[0]}	
		}
	|	DESC STRING {
			$$ = attribute{DESC, $2}
		}
	|	OBSOLETE {
			$$ = attribute{OBSOLETE, true}
		}
	|	SUP OID {
			$$ = attribute{SUP, []string{$2}}
		}
	|	SUP LBRACE oids RBRACE {
			$$ = attribute{SUP, $3}
		}
	|	ABSTRACT {
			$$ = attribute{ABSTRACT, true}
		}
	|	STRUCTURAL {
			$$ = attribute{STRUCTURAL, true}
		}
	|	AUXILIARY {
			$$ = attribute{AUXILIARY, true}
		}
	|	MUST OID {
			$$ = attribute{MUST, []string{$2}}
		}
	|	MUST LBRACE oids RBRACE {
			$$ = attribute{MUST, $3}
		}
	|	MAY OID {
			$$ = attribute{MAY, []string{$2}}
		}
	|	MAY LBRACE oids RBRACE {
			$$ = attribute{MAY, $3}
		}
	;

strings:	strings STRING {
			$$ = append($$, $2)
		}
	|	STRING {
			$$ = []string{$1}
		}
	;

oids:		oids DOLLAR OID {
			$$ = append($$, $3)
		}
	|	OID {
			$$ = []string{$1}
		}
	;
%%

