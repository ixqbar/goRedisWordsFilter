package filter

import (
	"errors"
	"github.com/jonnywang/go-kits/redis"
	"os"
	"os/signal"
	"syscall"
)

type WordsFilterRedisHandler struct {
	redis.RedisHandler
	wordsMapping *WordsMapping
}

func (obj *WordsFilterRedisHandler) Init() error {
	obj.Initiation(func() {
		obj.wordsMapping = NewWordsMapping()
		obj.wordsMapping.Load()
	})

	return nil
}

func (obj *WordsFilterRedisHandler) Shutdown() {
	Logger.Print("redis server will shutdown")
}

func (obj *WordsFilterRedisHandler) Version() (string, error) {
	return VERSION, nil
}

func (obj *WordsFilterRedisHandler) Ping(message string) (string, error) {
	if len(message) > 0 {
		return message, nil
	}

	return "PONG", nil
}

func (obj *WordsFilterRedisHandler) Filter(sentence string) (string, error) {
	if len(sentence) == 0 {
		return sentence, nil
	}

	return string(obj.wordsMapping.FilterSentence([]rune(sentence))), nil
}

func (obj *WordsFilterRedisHandler) Exists(word string) (int, error) {
	if  len(word) > 0 && obj.wordsMapping.ExistsWord(word) {
		return 1, nil
	}

	return 0, nil
}

func (obj *WordsFilterRedisHandler) Add(word string) error {
	if len(word) == 0 {
		return errors.New("error params")
	}

	if obj.wordsMapping.AddNewWord([]rune(word), true) {
		return nil
	}

	return errors.New("FAIL")
}

func (obj *WordsFilterRedisHandler) Delete(word string) error {
	if len(word) == 0 {
		return errors.New("error params")
	}

	if obj.wordsMapping.DeleteWord([]rune(word)) {
		return nil
	}

	return errors.New("FAIL")
}

func (obj *WordsFilterRedisHandler) FlushAll() error {
	obj.wordsMapping.Clear()

	return nil
}

func (obj *WordsFilterRedisHandler) Total() (int, error) {
	return obj.wordsMapping.TotalWords(), nil
}

func (obj *WordsFilterRedisHandler) Reload() error {
	obj.wordsMapping.Clear()
	obj.wordsMapping.Load()

	return nil
}

func Run() {
	wordsFilterHandler := &WordsFilterRedisHandler{}

	err := wordsFilterHandler.Init()
	if err != nil {
		Logger.Print(err)
		return
	}

	wordsFilterServer, err := redis.NewServer(GConfig.ListenServer, wordsFilterHandler)
	if err != nil {
		Logger.Print(err)
		return
	}

	serverStop := make(chan bool)
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-stopSignal
		wordsFilterServer.Stop(10)
		serverStop <- true
	}()

	err = wordsFilterServer.Start()
	if err != nil {
		Logger.Print(err)
		stopSignal <- syscall.SIGTERM
	}

	<-serverStop

	close(serverStop)
	close(stopSignal)

	Logger.Print("all server shutdown")
}
