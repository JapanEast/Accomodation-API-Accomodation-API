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
	uuid.Must