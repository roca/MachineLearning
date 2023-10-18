# TensorFlow in Go

- [Link](https://medium.com/@exyzzy/tensorflow-in-go-9387fc130f50)

- [Build an Image Recognition API with Go and TensorFlow](https://outcrawl.com/image-recognition-api-go-tensorflow)

 ```sh
export LIBRARY_PATH=$LIBRARY_PATH:$HOME/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$HOME/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export CPLUS_INCLUDE_PATH=$CPLUS_INCLUDE_PATH:$HOME/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/include
export C_INCLUDE_PATH=$C_INCLUDE_PATH:$HOME/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/include
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$HOME/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export CGO_ENABLED="1"
gcc hello_tf.c -ltensorflow -o hello_tf
 ```

## Notes

- [Source code](https://github.com/exyzzy/tenseimage)
