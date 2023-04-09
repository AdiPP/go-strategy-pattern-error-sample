package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnknown          = errors.New("unknown error")
)

func main() {
	err1 := fmt.Errorf("some error: %w", ErrNotFound)
	err2 := fmt.Errorf("some error: %w", ErrPermissionDenied)
	err3 := fmt.Errorf("some error: %w", ErrUnknown)

	resolver := NewErrorResolverStartegy()

	fmt.Println(resolver.ResolveError(err1))
	fmt.Println(resolver.ResolveError(err2))
	fmt.Println(resolver.ResolveError(err3))
}

type ErrorResolver interface {
	IsError(error) bool
	ResolveError(error) error
}

type NotFoundResolver struct{}

func (r NotFoundResolver) IsError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func (r NotFoundResolver) ResolveError(err error) error {
	return fmt.Errorf("resource not found: %w", err)
}

type PermissionDeniedResolver struct{}

func (r PermissionDeniedResolver) IsError(err error) bool {
	return errors.Is(err, ErrPermissionDenied)
}

func (r PermissionDeniedResolver) ResolveError(err error) error {
	return fmt.Errorf("permission denied: %w", err)
}

type OtherErrorResolver struct{}

func (r OtherErrorResolver) IsError(err error) bool {
	return true
}

func (r OtherErrorResolver) ResolveError(err error) error {
	return fmt.Errorf("unexpected error occured: %w", err)
}

type ErrorResolverStartegy struct {
	resolvers []ErrorResolver
}

func NewErrorResolverStartegy() ErrorResolverStartegy {
	return ErrorResolverStartegy{
		resolvers: []ErrorResolver{
			NotFoundResolver{},
			PermissionDeniedResolver{},
			OtherErrorResolver{},
		},
	}
}

func (s ErrorResolverStartegy) ResolveError(err error) error {
	for _, r := range s.resolvers {
		if r.IsError(err) {
			return r.ResolveError(err)
		}
	}

	return err
}
