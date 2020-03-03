package sauce

func loadErrorPage(err string) string {
	return `
	<html>
		<head>
			<title>
				Error!
			</title>
		</head>
		<body>
			<h1>` + err + `</h1>
		</body>
	</html>`
}
