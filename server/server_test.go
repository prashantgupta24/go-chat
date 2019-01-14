package server_test

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/go-chat/mock_service"
	"github.com/go-chat/server"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestChatServer struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestChatServer))
}

func (suite *TestChatServer) TestRegMessage() {

	// go func() {
	// 	err := http.ListenAndServe(":8001", nil)
	// 	if err != nil {
	// 		log.Fatal("error starting server: ", err)
	// 	}
	// }()

	// chatServer := server.CreateChatServer()
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	conn, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		log.Fatalf("Unable to start websockets : %v ", err)
	// 	}

	// 	go server.Read(conn)
	// 	go server.Write()
	// })

	t := suite.T()
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockConn := mock_service.NewMockConnection(mockCtrl)
	chatServer := server.CreateChatServer()

	//chatServer := serverMock.(*server.MyChatServer)

	message := &server.MessageJSON{
		MsgType: server.REGISTER,
		Sender:  "test",
		Message: "",
	}

	t1 := mockConn.EXPECT().Read().Return(message, nil).Times(1)
	mockConn.EXPECT().Read().After(t1).Return(nil, errors.New("error while parsing")).Times(1)

	mockConn.EXPECT().Conn().Return(nil).Times(1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case msg := <-chatServer.Messages():
			assert.Equal(t, msg.Message, "1")
			assert.Equal(t, msg.MsgType, message.MsgType)
			assert.Equal(t, msg.Sender, message.Sender)
			wg.Done()
		}
	}()

	chatServer.Read(mockConn)
	wg.Wait()
	//go chatServer.Write()

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (suite *TestChatServer) TestPerformance() {
	chatServer := server.CreateChatServer()
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	var conn *websocket.Conn
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("Unable to start websockets : %v ", err)
		}

		connectionStruct := &server.ConnectionStruct{
			conn: conn,
		}
		go chatServer.Read(connectionStruct)
		go chatServer.Write()
	})

	go http.ListenAndServe(":8001", nil)

}
