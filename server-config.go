package xlib

var deferFunc []func()

type ServerConfig struct {
	DB     IDBConn
	Bus    IBus
	Logger ILogger
	Email  IMailSender
	Mem    IMemCash
	isok   bool
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{isok: true}
}

func (sc *ServerConfig) Close() {
	for i := 0; i < len(deferFunc); i++ {
		deferFunc[i]()
	}
}

func (sc *ServerConfig) Build() *ServerConfig {
	if !sc.isok {
		return nil
	}
	return sc
}

func (sc *ServerConfig) WithBus() *ServerConfig {
	nc, bustidy, err := NewBus()
	if err != nil {
		sc.isok = false
	} else {
		deferFunc = append(deferFunc, bustidy)
	}
	sc.Bus = nc
	defer bustidy()
	return sc
}

func (sc *ServerConfig) WithIMailSender() *ServerConfig {
	//sc.Email = email
	return sc
}

func (sc *ServerConfig) WithLogger() *ServerConfig {
	l, tidylogger, err := NewLogger()
	if err != nil {
		sc.isok = false
	} else {
		deferFunc = append(deferFunc, tidylogger)
	}
	sc.Logger = l
	return sc
}

func (sc *ServerConfig) WithMemCash() *ServerConfig {
	mem, err := NewMemCash()
	if err != nil {
		sc.isok = false
	}
	sc.Mem = mem
	return sc
}

func (sc *ServerConfig) WithDB() *ServerConfig {

	db, dbtidy, err := InitDatabase()
	if err != nil {
		sc.isok = false
	} else {
		deferFunc = append(deferFunc, dbtidy)
	}
	sc.DB = db
	return sc
}
