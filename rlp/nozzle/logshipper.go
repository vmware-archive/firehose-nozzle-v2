package nozzle

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"encoding/json"
	"io"
)

type LogShipper interface {
	LogShip(log *loggregator_v2.Envelope) error
}

type SampleShipper struct {
	writer io.Writer
}

func NewSampleShipper(writer io.Writer) LogShipper {
	return &SampleShipper{
		writer: writer,
	}
}

func (ss *SampleShipper) LogShip(log *loggregator_v2.Envelope) error {
	logBytes, _ := json.Marshal(log)
	_, err := ss.writer.Write(logBytes)
	if err != nil {
		return err
	}

	_, err = ss.writer.Write([]byte("\n"))
	return err
}
