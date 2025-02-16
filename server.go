package main

import (

	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/consumer"
)



func main() {
	bashConsumer := consumer.BashCommandConsumer{}
	bashConsumer.Setup([]string{"localhost:9092"})
	print(bashConsumer.Serve().Error())
}