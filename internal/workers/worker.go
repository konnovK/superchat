package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/repository"
	"gorm.io/gorm"
)

type Worker struct {
	WorkerEnqueuer work.Enqueuer
	WorkerPool     work.WorkerPool
	chatRepository repository.Chat
}

const (
	namespace = "chats_ttl_check_worker_namespace"
	deleteJob = "delete_chat_job"
)

func NewWorker(redisPool *redis.Pool) *Worker {
	worker := &Worker{
		WorkerEnqueuer: *work.NewEnqueuer(namespace, redisPool),
		WorkerPool:     *work.NewWorkerPool(Worker{}, 3, namespace, redisPool),
	}

	worker.WorkerPool.Job(deleteJob, (*Worker).deleteJob)

	return worker
}

func (w *Worker) deleteJob(job *work.Job) error {
	chatId := job.ArgInt64("chatId")
	return w.chatRepository.Delete(&entity.Chat{
		Model: gorm.Model{
			ID: uint(chatId),
		},
	})
}

func (w *Worker) SetJobTimer(delay int64, chatId int) error {
	_, err := w.WorkerEnqueuer.EnqueueIn(deleteJob, delay, work.Q{"chatId": chatId})
	return err
}
