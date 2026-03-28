package utils

// Mainly for declaring function that could help,
// this one in particular for converting string to
// its pointer counterpart which is easy to do
// in its main function but i choose to make it
// separate function cuz it looks more clean this way
// lol...

func Stroptr(msg string) *string {
	result := msg
	return &result
}
