package lockitdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurn(t *testing.T) {

	testcases := []struct {
		direction TurnDirection
		expected  Pair
	}{
		{Left,
