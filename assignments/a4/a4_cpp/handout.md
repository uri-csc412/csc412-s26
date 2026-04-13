

#### Makefile to make your life easier

Use our included Makefile to help your workflow; after downloading the
stencil run **make remake** to see the basic 8x8 grid printed to the terminal
with one red cell with an arrow. The stencil will not generate new
frames or updates to the grid. You will need to implement a
multi-threaded application in this assignment to get this animation. :)

<img src="/web/ascii.png" width="400" alt="ASCII Terminal Art. This figure shows an example of an 8x8 grid with 4 threads in multiple colors. It is replaced at set intervals by ascii_art.cpp."/>

> ASCII Terminal Art. This figure shows an example of an 8x8 grid with 4 threads in multiple colors. It is replaced at set intervals by ascii_art.cpp.

## Part 1: Understanding the ASCII Art Code

This program uses a library we wrote to write output to the terminal in the form of ASCII art. This is a larger codebase, and it is more complex than previous assignments. Some sections of this code must not be touched. Most of the work will be done in the main.cpp file, and some small changes will and can be made to the ```ascii_art.cpp``` file.

> Part 1 is all about reading and understanding the codebase.

### ASCII Art in the Terminal

When the program runs you will see a colored grid in the terminal. This grid is going to show the colored positions of multi-threaded “inklings” The grid shows the state of the simulation: the number of currently “live” inkling threads (i.e., of threads that have not terminated after having reached a grid corner) and the current levels of three ink tanks. 

The following keyboard commands are enabled if you uncomment ```std::thread listenerThread(keyListener);``` in the stencil ```ascii_art.cpp```:

* **esc** - terminates the program.
* **’r’** - partially refills the red ink tank.
* **’g’** - partially refills the green ink tank.
* **’b’** - partially refills the blue ink tank.

:::tip[HINT: how does your program behave when uncommenting this line?]
Does it crash? Well that is expected! Getting this to work is worth minimal points in the autograder but it could be an area you choose to propose a fix and implement.
:::


## Part 2: Creating Inklings as Multiple Threads

Part 2 focuses on implementing multiple "inkling" threads that move randomly (with some constraints) on a grid in the terminal. See the ```TODO``` comments in the ```main.cpp``` for hints. Inklings will need to do the following:

* Be one of three kinds of inklings - red, green, or blue.
* Have an initial unique starting position.
* Generate a random amount of valid spaces to travel.
* Record their position of travel on the grid.
* Have access to 1 ink resource before being able to move.
* Stop moving if there are no more available ink resources.
* If the ink has run out, refill the tank.
* Terminate the inkling thread if they reach a corner.

> Part 2's sections try to give some context and descriptions about the code. You will see many of these descriptions represented as tests within Gradescope.

:::caution[Submit to Gradescope Early]
Be sure to check the Gradescope tests early. You will need to make specific considerations early to pass each Gradescope test. For example, can the inklings run long enough to terminate, or does your program error out?
:::

### Grid Dimensions and Number of Threads Can be Passed as Arguments

The program supports arguments for the grid width, height, and the number of inkling threads, which are arguments of the program rather than hard-coded constants. If no arguments are passed, the program will create an 8x8 grid and set the program to attempt to create 4 inkling threads.

:::tip[HINT: You will need to create these threads!]
:::

### Inkling Movement

* Inklings can move vertically or horizontally on the grid. This means that in addition to a position (row and column index) on the grid, each inkling has a current position (i.e., orientation).
* At each iteration of its main loop, an inkling thread randomly selects a new displacement direction and displacement length. The new direction of displacement must be perpendicular to the current direction. The length of displacement must keep the inkling within the dimensions of the grid.
* An inkling that reaches one of the corners of the grid terminates its execution. If you want to join that thread in the main thread, for whatever reason, you can only do that after the ```myEventLoop()``` call towards the end of the main function in **main.cpp**. This is because this event loop handles all events and interrupts. If you insert a ```pthread_join``` (or, ```join()```) call (basically, a blocking call) anywhere, you basically block the event loop from updating the terminal. Keyboard events won’t be handled anymore, and no more rendering will occur. There is no way around that, as myEventLoop() should run on the main thread.

### Work with Colors

* An inkling going through a cell in the grid will paint its designated color in its current cell.
* They will write their color and direction (as an arrow) in their assigned color in the grid.
* The grid is a 2D array of int. You can access the grid cell using ```grid[i][j]```.

### Inkling Threads

You must be able to create at least eight inkling threads and assign them a random color among the three available, an initial position on the grid, and an initial orientation. Use the data types defined in ascii_art.h: InklingType for the choice of ink color and TravelDirection for the orientation.

Two inklings may occupy the same grid cell. In other words, only the color information of a grid cell is considered a shared resource under race conditions. It is likely the “space” within the grid cell is not under race conditions for this simulation. How might you handle this for a possible area to improve the program?

For the InklingInfo struct, you may consider taking it as a starting point for the data that you want to store and update about an inkling thread. The isLive boolean indicates whether the inkling is still running or has already found a grid corner. An inkling is only displayed on the grid if they are still alive (isLive).

## Part 3: Create Log Files for the Inklings

### Inkling Log Files

Each inkling will have its own log file that tracks its actions (condition and movement). Your program must create a logFolder directory with 0755 permissions, where your program will create each file. The log files must be titled numerically, inkling1.txt, inkling2.txt, etc., with 0755 permissions, and will contain the following information:

* The inkling’s color and starting position.
* The inkling’s direction of movement and new position.
* The inkling’s termination.

Each of these lines must also include the time at which every action
occurred. The log files must be in the following format:  

```
21:20:30.562,inkling3,green,row3,col2  
21:20:30.611,inkling3,west,row2,col2  
21:20:31.673,inkling3,north,row2,col1  
21:20:31.716,inkling3,east,row3,col1  
21:20:30.824,inkling3,west,row1,col1  
21:20:33.901,inkling3,terminated  
```

The top-left (north-west) corner is row0 and col0. The timestamps are
formatted as:  
```
hour:minute:second.milliseconds  
21:20:31.673
```

### Combine Log Files

You will create a file called `logs.cpp` that when run with:  
```make logs```
Will read all log files from logFolder, combine them, sort all actions by time, and save the new file in logFolder with the name actions.txt with ```0755``` permissions. 
> ```make clean``` will remove the logs executable.