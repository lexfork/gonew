{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import .Imports "//log" "//fmt" "//os" }}

var opt = parseFlags()

{{ func "main" }}