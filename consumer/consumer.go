package consumer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/datastorer"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/launcher"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/stopper"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Consumer interface {
	Setup([]string)
	Serve()
}

type BashCommandConsumer struct {
	Store    datastorer.DataStorer
	Launch   launcher.Launcher
	Stop     stopper.Stopper
	KafkaUrl []string
}

func (bcc *BashCommandConsumer) Setup(KafkaUrl []string) {
	bcc.KafkaUrl = KafkaUrl
	bcc.Store = datastorer.NewBashDataSorer()
	bcc.Launch = launcher.NewBashLauncher(bcc.Store)
	bcc.Stop = stopper.NewBashStopper(bcc.Store)
}

func (bcc *BashCommandConsumer) Serve() error {
	consumer, err := sarama.NewConsumer(bcc.KafkaUrl, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	runBashConsumer, err := consumer.ConsumePartition("RunBashAction", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer runBashConsumer.Close()

	BashStatusProducer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_URL")}, nil)
	if err != nil {
		return err
	}
	defer BashStatusProducer.Close()

	stopBashConsumer, err := consumer.ConsumePartition("StopBashAction", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer stopBashConsumer.Close()

	statRequestConsumer, err := consumer.ConsumePartition("StatRequest", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer statRequestConsumer.Close()

	statRequestProduser, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_URL")}, nil)
	if err != nil {
		return err
	}
	defer statRequestProduser.Close()

	for {
		select {
		// (обработка входящего сообщения и отправка ответа в Kafka)
		case msg, ok := <-runBashConsumer.Messages():
			if !ok {
				return errors.New("Connection closed")
			}
			execStatus := &sarama.ProducerMessage{
				Topic: "BashStatus",
				Key:   sarama.ByteEncoder(msg.Key),
				Value: sarama.StringEncoder("Launched"),
			}
			if err := bcc.Launch.Launch(msg); err != nil {
				execStatus.Value = sarama.StringEncoder(err.Error())
			}
			_, _, err = BashStatusProducer.SendMessage(execStatus)
			if err != nil {
				log.Println("Could not notify start")
				return err
			}
			log.Println("Start notification sent")
		case msg, ok := <-stopBashConsumer.Messages():
			if !ok {
				return errors.New("Connection closed")
			}
			stopStatus := &sarama.ProducerMessage{
				Topic: "BashStatus",
				Key:   sarama.ByteEncoder(msg.Key),
				Value: sarama.StringEncoder("Stoped"),
			}
			if err := bcc.Stop.Stop(msg); err != nil {
				stopStatus.Value = sarama.StringEncoder(err.Error())
			}
			_, _, err = BashStatusProducer.SendMessage(stopStatus)
			if err != nil {
				return err
			}
		case _, ok := <-statRequestConsumer.Messages():
			if !ok {
				return errors.New("connection closed")
			}
			v, _ := mem.VirtualMemory()
			c, _ := cpu.Times(false) //.Avg() //(false) //(time.Microsecond*2, true)
			fmt.Printf("Cpu: %v \n", (c[0].User+c[0].System)/(c[0].User+c[0].System+c[0].Idle)*100)
			fmt.Printf("UsedPercent:%f%%\n", v.UsedPercent)
			cpu := strconv.FormatFloat((c[0].User+c[0].System)/(c[0].User+c[0].System+c[0].Idle)*100, 'f', -1, 64)
			mem := strconv.FormatFloat(v.UsedPercent, 'f', -1, 64)
			status := &sarama.ProducerMessage{
				Topic: "StatResponse",
				Key:   sarama.ByteEncoder{},
				Value: sarama.StringEncoder("{\"cpu\": " + cpu + ", \"mem\":" + mem + "}"),
			}
			_, _, err = BashStatusProducer.SendMessage(status)
			if err != nil {
				return err
			}
		}
	}
}
