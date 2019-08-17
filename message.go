package dwlog

type Message struct {
	Time string
	Level Level
	Message string
}

func (m *Message) TextFormat(log *DWLog) []byte {
	text := "["+ m.Time +"]"+
			"["+ log.Host +"]"+
			"["+ m.Level.String() +"]"+
			" "+ m.Message

	return []byte(text)
}

func (m *Message) CSVFormat(log *DWLog) []byte {
	csv := m.Time +","+
			log.Host +","+
			m.Level.String() +","+
			m.Message

	return []byte(csv)
}