import subprocess

program = '''
import (
	"fmt"
)

func main() {
    fmt.Println("PYTHON!")
}
'''

with open("./template.txt", "r") as f:
    contents = f.read()

begin = contents.split("// <-- begin run -->")[0]
end = contents.split("// <-- end run -->")[1]
injected =  begin  + program.replace("func main() {", "func run() {") + end

with open("hosted.go", "w") as f:
    f.write(injected)

subprocess.run(["go", "fmt", "."])