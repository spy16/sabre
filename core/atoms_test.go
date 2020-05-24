package core

var (
	_ Value = Nil{}
	_ Value = Bool(true)
	_ Value = Int64(0)
	_ Value = Float64(0)
	_ Value = Character('a')
	_ Value = String("specimen")
	_ Value = Keyword("hello")
	_ Value = Symbol{Value: "def"}
)
