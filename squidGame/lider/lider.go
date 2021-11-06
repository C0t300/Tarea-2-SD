package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "../comms"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedJuego1Server
}

var Vivos int32

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (s *Server) QuieroJugar(ctx context.Context, empty *pb.Empty) (*pb.Jugador, error) {
	Vivos += 1
	return &pb.Jugador{NumJug: int32(Vivos)}, nil
}

func (s *Server) Jugada1(ctx context.Context, jugadorJuego1 *pb.JugadorJuego1) (*pb.EstadoJuego, error) {
	nEsc := int(jugadorJuego1.EscogidoJugador)
	log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	nLider := rand.Intn(5) + 6
	dead := false
	if nLider >= nEsc {
		dead = true
	}

	return &pb.EstadoJuego{Vivo: dead, EscogidoLider: int32(nLider), Win: false, Round: int32(1)}, nil
}

func openRMQ() (amqp.Queue, error, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	return q, err, ch
}

func sendJugadorEliminadoPozo(q amqp.Queue, err error, ch *amqp.Channel, nJugador int, nRonda int) {
	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func main() {

	fmt.Println("Soy el Lider!")
	Vivos = 0
	//q, errr, ch := openRMQ()
	// Estas variables se usan cada vez que se elimina alguien
	// Se debe llamar a sendJugadorEliminadoPozo()

	//parte cliente Lider-nameNode
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se conecta con el nameNode: %s", err)
	}
	defer conn.Close()

	c := pb.NewMensajeDataLiderClient(conn)

	response, err := c.Jugada(context.Background(), &pb.JugadaDataNode{Jugador: 1, Ronda: 1, Jugada: 1})
	_ = response
	if err != nil {
		log.Fatalf("Error when calling nameNode: %s", err)
	}
	log.Printf("Respuesta desde nameNode.")
	//parte Servidor Lider-Jugadores
	//cantRondasJuego1 := 1
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterJuego1Server(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
