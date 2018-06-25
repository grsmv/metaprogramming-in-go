package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"plugin"
	"text/template"
	"time"
)

const pluginName = "plugin"

var code = `package main

func Add(x int) int {
  return x + {{.Number}}
}
`

func generateCode() string {
	var (
		tmpl, _ = template.New("t").Parse(code)
		buf     bytes.Buffer
	)
	tmpl.Execute(&buf, struct {
		Number int
	}{
		rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100),
	})
	return buf.String()
}

func writeFile(name, code string) {
	ioutil.WriteFile(name+".go", []byte(code), 0644)
}

func compilePlugin(name string) {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", name+".so", name+".go")
	cmd.Run()
}

func loadAndExecutePlugin(name string) {
	p, err := plugin.Open("./" + name + ".so")
	if err != nil {
		panic(err)
	}

	add, err := p.Lookup("Add")
	if err != nil {
		panic(err)
	}

	sum := add.(func(int) int)(100)
	println(".Add(int) result:", sum)
}

func cleanUp(name string) {
	os.Remove(name + ".go")
	os.Remove(name + ".so")
}

func main() {

	// Generating Go code from template
	var source = generateCode()
	writeFile(pluginName, source)

	// Compiling loadable .so plugin
	compilePlugin(pluginName)

	// Loading and executing .so plugin
	loadAndExecutePlugin(pluginName)

	// Tear down
	time.Sleep(1 * time.Second)
	cleanUp(pluginName)
}
