package mock

import "github.com/boson-project/grid"

type Adapter struct{}

func (n Adapter) Instances() (int, error)                       { return 0, nil }
func (n Adapter) SubscriptionManager() grid.SubscriptionManager { return SubscriptionManager{} }
func (n Adapter) EventManager() grid.EventManager               { return EventManager{} }

type SubscriptionManager struct{}

func (n SubscriptionManager) Create(string) error     { return nil }
func (n SubscriptionManager) Delete(string) error     { return nil }
func (n SubscriptionManager) List() ([]string, error) { return []string{}, nil }

type EventManager struct{}

func (n EventManager) Create(string) error     { return nil }
func (n EventManager) Delete(string) error     { return nil }
func (n EventManager) List() ([]string, error) { return []string{}, nil }
