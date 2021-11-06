package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "../comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMensajeDataLiderServer
}

var dn1 pb.DataNodeClient
var dn2 pb.DataNodeClient
var dn3 pb.DataNodeClient

var l1 [16][2]int
var l2 [16][2]int
var l3 [16][2]int

func (s *Server) Jugada(ctx context.Context, jugadaDataNode *pb.JugadaDataNode) (*pb.Empty, error) {
	nJug := strconv.Itoa(int(jugadaDataNode.Jugador))
	nRon := strconv.Itoa(int(jugadaDataNode.Ronda))
	//log.Printf("Recibido jugada, escogio el numero: %d", nEsc)
	frase := ("Jugador_") + nJug + (" Ronda_") + nRon
	f, err := os.Create("NameNodeData.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(frase + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("listo")
	return &pb.Empty{}, nil
}

func connectGRPC(puerto string) pb.DataNodeClient {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(puerto, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No Conectado a %s: %s", puerto, err)
	}
	defer conn.Close()

	c := pb.NewDataNodeClient(conn)
	return c
}

func Shuffle() {
	lista := [48][2]int{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, {12, 1}, {13, 1}, {14, 1}, {15, 1}, {16, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 2}, {9, 2}, {10, 2}, {11, 2}, {12, 2}, {13, 2}, {14, 2}, {15, 2}, {16, 2}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {7, 3}, {8, 3}, {9, 3}, {10, 3}, {11, 3}, {12, 3}, {13, 3}, {14, 3}, {15, 3}, {16, 3}}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(lista), func(i, j int) { lista[i], lista[j] = lista[j], lista[i] })
	l1 = lista[0:16]
	fmt.Println(l1)
}

func main() {

	fmt.Println("nameNode arriba")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	Shuffle()
	dn1 = connectGRPC(":9005")
	dn2 = connectGRPC(":9006")
	dn3 = connectGRPC(":9007")
	s := grpc.NewServer()

	pb.RegisterMensajeDataLiderServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
