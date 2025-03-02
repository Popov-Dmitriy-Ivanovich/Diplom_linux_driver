package main

import (
	"os"
	"github.com/joho/godotenv"

	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/consumer"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	bashConsumer := consumer.BashCommandConsumer{}
	bashConsumer.Setup([]string{os.Getenv("KAFKA_URL")})
	print(bashConsumer.Serve().Error())
}