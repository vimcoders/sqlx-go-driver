package sqlx

type Table interface {
	TableName() string
}

type Convertor interface {
	Convert() (sql string, args []interface{})
}

type Marshaler interface {
	Marshal() (str string, err error)
}

type Unmarshaler interface {
	Unmarshal(str string) error
}
