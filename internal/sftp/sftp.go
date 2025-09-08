package sftp

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Unfield/FileHopper/internal/auth"
	"github.com/Unfield/FileHopper/internal/db"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPServer struct {
	addr          string
	authenticator *auth.Authenticator
	db            db.DBDriver
	fsRoot        string
}

func NewSFTPServer(addr string, authenticator *auth.Authenticator, db *db.DBDriver, rootDir string) *SFTPServer {
	return &SFTPServer{
		addr:          addr,
		authenticator: authenticator,
		fsRoot:        rootDir,
		db:            *db,
	}
}

func (s *SFTPServer) Start() error {
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			log.Printf("Login attempt by %s", conn.User())
			if s.authenticator.Authenticate(conn.User(), string(pass)) {
				log.Printf("User %s authenticated", conn.User())
				user, err := s.db.GetUser(conn.User())
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve user '%s'", conn.User())
				}

				homeDir := fmt.Sprintf("%s/%s", s.fsRoot, user.HomeDir)

				if err := os.MkdirAll(homeDir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create home dir: %v", err)
				}

				return &ssh.Permissions{
					Extensions: map[string]string{
						"home": homeDir,
					},
				}, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		},
	}

	privateBytes, err := os.ReadFile("./server_key")
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	config.AddHostKey(private)

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}
	log.Printf("SFTP server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn, config)
	}
}

func (s *SFTPServer) handleConn(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Println("Failed handshake:", err)
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "only session channels are allowed")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Println("Could not accept channel:", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
					log.Println("Starting SFTP subsystem")
					req.Reply(true, nil)
				} else {
					req.Reply(false, nil)
				}
			}
		}(requests)

		homeDir := sshConn.Permissions.Extensions["home"]
		fs := NewChrootFS(homeDir)

		handlers := sftp.Handlers{
			FileGet:  fs,
			FilePut:  fs,
			FileCmd:  fs,
			FileList: fs,
		}

		server := sftp.NewRequestServer(channel, handlers)

		if err := server.Serve(); err == io.EOF {
			_ = server.Close()
			log.Println("SFTP client disconnected")
		} else if err != nil {
			log.Println("SFTP server completed with error:", err)
		}
	}
}
