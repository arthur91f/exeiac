package infra

import "os"

type StoreStdout struct {
	Output []byte
}

func (stdout *StoreStdout) Write(p []byte) (int, error) {
	stdout.Output = append(stdout.Output, p...)

	return len(p), nil
}

type StoreAndDisplayStdout struct {
	Output []byte
}

func (stdout *StoreAndDisplayStdout) Write(p []byte) (int, error) {
	stdout.Output = append(stdout.Output, p...)

	return os.Stdout.Write(p)
}

type StoreAndDisplayStderr struct {
	Output []byte
}

func (stderr *StoreAndDisplayStderr) Write(p []byte) (int, error) {
	stderr.Output = append(stderr.Output, p...)

	return os.Stderr.Write(p)
}
