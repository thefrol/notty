package app

import "errors"

var (
	ErrorSubscriptionNotFound = errors.New("subscription not found")
	ErrorSubscriptionExists   = errors.New("subscription exists")

	ErrorCustomerNotFound = errors.New("client not found")
	ErrorCustomerExists   = errors.New("client exists")
)
