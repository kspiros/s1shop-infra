package xlib

//IMailSender auto email sender interface
type IMailSender interface {
}

//MailSender structure
type mailSender struct {
	host     string
	port     string
	password string
}

//NewMailSender mail sender constructor
func NewMailSender(host string, port string, password string) IMailSender {
	return &mailSender{
		host:     host,
		port:     port,
		password: password,
	}
}
