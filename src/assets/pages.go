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
    I started this project just to see if it was really possible for me to build a big project like that and to learn
    how the web works under the hood.
  </p>
  <p>
    If you try to access a website it may not really load correctly, the reason is because the layout implementation is
    pretty lacking at the moment, I got a basic implemetation of a display block, list item and inline. For now there's
    no
    flexbox and no inline-block.
  </p>
  <p>
    There's no image support, no forms and no links. But I'm slowly working on bringing those things over. One thing
    that somewhat works is inline css via the style tag. For now something like this will work:
  </p>
  <div>
    <div style="color: white; background-color: coral;">
      Inline
    </div>
    <div style="color: white; background-color: cornflowerblue;">
      style
    </div>
    <div style="color: white; background-color: darkgreen;">
      with
    </div>
    <div style="color: white; background-color: purple;">
      block
    </div>
    <div style="color: black; background-color: yellow;">
      elements
    </div>
  </div>
  <p>
    If you want an example page to load, there are some tests at the test folder on the root of the project, there are
    also some links below.
  </p>
  <ul>
    <li><a href="https://motherfuckingwebsite.com">https://motherfuckingwebsite.com</a></li>
    <li><a href="http://lite.cnn.com/">http://lite.cnn.com/</a></li>
    <li><a href="http://serenityos.org/">http://serenityos.org/</a></li>
    <li><a href="thdwb://about/">About</a></li>
  </ul>
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
    You can find more info at my github
  </p>
  <a href="https://github.com/danfragoso/thdwb">https://github.com/danfragoso/thdwb</a>
</body>

</html>
`)
}

func DefaultPage() []byte {
	return []byte(`
    <html>
      <head>
        <title>THDWB</title>
      </head>
      <body>
        <div>
          <h3>Sorry, this page does not exist.</h3>
          <div>THDWB; The Hotdog web browser.</div>
          <ul>
            <li><a href="thdwb://history/">History</a></li>
            <li><a href="thdwb://about/">About</a></li>
          </ul>
          <div> - - - - - - - - -</div>
          <a href="thdwb://homepage/">Go back to home</a>
        </div>
      </body>
    </html>
  `)
}
