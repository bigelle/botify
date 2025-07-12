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

func TestMessage_GetCommand(t *testing.T) {
	testcases := []struct {
		Name string
		Input     Message
		Output    string
		ExpectErr bool
	}{
		{
			Name: "1 entity, has text, expect no err",
			Input: Message{
				Text: strPointer("/start"),
				Entities: &[]MessageEntity{
					{
						Type: "bot_command",
						Offset: 0,
						Length: 6,
					},
				},
			},
			Output: "/start",
			ExpectErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			cmd, err := tc.Input.GetCommand()

			if tc.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.Output, cmd)
		})
	}
}

func strPointer(str string) *string {
	return  &str
}
