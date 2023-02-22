package main_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/tests/utils"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SaveEcosystem struct {
	SaveEcosystem *model.LoginResult `json:"saveEcosystem"`
}

func TestSaveEcosystemCreatedEntity(t *testing.T) {
	query := "mutation SaveEcosystem($ecosystem: EcosystemInput!) { saveEcosystem(ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"ecosystem\": { \"name\": \"tEst eCosystem\", \"aquarium\": { \"dimensions\": { \"width\": 10 } } } }"

	var data graphql.Response
	state, token := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &token)
	user, err := state.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems, 1)
	assert.Equal(t, "tEst eCosystem", user.Ecosystems[0].Name)
	assert.Equal(t, 10, user.Ecosystems[0].Aquarium.Dimensions.Width)
	assert.Equal(t, 0, user.Ecosystems[0].Aquarium.Dimensions.Height)
}
