/*
Package teq is a Go library designed to enhance your testing experience by providing a flexible and customizable way to perform deep equality checks. It's especially useful when you need to compare complex objects or types that have specific equality conditions.

# Features

- Transforms: Register a "transform" function that can modify objects before comparison. This allows you to control how equality is determined. For example, by transforming time.Time objects to UTC, you can make your equality checks timezone-insensitive.
- Formats: Register a "format" function that defines how objects are displayed when they are not equal. This is useful for types like time.Time and decimal.Decimal that may not be human-readable in their default format. By registering your own format, you can make the output of your tests more understandable.

# Usage

To use teq, you first need to create a new instance:

	tq := teq.New()

Then, you can add your transforms and formats:

	// time.Time will be transformed into UTC time. So equality check with `tq` will be timezone-insensitive.
	tq.AddTransform(func(d time.Time) time.Time {
	    return d.In(d.UTC)
	})

	// time.Time will be shown in RFC3339 format when it appear in diff.
	tq.AddFormat(func(d time.Time) string {
	    return d.Format(time.RFC3339)
	})

Finally, you can use teq to perform deep equality checks in your tests:

	tq.Equal(t, expected, actual)
*/
package teq
