![thdwb](https://raw.githubusercontent.com/danfragoso/thdwb/master/thdwb.png))
___

This is the hotdog web browser, it's a web browser written from scratch entirely in golang.

<p align="center">
  <img src="screenshot.png?raw=true"></img>
</p>

### Components
- ketchup (html parser and DOM Tree builder)
- mayo (css parser and Render Tree builder)
- mustard (UI Toolkit, events and OpenGL)
- sauce (requests, cache and filesystem)
- bun (css layout calculator)
- [gg](https://github.com/fogleman/gg) (drawing routines and text rendering)

### Getting started
- Running 

  ```sh
  make
  ```
  This command will start the browser and load an example page
  
- Testing

  ```sh
  make test
  ```
  This command will run all the configurated unit tests
  
- Building

  ```sh
  make build
  ```
  This command will build the binary version
