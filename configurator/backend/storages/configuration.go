package storages

import (
	"errors"
	"time"
)

var ErrConfigurationNotFound = errors.New("Configuration wasn't found")

//ConfigurationsStorage - Collection here is used as a type of configuration - like destinations, api_keys, custom_domains, etc.
type ConfigurationsStorage interface {
	//Get returns a single configuration from collection
	//If configuration is not found, must return ErrConfigurationNotFound for correct response message
	Get(collection string, id string) ([]byte, error)
	//GetAllGroupedByID returns all the configurations of requested type grouped by id (result must be
	//deserializable to map[string]<entity_type>
	GetAllGroupedByID(collection string) (map[string][]byte, error)
	//GetCollectionLastUpdated returns time when collection was last updated
	//(max _lastUpdated field among entities)
	GetCollectionLastUpdated(collection string) (*time.Time, error)
	//UpdateCollectionLastUpdated updates time when collection was last updated
	UpdateCollectionLastUpdated(collection string) error
	//Store saves entity and also must update _lastUpdated field of the collection
	Store(collection string, id string, entity []byte) error
	//Delete deletes entity and also must update _lastUpdated field of the collection
	Delete(collection string, id string) error

	// GetRelationIndex returns all key IDs from the relation.
	GetRelationIndex(relation string) ([]string, error)
	// DeleteRelation removes all related IDs for the relation.
	DeleteRelation(relation, id string) error
	// GetRelatedIDs returns a slice of related IDs from the relation.
	GetRelatedIDs(relation, id string) ([]string, error)
	// AddRelatedIDs adds related IDs to the relation.
	AddRelatedIDs(relation, id string, relatedIDs ...string) error
	// DeleteRelatedIDs deletes related IDs from the relation.
	DeleteRelatedIDs(relation, id string, relatedIDs ...string) error

	//Close frees all the resources used by the storage (close connections etc.)
	Close() error
}
