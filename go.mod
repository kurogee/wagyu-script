module example.com/wagyu_script

go 1.22.5

replace example.com/wagyu_script_date => ./date

replace example.com/wagyu_script_get => ./get

replace example.com/wagyu_script_file => ./file

replace example.com/wagyu_script_array => ./array

replace example.com/wagyu_script_string => ./string

replace example.com/wagyu_script_math_sharp_functions => ./maths

require (
	example.com/wagyu_script_array v0.0.0-00010101000000-000000000000
	example.com/wagyu_script_date v0.0.0-00010101000000-000000000000
	example.com/wagyu_script_file v0.0.0-00010101000000-000000000000
	example.com/wagyu_script_get v0.0.0-00010101000000-000000000000
	example.com/wagyu_script_string v0.0.0-00010101000000-000000000000
	example.com/wagyu_script_math_sharp_functions v0.0.0-00010101000000-000000000000

	github.com/Knetic/govaluate v3.0.0+incompatible
)
