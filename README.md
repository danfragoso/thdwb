<p align="center">
  <img src="thdwb.png?raw=true"></img>
</p>

This is the hotdog web browser, it's a web browser written from scratch entirely in golang.

### Components
- ketchup (html parser, DOM and Render Tree builder)
- mayo (css parser, stylesheet and layout calculator)
- mustard (browser UI, window and renderer)
- sauce (requests, cache and filesystem)

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
  This command will build the binary and move it to the bin/ folder

### Screenshot
<p align="center">
  <img src="screenshot.png?raw=true"></img>
</p>
