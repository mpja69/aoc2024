# Day 22 - Monkey Market

## Iterators
Playing arouund with _Iterators_.

The outer function is a "wrapper", (it maybe can be called a "factory"??). This wrapper function returns a function, a _closure_, which is an on is an _Iterator_.

We have a text file, where each line has a number. 
The iterator reads line by line with `fmt.Fscanln`, and for each line produces an `integer`.of the number.


So we do **not** have to 
- read the whole file into a `[]byte`
- split the data into lines with `strings.Split(...)` into another `[]string`
- loop over the lines and convert every line/string to an `int` and store it in a `[]int`
...and then finally loop over the integer slice to perform the actual task at hand.

Even if this example doesn't have so many lines of data to be read, and either way of solving it is just `O(n)`, (where `n` is the length of the file). It makes the code _clearer_! It becomes easier to tread and understand the _intention_ of the different parts.

### Example code
Pretty nifty. Takes a file name and scans each line into an `int`, on the fly, when you need it.


``` 
func readNumbers(s string) func(func(int) bool) {
	f, err := os.Open(s)
	if err != nil {
		log.Fatal("readNumber(): ", err)
	}
	return func(yield func(int) bool) {
		defer f.Close()
		val := 0
		for _, err := fmt.Fscanln(f, &val); err != io.EOF; _, err = fmt.Fscanln(f, &val) {
			if err != nil {
				log.Fatal("readNumber(): ", err)
				return
			}
			if !yield(val) {
				return
			}
		}
	}

}

```

When solving the "actual problem", just use the Iterator with `range` in a loop:
```
for num := range readNumbers("myfile") { 
	... 
}
```

### Another way
The code inside the iterator could of course be used in the "actual problem loop" like this:

```
func main()
	f, err := os.Open("myFile")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	val := 0
	for _, err := fmt.Fscanln(f, &val); err != io.EOF; _, err = fmt.Fscanln(f, &val) {
		if err != nil {
			log.Fatal("readNumber(): ", err)
			return
			}
	
		// Do the actual task at hand with the value
		...
```
_But this makes to "overall code" unclear, since it's a mix of things solving completely different tasks!_ 

I'm not a "Clean Code" orthodox, but I like the _idea_ of "Single responsibilty Principle", and the Unix philosofy of "One tool one job".

I like to separate different tasks, with different intents, into different functions...which make sense to have this part as an Iterator.
(And of course the iterator approach would make even more sense if it were to be reused.)

At the same time I also like the idea of "Locoality of Code". <!-- TODO: Primagen? --> 

And I think of all these ideas as **LAIR**

## LAIR - Level of Abstraction, Intention and/or Responsibilty

_What problem are we trying to solve?_

Ask yourself that over and over, in different parts of the code.
And split/organize the code accordingly.

The goal is to get chunks of code that belongs together, and make sense to the "reader". (The "reader" is also yourself a week/month/year later.)

Aiming for clarity and simplicity.



>**Example:**
>
>- The `main()` is responsible for starting everything and running the 2 parts of the day, (in advent of code). Maybe it's responsible for getting the data (if it's the same for both parts).
>
>- The `Part1()` is responsible for solving "the high level problem" in part 1 and returning the answer.
>
>Part 1 maybe needs other supporting functions, e.g. iterators, calculators, etc., or it could be data structures, structs with methods. 
>And each of these support functionality should also be crafted to solve it's specific problem.
>- `readNumbers(string)`
>- `evolve()`
>- `RingBuffer`, with `.Write(int)` and `.Sequence()[4]int`
>- `seen map[...]bool`
 

What it's not:
- It's **not** about writing as little as possible, (in each function).
- It's **not** about making everything super abstract.
- It's **not** about that every support structure needs to be super general, and to be reused in other places.


## Process of trying/building stuff, high and low
Finding the ritght Level of Abstraction and organizing the code, is _hard to think out beforehand_!

The process of developing is **not** about trying to come up with the "correct" solution before you start typing code. 
It's almost impossible to just "think" about the solution...and get correct.

A good approach to the actual process is to jump betwwen different levels.
- Draw out the highlevel concept. _What problem are you trying to solve?_ 
	- What is the result we are after? 
	- What do we have as input?
	- What steps needs to be in between?
- Just start coding.
- Learn as you go.
- Adapt to what you discover.


