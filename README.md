# VoxelEngine

This is a small voxel engine built with golang.
In order to make the project cross platform it's
compiled to web assembly and then serverd with a
simple web server.

Note: This project makes heavy use of the syscall/js go library to interact with webGL.
This library is still under development and may introduce breaking changes which
could break this project. In order to run this project you may have to use an older
version of golan and syscall/js. This project last was tested on go version 1.14.2 on macOS

# Running The Project

First, you'll need to [install golang](https://golang.org/) if you have not already.
once go is installed navigate into the go folder. Then simply run ```./server/server```
from your terminal. This will start the server that will serve the dist directory. Note
that the server executable was built for macOS so if you are running this on another OS
you'll need to re-build the server. you can do this by simply running ```go build``` in
the go/server directory.

# Building The Project

You can build the project by running the build.sh file in the go/wasm directory. This will 
compile the go code into a wasm file and place it in the dist folder. You may also need to
update the wasm_exec.js file. This file must match the version of go used to compile the
wasm file. This file can be located in the go/misc/wasm/wasm_exec.js and can simply be
coppied from that directory into the dist directory.

# Saving & Loading Files

You can save and load the voxel project you are working on so you don't loose your project.
the files are stored in .vng format. You can see several example vng files in the example 
folder of this repo.

![Open File](https://github.com/bjatkin/VoxelEngine/blob/master/images/OpenFile.png)

# Controlls
 * Righ Click & drag - Zoom In
 * Left Click & drag - Select voxel faces
 * Left Click & drag + Left Shift - Pan the view
 * Middle Click & drag - Pan the view
 * Left Click & drag + Left Alt - Rotate the view
 * A Key - switch into add mode. Selected faces will turn red
 * S Key - switch into subtract mode. Selected faces will turn blue
 * Left Click on add mode faces - add voxels on top of the selected faces
 * Left Click on subtract mode faces - remove selected voxels and the voxels beneth them
 * ESC - deselect all voxel faces

### Add Mode
![Add Mode](https://github.com/bjatkin/VoxelEngine/blob/master/images/AddMode.png)

### Sub Mode
![Sub Mode](https://github.com/bjatkin/VoxelEngine/blob/master/images/SubtractMode.png)

# Examples
The following are example models created using this project

![Green Pipe](https://github.com/bjatkin/VoxelEngine/blob/master/images/GreenPipe.png)
![Gold Coin](https://github.com/bjatkin/VoxelEngine/blob/master/images/Coin.png)
![Mushroom](https://github.com/bjatkin/VoxelEngine/blob/master/images/Mushroom.png)
![Zelda Sword](https://github.com/bjatkin/VoxelEngine/blob/master/images/ZeldaSword.png)

# Future Project Goals
 * Create a standalone version using electron
 * Fix the strange camera rotation. (I may need to move the origin before rotating?)
 * Fix occasional nil pointer crashes when removing voxels
 * Optimize the voxel mesh so that I can handle more geometry 
    (this will probably need to happen asyncronously)
 * Improve the UI (Maybe it should all just exsist in the canvas?)
 * Save and load the color pallet
 * Could I save GPU and transfer speed by using indexed color?
 * Clean up the code. Code could always be cleaner.
 * Import/ export .vox files instead of custom .vng files
