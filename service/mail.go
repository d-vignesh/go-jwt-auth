package service

type Message struct {
	To 		[]string
	From	string
	Subject string
	Body 	string
	User	string
	Type 	string
	Massive bool
	Info 	string
}

