package hello

import (
	"expvar"
	"github.com/nsqio/go-nsq"
	logging "gopkg.in/tokopedia/logging.v1"
	"log"
	"net/http"
	"os"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type HelloWorldModule struct {
	cfg       *Config
	q         *nsq.Consumer
	something string
	stats     *expvar.Int
}

func NewHelloWorldModule() *HelloWorldModule {

	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfg.Server.Name)

	// contohnya: caranya ciptakan nsq consumer
	nsqCfg := nsq.NewConfig()
	q := createNewConsumer(nsqCfg, "random-topic", "test", handler)
	q.SetLogger(log.New(os.Stderr, "nsq:", log.Ltime), nsq.LogLevelError)
	q.ConnectToNSQLookupd("nsqlookupd.local:4161")

	return &HelloWorldModule{
		cfg:       &cfg,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
		q:         q,
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	w.Write([]byte("Hello " + hlm.something))
}

func handler(msg *nsq.Message) error {
	log.Println("got message :", string(msg.Body))
	msg.Finish()
	return nil
}

func createNewConsumer(nsqCfg *nsq.Config, topic string, channel string, handler nsq.HandlerFunc) *nsq.Consumer {
	q, err := nsq.NewConsumer(topic, channel, nsqCfg)
	if err != nil {
		log.Fatal("failed to create consumer for ", topic, channel, err)
	}
	q.AddHandler(handler)
	return q
}
