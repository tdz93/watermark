package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tdz93/watermark/api/v1/pb"
	watermark "github.com/tdz93/watermark/pkg"
	"github.com/tdz93/watermark/pkg/endpoints"
	"github.com/tdz93/watermark/pkg/transport"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

func main() {
	var (
		logger log.Logger
		//configuriamo gli indirizzi per i server HTTP e gRPC
		//net.JoinHostPort() è una funzione che combina un hostname e una porta in un singolo indirizzo. Per esempio, unisce "localhost" e "8081" in "localhost:8081"
		//envString() è una funzione helper definita alla fine del file
		httpAddr = net.JoinHostPort("0.0.0.0", envString("HTTP_PORT", defaultHTTPPort))
		grpcAddr = net.JoinHostPort("0.0.0.0", envString("GRPC_PORT", defaultGRPCPort))
	)

	//log.NewLogfmtLogger() crea un nuovo logger che usa il formato logfmt (un formato chiave-valore comune per i log)
	//log.NewSyncWriter(os.Stderr) crea un writer sincronizzato che scrive su stderr (output di errore standard)
	//il writer sincronizzato garantisce che le operazioni di scrittura dei log siano thread-safe
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	//log.With() aggiunge dei campi predefiniti che appariranno in ogni messaggio di log
	//"ts" è la chiave per il timestamp
	//log.DefaultTimestampUTC è una funzione che fornisce il timestamp in formato UTC
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var (
		service     = watermark.NewService()
		eps         = endpoints.NewEndpointSet(service)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	//group.Group è uno strumento che permette di:
	//	-gestire più goroutine contemporaneamente in modo coordinato
	//	-gestire il ciclo di vita dei vari componenti del servizio
	//	-implementare una chiusura ordinata (graceful shutdown) del servizio
	//di seguito viene utilizzato per gestire tre componenti principali:
	//	-il server HTTP
	//	-il server gRPC
	//	-un gestore di segnali per lo shutdown controllato
	//ogni componente viene aggiunto al gruppo con due funzioni:
	//	-una funzione di esecuzione (func() error)
	//	-una funzione di interruzione (func(error))
	//questo pattern è particolarmente utile perché:
	//	-garantisce che tutti i componenti vengano avviati e fermati in modo coordinato
	//	-gestisce correttamente gli errori di qualsiasi componente
	//	-permette una chiusura pulita del servizio quando riceve un segnale di interruzione (es. CTRL+C)
	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by the here demonstrated zipkin tracing middleware.
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterWatermarkServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

// funzione helper, cerca una variabile d'ambiente spechificata(HTTP_PORT o GRPC_PORT), se esiste ritorna quella, se no fallback
// approccio molto comune nei microservizi, permette di:
//
//	-avere valori predefiniti per lo sviluppo locale
//	-configurare facilmente porte diverse in ambienti di produzione attraverso variabili d'ambiente
//	-evitare conflitti di porta quando si eseguono più istanze del servizio
func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
