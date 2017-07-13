package kcache

import (
	"context"

	logutil "github.com/boz/go-logutil"
)

type FilterSubscription interface {
	Subscription
	Refilter(Filter)
}

type filterSubscription struct {
	parent Subscription

	refilterch chan Filter

	outch   chan Event
	readych chan struct{}
	stopch  chan struct{}
	cache   cache

	log logutil.Log
}

func NewFilterSubscription(log logutil.Log, parent Subscription, filter Filter) FilterSubscription {

	ctx := context.Background()

	stopch := make(chan struct{})

	s := &filterSubscription{
		parent:     parent,
		refilterch: make(chan Filter),
		outch:      make(chan Event, EventBufsiz),
		readych:    make(chan struct{}),
		stopch:     stopch,
		cache:      newCache(ctx, log, stopch, filter),
		log:        log,
	}

	go s.run()

	return s
}

func (s *filterSubscription) Cache() CacheReader {
	return s.cache
}
func (s *filterSubscription) Ready() <-chan struct{} {
	return s.readych
}
func (s *filterSubscription) Events() <-chan Event {
	return s.outch
}
func (s *filterSubscription) Close() {
	s.parent.Close()
}
func (s *filterSubscription) Done() <-chan struct{} {
	return s.parent.Done()
}

func (s *filterSubscription) Refilter(filter Filter) {
	select {
	case s.refilterch <- filter:
	case <-s.Done():
	}
}

func (s *filterSubscription) run() {
	defer s.log.Un(s.log.Trace("run"))
	defer close(s.outch)
	defer close(s.stopch)

	preadych := s.parent.Ready()

	for {
		select {
		case <-preadych:

			list, err := s.parent.Cache().List()

			if err != nil {
				s.log.Err(err, "parent.Cache().List()")
				s.parent.Close()
			} else {
				s.cache.sync(list)
				close(s.readych)
			}

			preadych = nil

		case filter := <-s.refilterch:

			list, err := s.parent.Cache().List()
			if err != nil {
				s.log.Err(err, "parent.Cache().List()")
				s.parent.Close()
				continue
			}

			s.distributeEvents(s.cache.refilter(list, filter))

		case evt, ok := <-s.parent.Events():
			if !ok {
				return
			}

			s.distributeEvents(s.cache.update(evt))
		}
	}
}

func (s *filterSubscription) distributeEvents(events []Event) {
	for _, evt := range events {
		select {
		case s.outch <- evt:
		default:
			s.log.Warnf("event buffer overrun")
		}
	}
}
