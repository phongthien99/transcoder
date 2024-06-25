package adapter

import (
	"example/src/module/origin/model"
	"net"

	"github.com/pkg/sftp"
	"github.com/spf13/afero"
	"github.com/spf13/afero/sftpfs"
	"golang.org/x/crypto/ssh"
)

func NewSftps(config *model.FTPConfig) (afero.Fs, error) {
	configSsh := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Dial your ssh server.
	conn, errSSH := ssh.Dial("tcp", config.Host, configSsh)
	if errSSH != nil {
		return nil, errSSH
	}

	client, errSftp := sftp.NewClient(conn)
	if errSftp != nil {
		return nil, errSSH
	}

	return sftpfs.New(client), nil
}
