package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	pb "../comms"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedPozoServer
}

// Anexo para tratar errores
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Revisa si el archivo del pozo existe
func fileExists() bool {
	if _, err := os.Stat("pozo.txt"); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		fmt.Println("wtf just happened (err 1)")
		fmt.Println(err.Error())
		return false
	}
}

// Crea el archivo del pozo
func createFile() {
	myfile, err := os.Create("pozo.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Se creó pozo.txt")
	myfile.Close()
}

// AGrega un jugador eliminado al archivo pozo
func addJugadorEliminado(nJugador int, nRonda int, prevMonto int) {
	text := strconv.Itoa(nJugador) + " " + strconv.Itoa(nRonda) + " " + strconv.Itoa(prevMonto+100000000)
	f, err := os.OpenFile("pozo.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

// Formatea los datos del archivo para entregarlos
func parseJugadorEliminadoString(texto string) (int, int) {
	separado := strings.Split(texto, " ")
	nJugador, err := strconv.Atoi(separado[0])
	if err != nil {
		// handle error
		fmt.Println(err)
		fmt.Println("err 4")
		os.Exit(2)
	}
	nRonda, err := strconv.Atoi(separado[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		fmt.Println("err 5")
		os.Exit(2)
	}
	return nJugador, nRonda
}

// Obtiene pozo
func getPozo() int {
	if fileExists() {
		content, err := ioutil.ReadFile("pozo.txt")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("err 2")
		}
		if len(content) == 0 {
			return 0
		}
		contentString := string(content)
		parsed := strings.Split(contentString, "\n")
		ultimaLinea := parsed[len(parsed)-1]
		pLParsed := strings.Split(ultimaLinea, " ")
		i, err := strconv.Atoi(pLParsed[2])
		if err != nil {
			// handle error
			fmt.Println(err)
			fmt.Println("err 3")
			os.Exit(2)
		}
		return i
	}
	fmt.Println("algo paso err 12389")
	return -1
}

// Maneja el pedido el pozo
func (s *Server) PedidoPozo(ctx context.Context, empty *pb.Empty) (*pb.MontoAcumulado, error) {

	return &pb.MontoAcumulado{Monto: int64(getPozo())}, nil
}

// Abre las conexiones para RPC
func opengRPC() {
	createFile()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9002))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterPozoServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// ABre las conexiones para RMQ
func openRMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@10.6.40.231:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			nJugador, nRonda := parseJugadorEliminadoString(string(d.Body))
			monto := getPozo()
			addJugadorEliminado(nJugador, nRonda, monto)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Main que corre todo
func main() {

	//cantRondasJuego1 := 1

	fmt.Println("Soy el Pozo!")

	opengRPC()
	openRMQ()

}
