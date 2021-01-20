package snowflake

var def = New(0)

func Config(node uint) {
	def.Set(node)
}

func Generate() ID {
	return def.Generate()
}
