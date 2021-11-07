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

var l1 [][2]int
var l2 [][2]int
var l3 [][2]int

var ip1 string
var ip2 string
var ip3 string

func enLista(jugador int, ronda int, lista [][2]int) bool {
	encontrar := [2]int{jugador, ronda}
	for _, valores := range lista {
		if valores == encontrar {
			return true
		}
	}
	return false
}

func (s *Server) Jugada(ctx context.Context, jugadaDataNode *pb.JugadaDataNode) (*pb.Empty, error) {
	nJug := strconv.Itoa(int(jugadaDataNode.Jugador))
	nRon := strconv.Itoa(int(jugadaDataNode.Ronda))
	var ip string

	if enLista(int(jugadaDataNode.Jugador), int(jugadaDataNode.Ronda), l1) {
		EnviarDatos(dn1, *jugadaDataNode)
		ip = ip1
	} else if enLista(int(jugadaDataNode.Jugador), int(jugadaDataNode.Ronda), l2) {
		EnviarDatos(dn2, *jugadaDataNode)
		ip = ip2
	} else {
		EnviarDatos(dn3, *jugadaDataNode)
		ip = ip3
	}
	frase := ("Jugador_") + nJug + (" Ronda_") + nRon + " " + ip
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

func (s *Server) HistorialJugador(ctx context.Context, jugador *pb.Jugador) (*pb.HistorialJugadas, error) {
	var ronda1 *pb.HistorialJugadas
	/* var ronda2 pb.HistorialJugadas
	var ronda3 pb.HistorialJugadas */

	/* ronda1, err := dn1.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 1})
	if err != nil {
		log.Fatalf("Error when calling PedirDatos: %s", err)
	}
	ronda2, err := dn2.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 2})
	if err != nil {
		log.Fatalf("Error when calling PedirDatos: %s", err)
	}
	ronda3, err := dn3.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 3})
	if err != nil {
		log.Fatalf("Error when calling PedirDatos: %s", err)
	} */

	if enLista(int(jugador.NumJug), 1, l1) {
		ronda1, _ = dn1.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 1})
	} else if enLista(int(jugador.NumJug), 1, l2) {
		ronda1, _ = dn2.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 1})
	} else {
		ronda1, _ = dn3.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 1})
	}
	/* if enLista(int(jugador.NumJug), 1, l1) {
		ronda1, err = dn1.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 2})
	} else if enLista(int(jugador.NumJug), 1, l2) {
		ronda1, err = dn2.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 2})
	} else if enLista(int(jugador.NumJug), 1, l3) {
		ronda1, err = dn3.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 2})
	} else {
		fmt.Println("Error 12039784")
	}

	if enLista(int(jugador.NumJug), 1, l1) {
		ronda1, err = dn1.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 3})
	} else if enLista(int(jugador.NumJug), 1, l2) {
		ronda1, err = dn2.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 3})
	} else if enLista(int(jugador.NumJug), 1, l3) {
		ronda1, err = dn3.PedirDatos(context.Background(), &pb.JugadorEtapa{Jugador: jugador.NumJug, Etapa: 3})
	} else {
		fmt.Println("Error 12039784")
	} */
	/* message HistorialJugadas{
		int32 jugador=1;
		int32 ronda1=2;
		int32 ronda2=3;
		int32 ronda3=4;
		int32 ronda4=5;
		int32 etapa2=6;
		int32 etapa3=7;
	  } */

	return &pb.HistorialJugadas{Jugador: int32(jugador.NumJug), Ronda1: int32(ronda1.Ronda1), Ronda2: int32(ronda1.Ronda2), Ronda3: int32(ronda1.Ronda3), Ronda4: int32(ronda1.Ronda4), Etapa2: 0, Etapa3: 0}, nil
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

func EnviarDatos(c pb.DataNodeClient, datos pb.JugadaDataNode) {
	fmt.Printf("Jugada enviada al DataNode.")
	response, err := c.GuardarDatos(context.Background(), &datos)
	_ = response
	if err != nil {
		log.Fatalf("Error when calling GuardarDatos: %s", err)
	}
}

func Shuffle() {
	lista := [48][2]int{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, {12, 1}, {13, 1}, {14, 1}, {15, 1}, {16, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 2}, {9, 2}, {10, 2}, {11, 2}, {12, 2}, {13, 2}, {14, 2}, {15, 2}, {16, 2}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {7, 3}, {8, 3}, {9, 3}, {10, 3}, {11, 3}, {12, 3}, {13, 3}, {14, 3}, {15, 3}, {16, 3}}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(lista), func(i, j int) { lista[i], lista[j] = lista[j], lista[i] })
	l1 = lista[0:16]
	l2 = lista[17:32]
	l3 = lista[33:48]
}

func main() {

	fmt.Println("nameNode arriba")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	Shuffle()
	ip1 = ""
	ip2 = ""
	ip3 = ""
	dn1 = connectGRPC(ip1 + ":9005")
	dn2 = connectGRPC(ip2 + ":9006")
	dn3 = connectGRPC(ip3 + ":9007")
	s := grpc.NewServer()

	pb.RegisterMensajeDataLiderServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
