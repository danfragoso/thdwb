package assets

func HomePage() []byte {
	return []byte(`
    <!DOCTYPE html>
    <html lang="en">

    <head>
        <title>THDWB Home page test</title>
    </head>

    <body>
       <h3>
            This is the hotdog web browser, it's a (toy) web rendering engine and web browser written from scratch in golang.
       </h3>
       <p>
            I started this project just to see if it was really possible for me to build a big project like that and to learn how the web works under the hood.
       </p>
       <p>
            If you try to access a website it may not really load correctly, the reason is because the layout implementation is pretty lacking at the moment, I got a basic implemetation of a display block and list item. For now there's no flexbox and no inline-block.
       </p>
        <p>
            There's no image support, no forms and no links. But I'm slowly working on bringing those things over. One thing that somewhat works is inline css via the style tag. For now something like this will work:
       </p>
       <div style="font-size: 30px; color: aqua; background-color: purple;" >
        style="font-size: 30px; color: aqua; background-color: purple;"
       </div>
        <p>
            If you want an example page to load, there are some tests at the test folder on the root of the project.
       </p>
        <p>
          https://motherfuckingwebsite.com is a site that works ok.
       </p>
        <p>
            Below is the list of components:
       </p>
       <ul>
            <li>ketchup (html parser and DOM Tree builder)</li>
            <li>mayo (css parser and Render Tree builder)</li>
            <li>mustard (UI Toolkit, events and OpenGL)</li>
            <li>sauce (requests, cache and filesystem)</li>
            <li>bun (css layout calculator)</li>
            <li>gg (drawing routines and text rendering)</li>
       </ul>
       <p>
            You can find more info at
       </p>
       <a style="color: blue;" href="https://github.com/danfragoso/thdwb">https://github.com/danfragoso/thdwb</a>
    </body>

    </html> 
	`)
}
