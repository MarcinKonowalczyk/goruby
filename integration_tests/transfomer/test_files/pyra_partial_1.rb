def indices(str, chr)
    (0 ... str.length).find_all { |i| str[i] == chr }
end

print indices("hello", "l") # => [2, 3]