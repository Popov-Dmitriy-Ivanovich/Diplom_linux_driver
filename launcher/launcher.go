package launcher

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/Popov-Dmitriy-Ivanovich/Diplom_linux_driver/datastorer"
)

type Launcher interface {
	Launch(*sarama.ConsumerMessage) error
}

type BashLauncher struct {
	Store datastorer.DataStorer
}

func NewBashLauncher(store datastorer.DataStorer) BashLauncher {
	return BashLauncher{Store: store}
}

var bashFileUniqueIndex uint64 = 0

func (bl BashLauncher) Launch(msg *sarama.ConsumerMessage) error {
	key := msg.Key
	value := msg.Value

	ID, err := strconv.ParseUint(string(key), 16, 64)
	if err != nil {
		return err
	}

	filename := "bash/bash_" + strconv.FormatInt(time.Now().Unix(), 16) + "_" + string(key) + ".sh"

	if err := os.WriteFile(filename, value, 0777); err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", filename)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Printf("Command has NOT started started %v", err)
		return err
	}
	log.Printf("Command started successfully")
	bl.Store.Set(uint(ID), datastorer.BashCommandData{
		FilePath: filename,
		Cmd:      cmd,
	})

	return nil
}
