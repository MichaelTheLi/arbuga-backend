package main_test

import (
	"arbuga/backend/api/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SaveEcosystem struct {
	SaveEcosystem *model.LoginResult `json:"saveEcosystem"`
}

func TestSaveEcosystemCreatedEntity(t *testing.T) {
	assert.Fail(t, "Not implemented")
}
