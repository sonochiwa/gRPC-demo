package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"server/internal/api"

	"google.golang.org/grpc"
)

// Генерируем случайный пароль
func genPassword(len int64) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, len)

	for i := range b {
		b[i] = letters[rand.Intn(int(len))]
	}
	return string(b)
}

// RandomServer - структура, названная в честь сервиса gRPC.
// Будет реализовывать интерфейс, требуемый сервером gRPC.
type RandomServer struct {
	api.UnimplementedRandomServer
}

// GetDate - метод, который возвращает клиенту дату и время сервера
// Принимает context.Context и указатель на сообщение RequestDateTime. Возвращает сообщение DateTime
func (RandomServer) GetDate(ctx context.Context, r *api.RequestDateTime) (*api.DateTime, error) {
	t := time.Now()

	// Форматируем время в iso8601
	currentTime := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02dZ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),
	)

	response := &api.DateTime{
		Value: currentTime,
	}

	return response, nil
}

// GetRandom - метод, который возвращает клиенту рандомное число
func (RandomServer) GetRandom(ctx context.Context, r *api.RandomParams) (*api.RandomInt, error) {
	// Переменные newSource и newRand используются для генерации случайного числа
	newSource := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(newSource)

	// На основе клиенского сообщения генерируется случайное число
	response := &api.RandomInt{
		Value: int64(newRand.Intn(int(r.GetValue()))),
	}

	return response, nil
}

func (RandomServer) GetRandomPass(ctx context.Context, r *api.RequestPass) (*api.RandomPass, error) {
	temp := genPassword(r.GetLength()) // Генерируем случайный пароль

	response := &api.RandomPass{
		Password: temp,
	}

	return response, nil
}

func main() {
	// Объявляем структуру Random реализующую наши методы
	var randomServer RandomServer

	// Создаем новый gRPC сервер
	server := grpc.NewServer()

	// Регистрируем новый сервер и структуру
	api.RegisterRandomServer(server, randomServer)

	// Слушаем сеть по указанному адресу
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Serving requests...")
	err = server.Serve(listen) // Запускаем gRPC сервер
	if err != nil {
		fmt.Println(err)
		return
	}
}
