package queue

import (
	"github.com/brutalzinn/ho-routine-obrigations/models"
)

/// a lazy and bugged queue

// var ObrigationPending *obrigation.Obrigation
var ObrigationsQueue chan models.ObrigationQueuePending

func StartQueue() {
	ObrigationsQueue = make(chan models.ObrigationQueuePending)
}
