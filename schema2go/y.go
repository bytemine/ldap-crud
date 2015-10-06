//line items.y:2
package main

import __yyfmt__ "fmt"

//line items.y:2
// Documenting yacc code isn't easy. This mostly implements parsing of schema
// definitions like described in http://www.rfc-editor.org/rfc/rfc2252.txt

import (
	"strings"
)

type attribute struct {
	name int
	val  interface{}
}

var attributetypedefs = make(map[string]*attributetype)
var objectclassdefs = make(map[string]*objectclass, 0)

//line items.y:21
type yySymType struct {
	yys     int
	val     string
	oids    []string
	strings []string
	attrs   *attribute
	objs    *objectclass
	aattr   attribute
	aattrs  map[int]interface{}
	oattr   attribute
	oattrs  map[int]interface{}
}

const ATTRIBUTETYPE = 57346
const OBJECTCLASS = 57347
const NAME = 57348
const DESC = 57349
const OBSOLETE = 57350
const SUP = 57351
const EQUALITY = 57352
const ORDERING = 57353
const SUBSTR = 57354
const SYNTAX = 57355
const SINGLEVALUE = 57356
const COLLECTIVE = 57357
const NOUSERMODIFICATION = 57358
const ABSTRACT = 57359
const STRUCTURAL = 57360
const AUXILIARY = 57361
const MUST = 57362
const MAY = 57363
const DOLLAR = 57364
const LBRACE = 57365
const RBRACE = 57366
const USAGE = 57367
const EOF = 57368
const STRING = 57369
const OID = 57370
const NUMERICOID = 57371

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ATTRIBUTETYPE",
	"OBJECTCLASS",
	"NAME",
	"DESC",
	"OBSOLETE",
	"SUP",
	"EQUALITY",
	"ORDERING",
	"SUBSTR",
	"SYNTAX",
	"SINGLEVALUE",
	"COLLECTIVE",
	"NOUSERMODIFICATION",
	"ABSTRACT",
	"STRUCTURAL",
	"AUXILIARY",
	"MUST",
	"MAY",
	"DOLLAR",
	"LBRACE",
	"RBRACE",
	"USAGE",
	"EOF",
	"STRING",
	"OID",
	"NUMERICOID",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line items.y:215

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 42
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 100

var yyAct = [...]int{

	63, 60, 16, 17, 18, 19, 20, 21, 22, 23,
	24, 25, 26, 59, 57, 48, 13, 12, 58, 56,
	39, 27, 16, 17, 18, 19, 20, 21, 22, 23,
	24, 25, 26, 30, 31, 32, 33, 55, 29, 74,
	64, 27, 54, 47, 34, 35, 36, 37, 38, 69,
	67, 49, 68, 68, 62, 52, 46, 45, 65, 51,
	66, 30, 31, 32, 33, 44, 42, 50, 6, 7,
	41, 61, 34, 35, 36, 37, 38, 53, 43, 71,
	71, 73, 72, 71, 15, 70, 11, 5, 10, 4,
	8, 6, 7, 3, 2, 1, 9, 28, 14, 40,
}
var yyPact = [...]int{

	87, -1000, 64, -1000, -1000, -1000, 65, 63, -1000, -1000,
	-12, -13, 16, 55, -4, -1000, 43, 51, -1000, 37,
	29, 28, 15, -14, -1000, -1000, -1000, -1000, 27, -1000,
	32, 50, -1000, 14, -1000, -1000, -1000, -9, -10, -1000,
	-1000, -1000, 44, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 44, -1000, -1000, 12, -1000, 12, -1000, 12,
	26, -1000, 25, 61, -1000, 58, 57, -1000, -1000, -1000,
	-1000, 11, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 1, 0, 98, 97, 84, 38, 95, 94, 93,
	89, 87,
}
var yyR1 = [...]int{

	0, 7, 8, 8, 9, 9, 10, 11, 3, 3,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 4, 4, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 1, 1,
	2, 2,
}
var yyR2 = [...]int{

	0, 2, 2, 1, 1, 1, 5, 5, 2, 1,
	2, 4, 2, 1, 2, 2, 2, 2, 2, 1,
	1, 1, 1, 2, 1, 2, 4, 2, 1, 2,
	4, 1, 1, 1, 2, 4, 2, 4, 2, 1,
	3, 1,
}
var yyChk = [...]int{

	-1000, -7, -8, -9, -10, -11, 4, 5, 26, -9,
	23, 23, 29, 29, -3, -5, 6, 7, 8, 9,
	10, 11, 12, 13, 14, 15, 16, 25, -4, -6,
	6, 7, 8, 9, 17, 18, 19, 20, 21, 24,
	-5, 27, 23, 27, 28, 28, 28, 28, 29, 24,
	-6, 27, 23, 27, 28, 23, 28, 23, 28, 23,
	-1, 27, -1, -2, 28, -2, -2, 24, 27, 24,
	24, 22, 24, 24, 28,
}
var yyDef = [...]int{

	0, -2, 0, 3, 4, 5, 0, 0, 1, 2,
	0, 0, 0, 0, 0, 9, 0, 0, 13, 0,
	0, 0, 0, 0, 19, 20, 21, 22, 0, 24,
	0, 0, 28, 0, 31, 32, 33, 0, 0, 6,
	8, 10, 0, 12, 14, 15, 16, 17, 18, 7,
	23, 25, 0, 27, 29, 0, 34, 0, 36, 0,
	0, 39, 0, 0, 41, 0, 0, 11, 38, 26,
	30, 0, 35, 37, 40,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line items.y:83
		{
			//		attributetypedefs = append(attributetypedefs, newAttributetype($4))
			a := newAttributetype(yyDollar[4].aattrs)
			for _, v := range yyDollar[4].aattrs[NAME].([]string) {
				attributetypedefs[strings.ToLower(v)] = a
			}
			//attributetypedefs[$4[NAME].(string)] = newAttributetype($4)
		}
	case 7:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line items.y:93
		{
			objectclassdefs[strings.ToLower(yyDollar[4].oattrs[NAME].(string))] = newObjectclass(yyDollar[4].oattrs)
			//		objectclassdefs = append(objectclassdefs, newObjectclass($4))
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:99
		{
			yyVAL.aattrs[yyDollar[2].aattr.name] = yyDollar[2].aattr.val
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:102
		{
			yyVAL.aattrs = make(map[int]interface{})
			yyVAL.aattrs[yyDollar[1].aattr.name] = yyDollar[1].aattr.val
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:108
		{
			yyVAL.aattr = attribute{NAME, []string{yyDollar[2].val}}
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line items.y:111
		{
			yyVAL.aattr = attribute{NAME, yyDollar[3].strings}
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:114
		{
			yyVAL.aattr = attribute{DESC, yyDollar[2].val}
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:117
		{
			yyVAL.aattr = attribute{OBSOLETE, true}
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:120
		{
			yyVAL.aattr = attribute{SUP, yyDollar[2].val}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:123
		{
			yyVAL.aattr = attribute{EQUALITY, yyDollar[2].val}
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:126
		{
			yyVAL.aattr = attribute{ORDERING, yyDollar[2].val}
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:129
		{
			yyVAL.aattr = attribute{SUBSTR, yyDollar[2].val}
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:132
		{
			yyVAL.aattr = attribute{SYNTAX, yyDollar[2].val}
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:135
		{
			yyVAL.aattr = attribute{SINGLEVALUE, true}
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:138
		{
			yyVAL.aattr = attribute{COLLECTIVE, true}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:141
		{
			yyVAL.aattr = attribute{NOUSERMODIFICATION, true}
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:144
		{
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:149
		{
			yyVAL.oattrs[yyDollar[2].oattr.name] = yyDollar[2].oattr.val
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:152
		{
			yyVAL.oattrs = make(map[int]interface{})
			yyVAL.oattrs[yyDollar[1].oattr.name] = yyDollar[1].oattr.val
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:158
		{
			yyVAL.oattr = attribute{NAME, yyDollar[2].val}
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line items.y:162
		{
			yyVAL.oattr = attribute{NAME, yyDollar[3].strings[0]}
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:165
		{
			yyVAL.oattr = attribute{DESC, yyDollar[2].val}
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:168
		{
			yyVAL.oattr = attribute{OBSOLETE, true}
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:171
		{
			yyVAL.oattr = attribute{SUP, []string{yyDollar[2].val}}
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line items.y:174
		{
			yyVAL.oattr = attribute{SUP, yyDollar[3].oids}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:177
		{
			yyVAL.oattr = attribute{ABSTRACT, true}
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:180
		{
			yyVAL.oattr = attribute{STRUCTURAL, true}
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:183
		{
			yyVAL.oattr = attribute{AUXILIARY, true}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:186
		{
			yyVAL.oattr = attribute{MUST, []string{yyDollar[2].val}}
		}
	case 35:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line items.y:189
		{
			yyVAL.oattr = attribute{MUST, yyDollar[3].oids}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:192
		{
			yyVAL.oattr = attribute{MAY, []string{yyDollar[2].val}}
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line items.y:195
		{
			yyVAL.oattr = attribute{MAY, yyDollar[3].oids}
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line items.y:200
		{
			yyVAL.strings = append(yyVAL.strings, yyDollar[2].val)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:203
		{
			yyVAL.strings = []string{yyDollar[1].val}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line items.y:208
		{
			yyVAL.oids = append(yyVAL.oids, yyDollar[3].val)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line items.y:211
		{
			yyVAL.oids = []string{yyDollar[1].val}
		}
	}
	goto yystack /* stack new state and value */
}
