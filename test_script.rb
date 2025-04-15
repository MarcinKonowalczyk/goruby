a = [1, 2, 3]
b1 = [*a]
b2 = [99, *a]
b3 = [*a, 99]
puts("b1:", b1)
puts("b2:", b2)
puts("b3:", b3)

a = [1, 2, 3]
b1 = [*a]


# # puts("zoo:", zoo["?"][3])
# $zoo = {
#     "?" => [0,1,2,3,4,5],
#     "/" => -> (a, b) { 1.0 * a / b }
# }
# def call_zoo(args)
#     return $zoo["/"][*args]
# end
# args = [355, 113]
# puts("zoo:", $zoo["/"][*args])
# puts("zoo:", call_zoo(args))

# def foo()
#     return 0
# end

# a = *foo()
# puts(a, *foo())
# puts("$LOADED_FEATURES", $LOADED_FEATURES)
# puts("ARGV.size", ARGV.size, ARGV)

# this_file = __FILE__
# content = File.read(this_file)
# puts("this_file:", this_file)
# puts("content:", content)

x = 0
# while x < 20
#     x += 1
#     puts("x:", x)
#     break if x == 10
#     puts("past break")
# end
loop {
    x += 1
    puts("x:", x)
    break unless x < 10
    puts("past break")
    raise "my little exception" if false or x == 5
}