//
//  main.cpp
//  inklings
//

 /*-------------------------------------------------------------------------+
 |	A graphic front end for a grid+state simulation.						|
 |																			|
 |	This application simply creates a a colored grid and displays           |
 |  some state information in the terminal using ASCII art.			        |
 |	Only mess with this after everything else works and making a backup		|
 |	copy of your project.                                                   |
 |																			|
 |	Current Keyboard Events                                     			|
 |		- 'ESC' --> exit the application									|
 |		- 'r' --> add red ink												|
 |		- 'g' --> add green ink												|
 |		- 'b' --> add blue ink												|
 +-------------------------------------------------------------------------*/

#include <random>
#include <vector>
#include <cstdlib>
#include <ctime>
#include <iostream>
#include <thread>
#include <unistd.h>
#include <mutex>

#include "ascii_art.h"

//==================================================================================
//	Function prototypes
//==================================================================================
void displayGridPane(void);
void displayStatePane(void);
void initializeApplication(void);
void threadFunction(InklingInfo* inkling);
void getNewDirection(InklingInfo* inkling);
bool checkIfInCorner(InklingInfo* inkling);
void redColorThreadFunc();
void greenColorThreadFunc();
void blueColorThreadFunc();
bool checkEnoughInk(InklingInfo* inkling, int moveAmount);

//==================================================================================
//	Application-level global variables
//==================================================================================

//	The state grid and its dimensions
int** grid;
int NUM_ROWS, NUM_COLS;

//	the number of live threads (that haven't terminated yet)
int MAX_NUM_TRAVELER_THREADS;
int numLiveThreads = 0;

//vector to store each struct
std::vector<InklingInfo> info;
bool DRAW_COLORED_TRAVELER_HEADS = true;

//	the ink levels
int MAX_LEVEL = 50;
int MAX_ADD_INK = 10;
int REFILL_INK = 10;
int redLevel = 20, greenLevel = 10, blueLevel = 40;

// create locks for color levels
std::mutex redLock;
std::mutex blueLock;
std::mutex greenLock;
std::mutex blueCellLock;
std::mutex redCellLock;
std::mutex greenCellLock;

// ink producer sleep time (in microseconds)
// [min sleep time is arbitrary]
const int MIN_SLEEP_TIME = 30000; // 30000
int producerSleepTime = 100000; // 100000

// inkling sleep time (in microseconds)
int inklingSleepTime = 1000000; // 1000000


//==================================================================================
//	These are the functions that tie the simulation with the rendering.
//	Some parts are "don't touch."  Other parts need your help to ensure
//	that access to critical data and the ASCII art are properly synchronized
//==================================================================================

void displayGridPane(void) {
	//---------------------------------------------------------
	//	This is the call that writes ASCII art to render the grid.
	//
	//	Should we synchronize this call?
	//---------------------------------------------------------
    drawGridAndInklingsASCII(grid, NUM_ROWS, NUM_COLS, info);
}

void displayStatePane(void) {
	//---------------------------------------------------------
	//	This is the call that updates state information
	//
	//	Should we synchronize this call?
	//---------------------------------------------------------
	drawState(numLiveThreads, redLevel, greenLevel, blueLevel);
}

//------------------------------------------------------------------------
//	These are the functions that would be called by a inkling thread in
//	order to acquire red/green/blue ink to trace its trail.
//	You *must* synchronize access to the ink levels (C++ lock and unlock)
//------------------------------------------------------------------------
// You probably want to edit these...
bool acquireRedInk(int theRed) {
	bool ok = false;
	if (redLevel >= theRed)
	{
		redLevel -= theRed;
		ok = true;
	}
	return ok;
}

bool acquireGreenInk(int theGreen) {
	bool ok = false;
	if (greenLevel >= theGreen)
	{
		greenLevel -= theGreen;
		ok = true;
	}
	return ok;
}

bool acquireBlueInk(int theBlue) {
	bool ok = false;
	if (blueLevel >= theBlue)
	{
		blueLevel -= theBlue;
		ok = true;
	}
	return ok;
}


//------------------------------------------------------------------------
//	These are the functions that would be called by a producer thread in
//	order to refill the red/green/blue ink tanks.
//	You *must* synchronize access to the ink levels (C++ lock and unlock)
//------------------------------------------------------------------------
// You probably want to edit these...
bool refillRedInk(int theRed) {
	bool ok = false;
	if (redLevel + theRed <= MAX_LEVEL)
	{
		redLevel += theRed;
		ok = true;
	}
	return ok;
}

bool refillGreenInk(int theGreen) {
	bool ok = false;
	if (greenLevel + theGreen <= MAX_LEVEL)
	{
		greenLevel += theGreen;
		ok = true;
	}
	return ok;
}

bool refillBlueInk(int theBlue) {
	bool ok = false;
	if (blueLevel + theBlue <= MAX_LEVEL)
	{
		blueLevel += theBlue;
		ok = true;
	}
	return ok;
}

//------------------------------------------------------------------------
//	You shouldn't have to touch this one.  Definitely if you do not
//	add the "producer" threads, and probably not even if you do.
//------------------------------------------------------------------------
void speedupProducers(void) {
	// decrease sleep time by 20%, but don't get too small
	int newSleepTime = (8 * producerSleepTime) / 10;
	
	if (newSleepTime > MIN_SLEEP_TIME) {
		producerSleepTime = newSleepTime;
	}
}

void slowdownProducers(void) {
	// increase sleep time by 20%
	producerSleepTime = (12 * producerSleepTime) / 10;
}

//-------------------------------------------------------------------------------------
//	You need to change the TODOS in the main function to pass the the autograder tests
//-------------------------------------------------------------------------------------
int main(int argc, char** argv) {
    // a try/catch block for debugging to catch weird errors in your code
    try {
        // check that arguments are valid, must be a 20x20 or greater and at least 8 threads/inklings
        if (argc == 4) {
            if (std::stoi(argv[1]) >= 20 && std::stoi(argv[2]) >= 20 && std::stoi(argv[3]) >= 8) {
                NUM_ROWS = std::stoi(argv[1]);
                NUM_COLS = std::stoi(argv[2]);
                MAX_NUM_TRAVELER_THREADS = std::stoi(argv[3]);
                numLiveThreads = std::stoi(argv[3]);
            }
        } else {
            std::cout << "No arguments provided, running with 8x8 grid and 4 threads.\n\tThis message will dissapear in 2 seconds... \n";
            sleep(2); // so the user can read the message in std::cout one line up
            clearTerminal();
            // some small defaults, will these run?
            NUM_ROWS = 8;
            NUM_COLS = 8;
            MAX_NUM_TRAVELER_THREADS = 4;
            numLiveThreads = 4;
        }
        
        initializeFrontEnd(argc, argv, displayGridPane, displayStatePane);
        
        initializeApplication();

        // TODO: create producer threads that check the levels of each ink

        
        // TODO: create threads for the inklings

        
        // now we enter the main event loop of the program
        myEventLoop(0);

        // ensure main does not return immediately, killing detached threads
        std::this_thread::sleep_for(std::chrono::seconds(30));
        
    } catch (const std::exception& e) {
        std::cerr << "ERROR :: Oh snap! unhandled exception: " << e.what() << std::endl;
    } catch (...) {
        std::cerr << "ERROR :: Red handed! unknown exception caught" << std::endl;
    }

	return 0;
}


//==================================================================================
//
//	TODO this is a part that you have to edit and add to.
//
//==================================================================================

void cleanupAndQuit(const std::string& msg) {
    std::cout << "Somebody called quits, goodbye sweet digital world, this was their message: \n" << msg;
	// should we join all the threads before you free the grid and other allocated data structures.  
    // you may run into seg-fault and other ugly termination issues otherwise.
	
	// also, if you crash there, you know something is wrong in your code.
	for (int i=0; i< NUM_ROWS; i++)
		delete []grid[i];
	delete []grid;

	// clear the inkling list
    exit(0);
}

void initializeApplication(void) {
	//	Allocate the grid
	grid = new int*[NUM_ROWS];
	for (int i=0; i<NUM_ROWS; i++)
		grid[i] = new int[NUM_COLS];
	
	//---------------------------------------------------------------
	//	The code block below to be replaced/removed
	//	I initialize the grid's pixels to have something to look at
	//---------------------------------------------------------------
	//	Random generator for the initial colors
	//	The initial cell color generation code should go away, but you
	//	should use tese for position and inkling type generation.
	std::random_device myRandDev;
	std::default_random_engine myEngine(myRandDev());
	
	for (int i=0; i<NUM_ROWS; i++) {
		for (int j=0; j<NUM_COLS; j++) {
			// the intialization you should use
            grid[i][j] = 0;
		}
	}

	//---------------------------------------------------------------
	//	TODO you are going to have to properly initialize your  
	//	inklings at random locations on the grid:
	//		- not at a corner
	//		- not at the same location as an existing inkling
	//---------------------------------------------------------------

    // TODO random number generation initialization

    // TODO maybe some kind of loop here?

    // TODO here is some starter code so everything does not crash
    // we prob won't need this exactly in the final solution
    InklingInfo inked = {	RED_TRAV,
							1, 5,
							EAST, true};
	info.push_back(inked); // aka the inklings
}

// TODO help me please, I am a vagabond
void threadFunction(InklingInfo* inkling) {
    
}

// TODO help me please, I have no direction
void getNewDirection(InklingInfo* inkling) {
    
}

// TODO help me please, I am blank and need some code
bool checkIfInCorner(InklingInfo* inkling) {

}

// TODO check if you have enough ink depending on what kind of inkling you are, maybe some locks/unlocks
bool checkEnoughInk(InklingInfo* inkling, int moveAmount) {

}

// TODO thread function for a red ink producer
void redColorThreadFunc() {

}

// TODO thread function for a green ink producer
void greenColorThreadFunc() {

}

// TODO thread function for a blue ink producer
void blueColorThreadFunc() {

}
