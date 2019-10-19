package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
)

const (
	server            = "rps.vhenne.de:6000"
	maxTimePerRound   = 2500 * time.Millisecond // limit 3s, but we have some lag
	logServerMessages = false
)

type Player struct {
	Name     string
	Password string
	Play     func(tick <-chan *Tick)
}

type Tick struct {
	Gamestate *rules.Gamestate
	Ctx       context.Context // with deadline
	Move      chan<- rules.Move
}

func Main(player *Player) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	var shouldShutDown uint64
	go func() {
		log.Printf("received signal %s, terminating after current game", <-sigChan)
		signal.Stop(sigChan)
		atomic.StoreUint64(&shouldShutDown, 1)
	}()
	for {
		if atomic.LoadUint64(&shouldShutDown) != 0 {
			log.Print("terminating")
			return
		}
		func() {
			ticks := make(chan *Tick, 1)
			defer close(ticks)
			go func() {
				player.Play(ticks)
			}()
			conn, err := net.Dial("tcp", server)
			defer func() {
				err = conn.Close()
				if err != nil {
					log.Printf("error closing connection: %s", err)
				}
			}()
			if err != nil {
				log.Printf("could not connect to server: %v", err)
				return
			}
			//login
			_, err = fmt.Fprintf(conn, "login %s %s\n", player.Name, player.Password)
			if err != nil {
				log.Printf("could not login: %s", err)
				return
			}
			log.Print("logged in, waiting for game")
			bufReader := bufio.NewReader(conn)
			for round := 0; ; round++ {
				message, err := bufReader.ReadBytes('\n')
				switch err {
				case io.EOF:
					// connection closed
					// may still have received one last gamestate without trailing newline?
					if len(message) == 0 {
						return
					}
				case nil:
				default:
					log.Printf("round %d: error reading new gamestate: %s", round, err)
					return
				}
				if len(message) == 0 || message[0] != '{' { // ignore non-gamestate messages
					if logServerMessages {
						log.Printf("round %d: server message: %s", round, message)
					}
					round-- // this round doesn't count
					continue
				}
				stop := func() bool {
					ctx, cancel := context.WithTimeout(context.Background(), maxTimePerRound)
					defer cancel()
					var gs gamestate
					err = json.Unmarshal(message, &gs)
					if err != nil {
						log.Printf("round %d: error unmarshaling gamestate: %s", round, err)
						return true
					}
					if round == 0 {
						log.Printf("starting game with players %v", gs.Players)
					}
					if gs.GameOver {
						if gs.Winner != nil {
							log.Printf("round %d: player %d won", round, *gs.Winner)
						} else {
							log.Printf("round %d: draw", round)
						}
						return true
					}
					moveChan := make(chan rules.Move, 1)
					ticks <- &Tick{
						Ctx:       ctx,
						Gamestate: gs.Preprocess(),
						Move:      moveChan,
					}
					var move rules.Move
					select {
					case m, ok := <-moveChan:
						if ok {
							move = m
						} else {
							log.Printf("round %d: got no move, doing nothing", round)
							move = rules.Nop
						}
					case <-ctx.Done():
						log.Printf("round %d: no respons in time, doing nothing", round)
						move = rules.Nop
					}
					_, err = move.WriteTo(conn)
					if err != nil {
						log.Printf("round %d: error writing response: %s", round, err)
						return true
					}
					return false
				}()
				if stop {
					return
				}
			}
		}()
	}
}
