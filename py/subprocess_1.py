# func openOSEditor(os, notesFile string) * exec.Cmd {
# 	if os == "windows" {
# 		return exec.Command("notepad", notesFile)
# 	}
# 	return exec.Command("/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", "+$", notesFile)
# }

import subprocess
# subprocess.call("command1")
# subprocess.call(["/usr/bin/xterm -fa", "arg1", "arg2"])
subprocess.call(
    ["/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", '+$', '/tmp/x.txt'])


# /usr/bin/xterm -fa Monospace -fs 14 -e /usr/bin/vim +$ '/tmp/x.txt'
