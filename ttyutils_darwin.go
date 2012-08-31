package ttyutils

import (
	"syscall"
	"errors"
	"fmt"
	"os"
	"unsafe"
)

const (
	sys_ISTRIP = 0x20
	sys_INLCR = 0x40
	sys_ICRNL = 0x100
	sys_IGNCR = 0x80
	sys_IXON = 0x200
	sys_IXOFF = 0x400
	sys_ICANON = 0x100
	sys_ISIG = 0x80
	termios_NCCS = 20
)

type tcflag_t uint64 // unsigned long
type cc_t byte       // unsigned char
type speed_t uint64  // unsigned long

type Termios struct {
	Iflag tcflag_t       /* input flags */
	Oflag tcflag_t       /* output flags */
	Cflag tcflag_t       /* control flags */
	Lflag tcflag_t       /* local flags */
	Cc[termios_NCCS] cc_t /* control chars */
	Ispeed speed_t      /* input speed */
	Ospeed speed_t      /* output speed */
}

func IsTerminal(fd uintptr) bool {
	var termios Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func MirrorWinsize(from, to *os.File) error {
	var n int
	err := ioctl(from.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return err
	}
	err = ioctl(to.Fd(), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return err
	}
	return nil
}

func ioctl(fd uintptr, cmd uintptr, ptr uintptr) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		cmd,
		uintptr(unsafe.Pointer(ptr)),
	)
	if e != 0 {
		return errors.New(fmt.Sprintf("ioctl failed! %s", e))
	}
	return nil
}

func MakeTerminalRaw(fd uintptr) (*Termios, error) {
	var s Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&s)), 0, 0, 0); err != 0 {
		return nil, err
	}

	oldState := s
	s.Iflag &^= sys_ISTRIP | sys_INLCR | sys_ICRNL | sys_IGNCR | sys_IXON | sys_IXOFF
	s.Lflag &^= syscall.ECHO | sys_ICANON | sys_ISIG
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&s)), 0, 0, 0); err != 0 {
		return nil, err
	}

	return &oldState, nil
}

func RestoreTerminalState(fd uintptr, termios *Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	return err
}

