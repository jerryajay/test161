package test161

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildFull(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	conf := &BuildConf{
		Repo:     "https://github.com/ops-class/os161.git",
		CommitID: "HEAD",
		KConfig:  "DUMBVM",
	}

	env := defaultEnv.CopyEnvironment()
	env.RootDir = "./fixtures/root"

	test, err := conf.ToBuildTest(env)
	assert.Nil(err)
	assert.NotNil(test)

	if test == nil {
		t.Log(err)
		t.FailNow()
	}

	_, err = test.Run(env)
	assert.Nil(err)

	t.Log(test.OutputJSON())

	for k, v := range env.keyMap {
		t.Log(k, v)
	}

}

type confDetail struct {
	repo      string
	commit    string
	config    string
	reqCommit string
}

func TestBuildFailures(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	configs := []confDetail{
		confDetail{"https://notgithub.com/ops-class/os161111.git", "HEAD", "DUMBVM", ""},
		confDetail{"https://github.com/ops-class/os161.git", "aaaaaaaaaaa111111112222", "FOO", ""},
		confDetail{"https://github.com/ops-class/os161.git", "HEAD", "FOO", ""},
		confDetail{"https://github.com/ops-class/os161.git", "HEAD", "DUMBVM", "notavalidcommitit"},
	}

	for _, c := range configs {

		conf := &BuildConf{
			Repo:           c.repo,
			CommitID:       c.commit,
			KConfig:        c.config,
			RequiredCommit: c.reqCommit,
		}

		test, err := conf.ToBuildTest(defaultEnv)
		assert.NotNil(test)

		res, err := test.Run(defaultEnv)
		assert.NotNil(err)
		assert.Nil(res)
	}
}

func TestHexString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testCases := []struct {
		input    string
		expected bool
	}{
		{"0123456789abcdef", true},
		{"0123456789ABCDEF", true},
		{"e1a2fbd038c618b6d9e636a94e1907dc92e94ca6", true},
		{"", false},
		{"foo", false},
		{"0123456789abcdefg", false},
		{"!", false},
	}

	for _, test := range testCases {
		t.Log(test.input)
		assert.Equal(test.expected, isHexString(test.input))
	}
}
