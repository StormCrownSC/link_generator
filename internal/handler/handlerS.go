package handler

import (
	"context"
	"log"
	"server/pg"
	"strconv"
	"sync"
	"utils"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

func main() {
	// Инициализация
	conn := pg.СonnectToDB()
	pg.СreateTableIfNotExists(conn)
	pg.СreateRobots(conn)

	// Rsocket
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := rsocket.Receive().
			OnStart(func() {
				log.Println("Server Started")
			}).
			Acceptor(func(_ context.Context, _ payload.SetupPayload, _ rsocket.CloseableRSocket) (rsocket.RSocket, error) {
				return rsocket.NewAbstractSocket(
					////////////////////////////////////////////////////////////////////////
					// Request-Response
					rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
						// Получить ID - вернуть робота с таким ID
						id, err := strconv.Atoi(msg.DataUTF8())
						if err != nil {
							log.Fatalln(err)
							return nil
						}
						utils.ColorPrintln(utils.ColorBlue, "ReqRes - Получение робота с ID:", id)
						robot := pg.RobotByID(conn, id)
						return mono.Just(payload.NewString(robot.String(), ""))
					}),

					// Request-Stream
					rsocket.RequestStream(func(request payload.Payload) flux.Flux {
						// Получить модель...
						model := request.DataUTF8()
						utils.ColorPrintln(utils.ColorBlue, "ReqStr - Получение роботов с моделью:", model)
						robots := make(chan pg.Robot)
						go pg.RobotsByModel(conn, model, robots)
						// ...Вернуть роботов с такой моделью
						return flux.Create(func(_ context.Context, s flux.Sink) {
							for robot := range robots {
								s.Next(payload.NewString(robot.String(), ""))
							}
							s.Complete()
						})
					}),

					// Request-Channel
					rsocket.RequestChannel(func(payloads flux.Flux) flux.Flux {
						robots := make(chan pg.Robot)
						ids := make(chan int)
						var fakeColor string
						// Получить множество ID - вернуть множество роботов с соседними ID
						// и непонастоящему перекрашивать их по запросу
						payloads.
							DoOnComplete(func() {
								close(ids)
							}).
							DoOnNext(func(msg payload.Payload) error {
								//if slices.Contains(pg.Colors, msg.DataUTF8()) {
								// Перекраска роботов в черный
								if msg.DataUTF8() == "Black" {
									fakeColor = msg.DataUTF8()
									utils.ColorPrintln(utils.ColorBlue, "ReqCh - Смена цвета на:", fakeColor)
								} else {
									id, err := strconv.Atoi(msg.DataUTF8())
									if err != nil {
										log.Fatalln(err)
										return nil
									}
									ids <- id
								}
								return nil
							}).Subscribe(context.Background())
						// Функция смены цвета
						fakeUpdateRobotColor := func(robot pg.Robot) pg.Robot {
							robot.Color = fakeColor
							return robot
						}
						// Получаем роботов
						go func() {
							for id := range ids {
								utils.ColorPrintln(utils.ColorBlue, "ReqCh - Получение робота с ID:", id)
								if fakeColor != "" {
									robots <- fakeUpdateRobotColor(pg.RobotByID(conn, id-1))
									robots <- fakeUpdateRobotColor(pg.RobotByID(conn, id+1))
								} else {
									robots <- pg.RobotByID(conn, id-1)
									robots <- pg.RobotByID(conn, id+1)
								}
							}
							close(robots)
						}()
						// Возвращаем роботов
						return flux.Create(func(_ context.Context, s flux.Sink) {
							for robot := range robots {
								s.Next(payload.NewString(robot.String(), ""))
							}
							s.Complete()
						})
					}),

					// Fire-and-forget
					rsocket.FireAndForget(func(msg payload.Payload) {
						// Перекрасить всех роботов в указанный цвет
						utils.ColorPrintln(utils.ColorBlue, "FaF - Перекраска роботов на", msg.DataUTF8())
						pg.ChangeAllRobotsColor(conn, msg.DataUTF8())
					}),

					////////////////////////////////////////////////////////////////////////
				), nil
			}).
			Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
			Serve(ctx)

		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()
	utils.EnterExit()
	cancel()
	wg.Wait()
}
