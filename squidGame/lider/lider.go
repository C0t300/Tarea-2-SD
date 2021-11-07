package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

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
	/* for Vivos < 15 {
		time.Sleep(1 * time.Second) // Ojala funcione [si no chao]
	} */
	return &pb.Jugador{NumJug: int32(Vivos)}, nil
}

func sendJugadorEliminadoPozo(nJugador int, nRonda int) {
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
	body := strconv.Itoa(nJugador) + " " + strconv.Itoa(nRonda)
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

func enviarDatosJugada(nJug int, nRonda int, nJugada int) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se conecta con el nameNode: %s", err)
	}
	defer conn.Close()

	c := pb.NewMensajeDataLiderClient(conn)
	response, err := c.Jugada(context.Background(), &pb.JugadaDataNode{Jugador: int32(nJug), Ronda: int32(nRonda), Jugada: int32(nJugada)})
	if err != nil {
		log.Fatalf("Error when calling nameNode: %s", err)
	}
	_ = response
}

func getEstadoEtapa1(r1 int, r2 int, r3 int, r4 int) (int, int) {
	suma := r1 + r2 + r3 + r4
	if r1 == 0 {
		return 1, suma
	} else if r2 == 0 {
		return 2, suma
	} else if r3 == 0 {
		return 3, suma
	} else if r4 == 0 {
		return 4, suma
	} else {
		return 0, suma
	}
}

func getHistorialJugador(nJug int) (int, int, int, int, int, int, int) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se conecta con el nameNode: %s", err)
	}
	defer conn.Close()

	c := pb.NewMensajeDataLiderClient(conn)
	response, err := c.HistorialJugador(context.Background(), &pb.Jugador{NumJug: int32(nJug)})
	if err != nil {
		log.Fatalf("Error when calling nameNode: %s", err)
	}
	fmt.Println("response ronda 1", response.Ronda1)
	return int(response.Jugador), int(response.Ronda1), int(response.Ronda2), int(response.Ronda3), int(response.Ronda4), int(response.Etapa2), int(response.Etapa3)
}

func (s *Server) Etapa1(ctx context.Context, jugadorEtapa1 *pb.JugadorEtapa1) (*pb.EstadoEtapa1, error) {
	nEsc := int(jugadorEtapa1.EscogidoJugador)
	numJug := int(jugadorEtapa1.NumJug)

	uno, r1, r2, r3, r4, tres, cuatro := getHistorialJugador(numJug)
	_ = uno
	_ = tres
	_ = cuatro
	nRonda, suma := getEstadoEtapa1(r1, r2, r3, r4)

	win := false
	if suma >= 21 {
		win = true
	}

	log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	nLider := rand.Intn(5) + 6
	alive := true
	if nLider <= nEsc {
		alive = false
		win = false
	}
	enviarDatosJugada(numJug, nRonda, nEsc)

	if nRonda >= 4 && suma < 21 {
		alive = false
		win = false
	}

	if !alive {
		sendJugadorEliminadoPozo(numJug, nRonda)
	}

	return &pb.EstadoEtapa1{Vivo: alive, EscogidoLider: int32(nLider), Win: win, Round: int32(nRonda)}, nil
}

func main() {

	fmt.Println("Soy el Lider!")
	Vivos = 0
	//q, errr, ch := openRMQ()
	// Estas variables se usan cada vez que se elimina alguien
	// Se debe llamar a sendJugadorEliminadoPozo()

	//parte cliente Lider-nameNode
	//parte Servidor Lider-Jugadores
	//cantRondasJuego1 := 1
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9003))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterJuego1Server(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
