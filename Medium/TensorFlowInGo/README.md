# TensorFlow in Go
 - [Link](https://medium.com/@exyzzy/tensorflow-in-go-9387fc130f50)

 ```sh
export LIBRARY_PATH=$LIBRARY_PATH:/Users/romel.campbell/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:/Users/romel.campbell/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export CPLUS_INCLUDE_PATH=$CPLUS_INCLUDE_PATH:/Users/romel.campbell/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/include
export C_INCLUDE_PATH=$C_INCLUDE_PATH:/Users/romel.campbell/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/include
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/Users/romel.campbell/GitHub/Repos/MachineLearning/Medium/TensorFlowInGo/lib/lib
export CGO_ENABLED="1"

gcc hello_tf.c -ltensorflow -o hello_t -I lib/include -L lib/lib



 ```
