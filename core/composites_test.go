package core

var (
	// Assert Set compatibilities
	_ Value      = (*Set)(nil)
	_ Expr       = (*Set)(nil)
	_ Seq        = (*Set)(nil)
	_ Comparable = (*Set)(nil)

	// Assert HashMap compatibilities
	_ Value = (*HashMap)(nil)
	_ Expr  = (*HashMap)(nil)

	// Assert Module compatibilities
	_ Value      = Module(nil)
	_ Expr       = Module(nil)
	_ Comparable = Module(nil)

	// Assert Vector compatibilities
	_ Value      = (*Vector)(nil)
	_ Expr       = (*Vector)(nil)
	_ Invokable  = (*Vector)(nil)
	_ Seq        = (*Vector)(nil)
	_ Comparable = (*Vector)(nil)

	// Assert List compatibilities
	_ Value      = (*List)(nil)
	_ Seq        = (*List)(nil)
	_ Expr       = (*List)(nil)
	_ Comparable = (*List)(nil)
)
