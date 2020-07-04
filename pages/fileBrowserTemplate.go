package pages

func fileBrowserTemplate() string {
	var template string
	template = `
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>File Browser</title>
	</head>
	<body>
		<div>
			<h1>Index of {{.Path}}</h1>
		</div>
		{{range .Dirs}}
			<div>
				<a href="file://{{$.Path}}/{{.}}/" style="font-size: 16px;">{{.}}/</a>
			</div>
		{{end}}
		{{range .Files}}
			<div>
				<a href="file://{{$.Path}}/{{.}}" style="font-size: 16px;">{{.}}</a>
			</div>
		{{end}}
	</html>
	`
	return template
}
