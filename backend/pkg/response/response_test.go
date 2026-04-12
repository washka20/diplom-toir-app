package response_test

import (
	"encoding/json"
	"testing"

	"toir-app/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccess_MarshalJSON(t *testing.T) {
	data := map[string]string{"name": "test"}
	resp := response.Success(data)

	b, err := json.Marshal(resp)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &result))

	assert.Equal(t, true, result["success"])
	assert.NotNil(t, result["data"])
	assert.Nil(t, result["error"])
	assert.Nil(t, result["meta"])
}

func TestError_MarshalJSON(t *testing.T) {
	resp := response.Error("something went wrong")

	b, err := json.Marshal(resp)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &result))

	assert.Equal(t, false, result["success"])
	assert.Nil(t, result["data"])
	assert.Equal(t, "something went wrong", result["error"])
	assert.Nil(t, result["meta"])
}

func TestPaginated_IncludesMeta(t *testing.T) {
	data := []string{"a", "b", "c"}
	resp := response.Paginated(data, 1, 10, 100)

	b, err := json.Marshal(resp)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &result))

	assert.Equal(t, true, result["success"])
	assert.NotNil(t, result["data"])
	assert.Nil(t, result["error"])

	meta, ok := result["meta"].(map[string]interface{})
	require.True(t, ok, "meta should be an object")
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(10), meta["per_page"])
	assert.Equal(t, float64(100), meta["total"])
}
