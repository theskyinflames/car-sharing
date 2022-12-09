package app

import (
	"context"

	"theskyinflames/car-sharing/internal/domain"

	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

// LocateResponse is a DTO
type LocateResponse struct {
	IsInJourney bool
	Ev          domain.Car
}

// LocateQuery is a query
type LocateQuery struct {
	GroupID int
}

// LocateName is sefl-described
var LocateName = "locate.group"

// Name implements Query interface
func (q LocateQuery) Name() string {
	return LocateName
}

// Locale is a query handler
type Locate struct {
	gr  GroupsRepository
	evr CarsRepository
}

// NewLocate is a constructor
func NewLocate(gr GroupsRepository, evr CarsRepository) Locate {
	return Locate{gr: gr, evr: evr}
}

// Handle implements the QueryHandler interface
func (qh Locate) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	q, ok := query.(LocateQuery)
	if !ok {
		return nil, NewInvalidQueryError(LocateName, query.Name())
	}

	g, err := qh.gr.FindByID(ctx, q.GroupID)
	if err != nil {
		return nil, err
	}

	if !g.IsOnJourney() {
		return LocateResponse{}, nil
	}

	ev, err := qh.evr.FindByID(ctx, g.Ev().ID())
	if err != nil {
		return nil, err
	}

	return LocateResponse{
		IsInJourney: true,
		Ev:          ev,
	}, nil
}
