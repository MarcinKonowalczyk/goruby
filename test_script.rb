# a = [1, 2, 3]
# b1 = [*a]
# b2 = [99, *a]
# b3 = [*a, 99]
# puts("b1:", b1)
# puts("b2:", b2)
# puts("b3:", b3)

# a = [1, 2, 3]
# b1 = [*a]


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

# x = 0
# while x < 20
#     x += 1
#     puts("x:", x)
#     break if x == 10
#     puts("past break")
# end
# loop {
#     x += 1
#     puts("x:", x)
#     break unless x < 10
#     puts("past break")
#     # raise "my little exception" if false or x == 5
# }

# $R_SIDE = "\\"
# puts("R_SIDE:", $R_SIDE, $R_SIDE.length) 

# $L_SIDE = "/"

# if "/" != $L_SIDE
#     puts("not equal")
# else
#     puts("equal")
# end

# puts([0, 1, 2, 3, 4, 5][2, 3])
# puts()
# puts([0, 1, 2, 3, 4, 5][2..-1])
# puts("hello"[2, 3])
# puts("hello"[2..3])
# puts("hello"[3..-1])
# puts([1, 2, 3].is_a?(Array))
# puts([1, 2, 3].is_a?(String))

# puts([0,1,2,3,4,5,6,7].all? { |e| e < 8 })
# puts([0,1,2,3,4,5,6,7].join(","))
# puts(*(0..3))
# puts(*(1..-1))

# raise "bzzzz"


# $zoo = {
#     "?" => [0,1,2,3,4,5],
#     "/" => -> (a, b) { 1.0 * a / b }
# }

# op = []
# puts("undefined operation `#{op}`")

# chain = ['set', 'v', 312312]
# op, args = chain
# puts("> chain", chain, "> op", op, "> args", args)

# if false then
#     puts("1")
# elsif false then
#     puts("2")
# elsif true then
#     puts("3")
# else
#     puts("4")
# end

# foo = [2, 16.0]
# bar = -> (a, b) {a ** b}
# puts(bar[*foo])
# puts("".to_f)

# $UNDEF = :UNDEF

# def falsey(val)
#     [0, [], "", $UNDEF, "\x00", nil, 123].include?(val)
# end

# # puts(0 == false)
# # puts(0.0 == false)
# # puts([] == false)
# # puts ("" == false)
# # puts($UNDEF == false)
# # puts("\x00" == false)
# # puts(nil == false)
# # puts(123 == false)
# # puts(false == 0)
# # puts(false == 0.0)
# # puts(false == [])
# # puts(false == "")
# # puts(false == $UNDEF)
# # puts(false == "\x00")
# # puts(false == nil)
# # puts(false == 123)
# # puts([0, [], "", $UNDEF, "\x00", nil, 123].include?(false))

# # puts(falsey(false))
# # puts(falsey(true))
# # puts(falsey(0))
# # puts(falsey([]))
# # puts(falsey(""))
# # puts(falsey($UNDEF))
# # puts(falsey("\x00"))
# puts(falsey(nil))
# puts(falsey(1))
# puts(falsey(1.0))
# puts(falsey([1, 2, 3]))
# puts(falsey("hello"))
# puts(falsey(0.0))
# puts(falsey(0.1))
# puts(falsey(123))
# puts(falsey(123.0))

# puts(0.0 == 0)

# def unwrap(t)
#     t.size == 1 ? t[0] : t
# end

# $zoo = {
#     "" => -> (*a) { unwrap(a) }
# }

# arg = [$UNDEF, $UNDEF]
# # puts($zoo[""][*arg])
# puts(unwrap(arg))
# puts($zoo[""][*arg])
# puts($UNDEF.size)

# b = arg.each { |e| print(1) }
# puts(b)

# a = [1, 2, 3]
# puts(a.pop)
# puts(a.pop)
# puts(a.pop)
# puts(a.pop)
# puts(nil)
# print(nil)

# plop = nil

# if plop && plop[1] == "d"
#     puts("plop")
# else
#     puts("no plop")
# end

# puts([])
# print([])
# puts()
# # puts([], 1)
# # print([], 1)
# # puts(*[])
# # print(*[])

# print([1, 2, 3] - [])
# puts()
# print([1, 2, 3] - [1])
# puts()
# print([1, 2, 3] - [2])
# puts()
# print([1, 2, 3] - [3])
# puts()
# print([1, 2, 3] - [4])
# puts()
# print([1, 2, 3] - [1, 2])


# print([1, 2, 3] + [])
# puts()
# print([1, 2, 3] + [1])
# puts()
# print([1, 2, 3] + [2])
# puts()
# print([1, 2, 3] + [3])
# puts()
# print([1, 2, 3] + [4])
# puts()
# print([1, 2, 3] + [1, 2])
# puts()

# print([1, 2, 3] * 1)
# puts()
# print([1, 2, 3] * 2)
# puts()
# print([1, 2, 3] * 3.999)
# puts()
# print([1, 2, 3] * ',')
# puts()

# puts("hello", 1)
# print("hello", 1)
# puts(["hello", 1])
# print(["hello", 1])
# puts()

# puts(0.1+0.2)

# Still not working:
# - favourite_number
# - golf_pyramid_scheme_negation
# - xor

# def foo(a)
#     return a
# end
# def bar(a)
#     return foo a
# end

# puts foo 'a'
# puts bar('a')
# puts bar 'a'
puts 'a'
puts [1,2,3]

if true || false 
    puts("true")
else
    puts("false")
end

a = (b = 99)

puts("a:", a)
puts("b:", b)


$ops = {
    "+" => -> (a, b) { a + b },
    "*" => -> (a, b) { a * b },
    "-" => -> (a, b) { a - b },
    "/" => -> (a, b) { 1.0 * a / b },
    "^" => -> (a, b) { a ** b },
    "=" => -> (a, b) { (a == b).to_i },
    "<=>" => -> (a, b) { a <=> b },
    "out" => -> (*a) { $outted = true; a.each { |e| print e }; },
    "chr" => -> (a) { a.to_i.chr },
    "arg" => -> (*a) { a.size == 1 ? ARGV[a[0]] : a[0][a[1]] },
    "#" => -> (a) { str_to_val a },
    "\"" => -> (a) { val_to_str a },
    "" => -> (*a) { unwrap a },
    "!" => -> (a) { falsey(a).to_i },
    "[" => -> (a, b) { a },
    "]" => -> (a, b) { b },
}

args = [1, 2]
print $ops["+"][*args]
print $ops["out"][*args]