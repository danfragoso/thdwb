package pages

import (
	structs "thdwb/structs"
)

func RenderAboutPage(buildInfo *structs.BuildInfo) string {
	var template string
	template = `
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>About</title>
	</head>
	<body>
		<img src="https://raw.githubusercontent.com/danfragoso/thdwb/master/assets/thdwb.png">
		<div>
			<ul>
				<li>REVISION: ` + buildInfo.GitRevision + `</li>
				<li>BRANCH: ` + buildInfo.GitBranch + `</li>
				<li>HOST: ` + buildInfo.HostInfo + `</li>
				<li>BUILD_TIME: ` + buildInfo.BuildTime + `</li>
			</ul>
		</div>
		<a href="thdwb://homepage/">Go back to home</a>
	</html>
	`
	return template
}
