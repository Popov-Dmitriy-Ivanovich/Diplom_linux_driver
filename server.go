package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/consumer"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	v, _ := mem.VirtualMemory()
	c, _ := cpu.Times(false) //.Avg() //(false) //(time.Microsecond*2, true)
	fmt.Printf("Cpu: %v \n", (c[0].User+c[0].System)/(c[0].User+c[0].System+c[0].Idle)*100)
	fmt.Printf("UsedPercent:%f%%\n", v.UsedPercent)

	bashConsumer := consumer.BashCommandConsumer{}
	bashConsumer.Setup([]string{os.Getenv("KAFKA_URL")})
	print(bashConsumer.Serve().Error())
}
