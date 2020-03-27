package assets

func HomePage() []byte {
	return []byte(`
    <!DOCTYPE html>
    <html lang="en">

    <head>
        <title>THDWB</title>
    </head>

    <body>
        <div style="background-color: red; height: 50px;">
            <h1 style="color: white;">The HotDog WebBrowser!</h1>
        </div>
        <div style="background-color: green; height: 400px;">
            <ul style="margin: 0;">
                <li color: white;>uns</li>
                <li>sao</li>
                <li>bons</li>
                <li>e outros</li>
                <li>nao</li>
            </ul>
        </div>
    </body>

    </html> 
	`)
}
