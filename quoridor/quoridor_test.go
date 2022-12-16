package quoridor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestIds = []uuid.UUID{
	uuid.MustParse("98ae983e-3f04-42ab-928a-c399d6d18375"),
	uuid.MustParse("5341acab-6e28-4d28-8530-8716e0c3dd2e"),
	uuid.MustParse("790bcc3f-6e72-4a0e-a6ea-bc806aa8aa03"),
	uuid.MustParse("6c8420b5-e7f5-4328-ae29-4dbdf7537612"),
	uuid.MustParse("f3282245-f546-4c71-92ca-5bada1f9c037"),
	uuid.MustParse("9cae2aa0-d21a-48ab-a877-4b78942259e4"),
	uuid.MustParse("0ad943b2-6ea9-45ad-9098-f67714652fcd"),
	uuid.MustParse("93ded37f-57d3-4b43-8933-1164e086a881"),
	uuid.MustParse("5b399bd3-aa3e-4754-bb51-175b30b77400"),
	uuid.MustParse("f7ea9019-033b-41e7-a671-26231952cd8c"),
}

func Test_NewGame(t *testing.T) {
	var testCases = []struct {
		id             uuid.UUID
		gameName       string
		expectGame     bool
		expectedErrMsg string
	}{
		{
			TestIds[0], "game 1", true, "",
		},
		{
			uuid.Nil, "game name", false, "unable to create game, need valid id",
		},
		{
			TestIds[0], "", false, "unable to create game, need non-empty nam