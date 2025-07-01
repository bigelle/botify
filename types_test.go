package botify

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatBoostSource_Deserialization(t *testing.T) {
	testcases := []struct {
		Name       string
		Input      []byte
		Output     string // e.g. "premium" if should be able to convert into ChatBoostSourcePremium
		ExpectFail bool
	}{
		{
			Name:       "simple check",
			Input:      []byte(`{"source":"premium","user":{}}`),
			Output:     "premium",
			ExpectFail: false,
		},
	}

	for _, tcase := range testcases {
		t.Run(tcase.Name, func(t *testing.T) {
			var result ChatBoostSource
			err := json.Unmarshal(tcase.Input, &result)

			if tcase.ExpectFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			switch tcase.Output {
			case "premium":
				_, err = result.AsPremium()
				if tcase.ExpectFail {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
