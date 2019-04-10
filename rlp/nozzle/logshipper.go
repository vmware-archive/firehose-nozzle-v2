package nozzle

import "io"

type LogShipper interface {
	LogShip(string) error
}

type SampleShipper struct {
	writer io.Writer
}

func NewSampleShipper(writer io.Writer) LogShipper {
	return &SampleShipper{
		writer: writer,
	}
}

func (ss *SampleShipper) LogShip(log string) error {
	_, err := ss.writer.Write([]byte(log))
	if err != nil {
		return err
	}

	_, err = ss.writer.Write([]byte("\n"))
	return err
}
