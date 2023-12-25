package main

import (
	"client/internal/api"
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AskingDateTime - запрашивает время у удаленного сервера
// Сигнатура функции должна содержать параметры context.Context и Client нашего api сервиса - RandomClient
// Далее могут идти параметры из сигнатуры rpc метода
func AskingDateTime(ctx context.Context, m api.RandomClient) (*api.DateTime, error) {
	// Передаем клиентские данные для обращения к методу удаленного сервера GetDate
	// В данном случае это заглушка ввиде Value
	request := &api.RequestDateTime{
		Value: "Please send me the date and time",
	}

	// Запрашиваем дату и возвращаем ее
	return m.GetDate(ctx, request)
}

// AskRandom - запрашивает рандомное число у удаленного сервера
func AskRandom(ctx context.Context, m api.RandomClient, value int64) (*api.RandomInt, error) {
	request := &api.RandomParams{
		Value: value,
	}

	return m.GetRandom(ctx, request)
}

// AskPass - запрашивает пароль у удаленного сервера
func AskPass(ctx context.Context, m api.RandomClient, length int64) (*api.RandomPass, error) {
	request := &api.RequestPass{
		Length: length,
	}

	return m.GetRandomPass(ctx, request)
}

func main() {
	// Создаем соединение с GRPC сервером
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}

	// Генерируем случайное число
	newSource := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(newSource)

	client := api.NewRandomClient(conn) // Создаем экземпляр клиента сервиса

	r, err := AskingDateTime(context.Background(), client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server timestamp:", r.Value)

	i1, err := AskRandom(context.Background(), client, int64(newRand.Intn(100))+1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Random Integer 1:", i1.Value)

	i2, err := AskRandom(context.Background(), client, int64(newRand.Intn(100)))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Random Integer 2:", i2.Value)

	p, err := AskPass(context.Background(), client, int64(newRand.Intn(32)+8))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Random Password:", p.Password)
}
