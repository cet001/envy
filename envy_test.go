package envy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const stringKey = "STRKEY"
const intKey = "INTKEY"
const nonexistentKey = "NONEXISTENT_KEY!"
const anotherNonexistentKey = "ANOTHER_NONEXISTENT_KEY!"

func ExampleEnv() {
	os.Setenv("HOST", "example.com")
	os.Setenv("DBNAME", "customerdb")
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("DBNAME")
	}()

	env := NewEnv()
	fmt.Println(env.Getenv("HOST").String())
	fmt.Println(env.Getenv("DBNAME").String())
	fmt.Println(env.Getenv("PORT").DefaultValue("8080").Int()) // PORT wasn't set, so use default

	// Output:
	// example.com
	// customerdb
	// 8080
}

func TestGetenv_String(t *testing.T) {
	testSetup(func(env *Env) {
		assert.Equal(t, "ABC", env.Getenv(stringKey).String())
		assert.Equal(t, "ABC", env.Getenv(stringKey).DefaultValue("XYZ").String())
		assert.Equal(t, "XYZ", env.Getenv(nonexistentKey).DefaultValue("XYZ").String())
		assert.Equal(t, 0, len(env.Errors()))

		assert.Equal(t, "", env.Getenv(nonexistentKey).String())
		assert.Equal(t, 1, len(env.Errors()))
	})
}

func TestGetenv_Int(t *testing.T) {
	testSetup(func(env *Env) {
		assert.Equal(t, 123, env.Getenv(intKey).Int())
		assert.Equal(t, 123, env.Getenv(intKey).DefaultValue("456").Int())
		assert.Equal(t, 456, env.Getenv(nonexistentKey).DefaultValue("456").Int())
		assert.Equal(t, 0, len(env.Errors()))

		assert.Equal(t, 0, env.Getenv(nonexistentKey).Int())
		assert.Equal(t, 0, env.Getenv(anotherNonexistentKey).Int())
		assert.Equal(t, 2, len(env.Errors()))
	})
}

func TestUnparsableIntVar(t *testing.T) {
	testSetup(func(env *Env) {
		// Both default ints values should cause a int parsing error
		assert.Equal(t, 0, env.Getenv(nonexistentKey).DefaultValue("XXXXX").Int())
		assert.Equal(t, 0, env.Getenv(nonexistentKey).DefaultValue("12.34").Int())
		assert.Equal(t, 2, len(env.Errors()))
	})
}

func testSetup(runTest func(env *Env)) {
	os.Setenv(stringKey, "ABC")
	os.Setenv(intKey, "123")
	defer func() {
		os.Unsetenv(stringKey)
		os.Unsetenv(intKey)
	}()

	runTest(NewEnv())
}
