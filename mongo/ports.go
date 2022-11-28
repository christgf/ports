package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/christgf/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// port is the representation of ports.Port as a BSON document.
type port struct {
	ID       string    `bson:"id"`
	Name     string    `bson:"name"`
	Code     string    `bson:"code"`
	City     string    `bson:"city"`
	Province string    `bson:"province"`
	Country  string    `bson:"country"`
	Alias    []string  `bson:"alias"`
	Regions  []string  `bson:"regions"`
	Timezone string    `bson:"timezone"`
	UNLocs   []string  `bson:"UNLocs"`
	Coords   []float64 `bson:"coords"`
}

// InsertPort will insert a new BSON document in the Ports collection, based on
// the information provided. If a document already exists with the same port.ID,
// then the existing BSON document is replaced with a new one, even if the two
// are exactly the same.
func (db *DB) InsertPort(ctx context.Context, p ports.Port) error {
	if _, err := db.Ports().ReplaceOne(ctx, bson.D{{Key: "id", Value: p.ID}}, port{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		City:     p.City,
		Province: p.Province,
		Country:  p.Country,
		Alias:    p.Alias,
		Regions:  p.Regions,
		Timezone: p.Timezone,
		UNLocs:   p.UNLocs,
		Coords:   p.Coords,
	}, options.Replace().SetUpsert(true)); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

// FindPort will attempt to retrieve a single BSON document from the Ports
// collection, based on the identifier provided, and return the corresponding
// information as ports.Port. It returns an error if a port document with the
// provided identifier could not be found.
func (db *DB) FindPort(ctx context.Context, portID string) (*ports.Port, error) {
	res := db.Ports().FindOne(ctx, bson.D{{Key: "id", Value: portID}})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &ports.Error{Code: ports.ErrCodeNotFound, Msg: "port not found"}
		}

		return nil, fmt.Errorf("find: %w", err)
	}

	p := new(port)
	if err := res.Decode(p); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &ports.Port{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		City:     p.City,
		Province: p.Province,
		Country:  p.Country,
		Alias:    p.Alias,
		Regions:  p.Regions,
		Timezone: p.Timezone,
		UNLocs:   p.UNLocs,
		Coords:   p.Coords,
	}, nil
}
