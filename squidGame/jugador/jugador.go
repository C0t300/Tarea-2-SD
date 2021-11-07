package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	pb "../comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/* func connectGRPC() pb.Juego1Client {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewJuego1Client(conn)
	return c
} */

// Envia al lider una peticion para jugar
func QuieroJugarJugador() int {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.6.40.231:9003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewJuego1Client(conn)
	fmt.Printf("Enviada solicitud de juego.")
	response, err := c.QuieroJugar(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error when calling QuieroJugarJugador: %s", err)
	}
	return int(response.NumJug)
}

func etapa1(nJug int) (bool, bool) { // keep loop and win
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.6.40.231:9003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewJuego1Client(conn)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese numero para Luz Roja, Luz Verde.")
	fmt.Print("Recuerde que debe ser entre 1 y 10: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	num, _ := strconv.Atoi(string(text))

	response, err := c.Etapa1(context.Background(), &pb.JugadorEtapa1{EscogidoJugador: int32(num), NumJug: int32(nJug)})
	if err != nil {
		log.Fatalf("Error when calling Etapa1: %s", err)
	}
	lider := response.EscogidoLider
	text = fmt.Sprintf("El lider escogio %d", lider)
	fmt.Println(text)
	if !response.Vivo {
		fmt.Println("Haz muerto.")
		return false, false
	} else if response.Win {
		fmt.Println("Ganaste!")
		return false, true
	} else {
		text = fmt.Sprintf("Pasando a la ronda %d", response.Round)
		fmt.Println(text)
		return true, false
	}

}

func main() {

	nJugador := QuieroJugarJugador()
	fmt.Printf("Juego iniciado. Jugador numero: %d", nJugador)
	loop, win := etapa1(nJugador)
	for loop {
		loop, win = etapa1(nJugador)
	}
	if !win {
		os.Exit(0)
	}

}
