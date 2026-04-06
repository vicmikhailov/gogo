// Package commonlibs_test showcases how Go developers use the 'testify' library for assertions.
// In Java, you'd use JUnit with AssertJ or Hamcrest. 'testify/assert' is the most common equivalent in Go.
package commonlibs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestifyAssertions(t *testing.T) {
	// 1. Using 'assert' (non-failing/continues execution, like soft assertions in Java)
	// Java comparison: assertThat(actual).isEqualTo(expected);
	name := "Gopher"
	assert.Equal(t, "Gopher", name, "They should be equal")

	// 2. Using 'require' (failing/stops execution, like JUnit's Assert.assertEquals)
	// Java comparison: assertEquals(expected, actual);
	require.NotEmpty(t, name)

	// 3. More complex assertions
	slice := []int{1, 2, 3}
	assert.Contains(t, slice, 2)
	assert.Len(t, slice, 3)

	// Map assertions
	m := map[string]int{"one": 1, "two": 2}
	assert.Contains(t, m, "one")
	assert.NotContains(t, m, "three")

	// 4. Type assertions (similar to instanceof in Java)
	var something interface{} = "I am a string"
	assert.IsType(t, "", something)

	// 5. Error assertions (important in Go)
	// Java comparison: assertThrows(Exception.class, () -> ...)
	var err error = nil
	assert.NoError(t, err)
}

// Example of a table-driven test (common Go idiom) with testify
func TestTableDrivenWithTestify(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive", 10, 20},
		{"zero", 0, 0},
		{"negative", -5, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input * 2
			assert.Equal(t, tt.expected, result)
		})
	}
}
