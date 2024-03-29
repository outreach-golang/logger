package gorm_V2

type SqlInfo map[string][]interface{}

func (i SqlInfo) Set(key string, val interface{}) {
	i[key] = []interface{}{val}
}

func (i SqlInfo) Get(key string) interface{} {
	return i[key][0]
}
