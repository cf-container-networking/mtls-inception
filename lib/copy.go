package lib

import (
	"fmt"
	"io"
)

func CopyBoth(remoteConn io.ReadWriteCloser, localSource io.Reader, localSink io.Writer) error {
	errs := make(chan error)

	go func() {
		_, err := io.Copy(remoteConn, localSource)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(localSink, remoteConn)
		errs <- err
	}()

	err1 := <-errs
	err2 := <-errs
	remoteConn.Close()
	if err1 == nil && err2 == nil {
		return nil
	}
	return fmt.Errorf("copy errors: %s\n%s", err1, err2)
}
