# Day 22 - Monkey Market

## Iterators
Playing arouund with _Iterators_.

The outer function is a "wrapper", (maybe can be called a "factory"). This wrapper function returns a function, a _closure_.
The returned function is an _Iterator_.

We have a text file, where each line has a number. 
The iterator reads line by line with `fmt.Fscanln`, and for each line produces an `integer`.of the number.


So we don't have to 
- read the whole file into a `[]byte`
- split the data into lines with `strings.Split(...)` into another `[]string`
- loop over the lines and convert every line to `int` and store it in a `[]int`
...and then loop over the integer slice to perform the actual task at hand.


``` 
func readNumber(s string) func(func(int) bool) {
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

### Another way
The code inside the iterator could of course be used "on the fly" like this:

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

#### LAIR - Level of Abstraction, Intention and/or Responsibilty

_What problem are we trying to solve?_


>Example:
The `main()` is responsible for starting everything and running the 2 parts of the day, (in advent of code). Maybe it's responsible for getting the data (if it's the sam for both parts).
>
>The `Part1()` is responsible for solving "the high level problem" in part 1 and returning the answer.
>
>Part 1 maybe needs other supporting functions, e.g. iterators, calculators, etc., or it could be data structures, structs with methods. 
>And each of these support functionality should also be crafted to solve it's specific problem.
>
>It's not about writing as little as possible.
>It's not about following "Clean Code" as Bible.
>It's not about that every other support structure needs to super general, and to be reused in other places.


#### Process of trying/building stuff, high and low

