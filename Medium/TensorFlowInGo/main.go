package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
)

func main() {
	// // Construct a graph with an operation that produces a string constant.
	// s := op.NewScope()
	// c := op.Const(s, "Hello from TensorFlow version "+tf.Version())
	// graph, err := s.Finalize()
	// if err != nil {
	// 	panic(err)
	// }

	// // Execute the graph in a session.
	// sess, err := tf.NewSession(graph, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// output, err := sess.Run(nil, []tf.Output{c}, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(output[0].Value())


	// Convert a image to a tensor
	filename := "6953297_8576bf4ea3.jpg"
	tensor, err := makeTensorFromImage(filename, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Shape:",tensor.Shape())
	
}

// Convert the image in filename to a Tensor suitable as input to the Inception model.
func makeTensorFromImage(filename string, url bool) (*tf.Tensor, error) {

	var bytes []byte
	if url {
		res, err := http.Get(filename)
		if err != nil {
			return nil, errors.New("invalid image url")
		}
		defer res.Body.Close()
		bytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.New("invalid image url")
		}
	} else {
		var err error
		bytes, err = os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	// Construct a graph to normalize the image
	graph, input, output, err := constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}
	// Execute that graph to normalize this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
func constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {
	// Some constants specific to the pre-trained model at:
	// https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip
	//
	// - The model was trained after with images scaled to 224x224 pixels.
	// - The colors, represented as R, G, B in 1-byte each were converted to
	//   float using (value - Mean)/Scale.
	const (
		H, W  = 224, 224
		Mean  = float32(117)
		Scale = float32(1)
	)
	// - input is a String-Tensor, where the string the JPEG-encoded image.
	// - The inception model takes a 4D tensor of shape
	//   [BatchSize, Height, Width, Colors=3], where each pixel is
	//   represented as a triplet of floats
	// - Apply normalization on each pixel and use ExpandDims to make
	//   this single image be a "batch" of size 1 for ResizeBilinear.
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s,
						op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)), tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}
