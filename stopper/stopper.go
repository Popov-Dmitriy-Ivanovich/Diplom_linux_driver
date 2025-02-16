package stopper

import (
	"errors"
	"os"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/datastorer"
)

type Stopper interface {
	Stop(*sarama.ConsumerMessage)
}

type BashStopper struct {
	Store datastorer.DataStorer
}

func NewBashStopper(store datastorer.DataStorer) BashStopper {
	return BashStopper{Store: store}
}

func (bs BashStopper) Stop(msg *sarama.ConsumerMessage) error {
	key := msg.Key
	// value := msg.Value

	ID, err := strconv.ParseUint(string(key),16,64)
	if err != nil {
		return err
	}

	
	data, err := bs.Store.Get(uint(ID)) 

	if err != nil {
		return err
	}

	cmdData, ok := data.(datastorer.BashCommandData)
	if !ok {
		return errors.New("передан неверный формат данных")
	}

	if err := cmdData.Cmd.Wait(); err != nil {
		return err
	}

	if err := os.Remove(cmdData.FilePath); err != nil {
		return err
	}

	return nil
}