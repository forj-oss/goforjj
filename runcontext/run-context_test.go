package runcontext_test

import (
	"goforjj/runcontext"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type containerMock struct {
	envs    map[string]rune
	opts    []string
	volumes map[string]rune
}

func newContainerMock() (ret *containerMock) {
	ret = new(containerMock)
	ret.opts = make([]string, 0)
	ret.envs = make(map[string]rune)
	ret.volumes = make(map[string]rune)
	return
}

func (d *containerMock) addOpts(opts ...string) {
	d.opts = append(d.opts, opts...)
}

func (d *containerMock) addEnv(key, value string) {
	env := key + "=" + value
	d.envs[env] = 'e'
}

func (d *containerMock) AddVolume(volume string) {
	d.volumes[volume] = 'v'
}

func TestRunContext(t *testing.T) {
	assert := assert.New(t)

	// ---------------------------------------
	testCase := "when initialize run Context"
	runC := runcontext.NewRunContext("TEST", 3)

	assert.NotNilf(runC, "Must be non nil %s", testCase)

	// ---------------------------------------
	testCase = "when options are added"
	runC.AddEnv("key1", "value1")
	runC.AddVolume("volume1")
	runC.AddOptions("Opt1")

	result := runC.BuildOptions()

	assert.ElementsMatchf(result, []string{"-e", "key1=value1", "-v", "volume1", "Opt1", "-e", "TEST=-v 'volume1' 'Opt1' -e 'key1=value1'"}, "Expect to find all added values %s", testCase)

	// ---------------------------------------
	testCase = "when container is set"
	container := newContainerMock()

	runC = runcontext.NewRunContext("TEST", 3)
	runC.DefineContainerFuncs(container.AddVolume, container.addEnv, container.addOpts)
	runC.AddEnv("key1", "value1")
	runC.AddVolume("volume1")
	runC.AddOptions("-h")

	assert.Containsf(container.envs, "key1=value1", "Expect key1=value to be found %s", testCase)
	assert.Containsf(container.volumes, "volume1", "Expect volume1 to be found %s", testCase)
	assert.Containsf(container.opts, "-h", "Expect -h to be found %s", testCase)

	// ---------------------------------------
	container = newContainerMock()
	runC = runcontext.NewRunContext("TEST", 3)
	os.Setenv("TEST", "-e key2=value1 -v volume2 -x")
	runC = runcontext.NewRunContext("TEST", 3)
	runC.DefineContainerFuncs(container.AddVolume, container.addEnv, container.addOpts)
	ret := runC.GetFrom()

	assert.Truef(ret, "Expect to get values from TEST")
	assert.Containsf(container.envs, "key2=value1", "Expect key2=value1 to be found %s", testCase)
	assert.Containsf(container.volumes, "volume2", "Expect volume1 to be found %s", testCase)
	assert.Containsf(container.opts, "-x", "Expect -x to be found %s", testCase)

}
