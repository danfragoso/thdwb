package mustard

import (
	"fmt"
	"image"
	"image/draw"
	assets "thdwb/assets"
	"thdwb/gg"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/goki/freetype/truetype"
)

//CreateTreeWidget - Creates and returns a new Tree Widget
func CreateTreeWidget() *TreeWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	return &TreeWidget{
		baseWidget: baseWidget{

			needsRepaint: true,
			widgets:      widgets,

			widgetType: treeWidget,

			cursor: glfw.CreateStandardCursor(glfw.ArrowCursor),

			backgroundColor: "#fff",

			font: font,
		},

		fontSize:  20,
		fontColor: "#000",
	}
}

//SetWidth - Sets the tree width
func (tree *TreeWidget) SetWidth(width int) {
	tree.box.width = width
	tree.fixedWidth = true
	tree.RequestReflow()
}

//SetHeight - Sets the tree height
func (tree *TreeWidget) SetHeight(height int) {
	tree.box.height = height
	tree.fixedHeight = true
	tree.RequestReflow()
}

//SetFontSize - Sets the tree font size
func (tree *TreeWidget) SetFontSize(fontSize float64) {
	tree.fontSize = fontSize
	tree.needsRepaint = true
}

//SetFontColor - Sets the tree font color
func (tree *TreeWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		tree.fontColor = fontColor
		tree.needsRepaint = true
	}
}

//SetBackgroundColor - Sets the tree background color
func (tree *TreeWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		tree.backgroundColor = backgroundColor
		tree.needsRepaint = true
	}
}

func (tree *TreeWidget) draw() {
	context := tree.window.context
	top, left, width, height := tree.computedBox.GetCoords()

	context.SetHexColor(tree.backgroundColor)
	context.DrawRectangle(float64(left), float64(top), float64(width), float64(height))
	context.Fill()

	for _, node := range tree.nodes {
		flowNode(context, node, tree, 0)
		drawNode(context, node, tree, 0)
	}

	if tree.buffer == nil || tree.buffer.Bounds().Max.X != width && tree.buffer.Bounds().Max.Y != height {
		tree.buffer = image.NewRGBA(image.Rectangle{
			image.Point{}, image.Point{width, height},
		})
	}

	draw.Draw(tree.buffer, image.Rectangle{
		image.Point{},
		image.Point{width, height},
	}, context.Image(), image.Point{left, top}, draw.Over)
}

func flowNode(context *gg.Context, node *TreeWidgetNode, tree *TreeWidget, level int) {
	node.box.left = level * int(tree.fontSize)
	node.box.height = int(tree.fontSize)

	prevSibling := node.PreviousSibling()
	if node.Parent == nil {
		if prevSibling != nil {
			node.box.top = prevSibling.box.top + prevSibling.box.height
		} else {
			node.box.top = 0
		}
	} else {
		if prevSibling != nil {
			node.box.top = prevSibling.box.top + prevSibling.box.height
		} else {
			fmt.Println("parent:", node.Parent.Content, node.Parent.box.top, node.Parent.box.height)
			node.box.top = node.Parent.box.top + node.Parent.box.height
			fmt.Println("node", node.Content, node.box.top, node.box.height)
		}
	}

	for _, childNode := range node.Children {
		flowNode(context, childNode, tree, level+1)
	}
}

func drawNode(context *gg.Context, node *TreeWidgetNode, tree *TreeWidget, level int) {
	top, left, _, _ := node.box.GetCoords()

	context.SetHexColor(tree.fontColor)
	context.SetFont(tree.font, tree.fontSize)
	context.DrawString("-> "+node.Content, float64(left)+tree.fontSize/4, float64(top)+tree.fontSize*2/2)
	context.Fill()

	for _, childNode := range node.Children {
		drawNode(context, childNode, tree, level+1)
	}
}
