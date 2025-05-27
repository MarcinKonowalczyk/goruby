_grgr_find_all = -> (range, fun) {
    out = []
    i = 0
    loop {
        break if i >= range.size
        if fun[i]
            out.push(i)
        end
        i += 1
    }
    out
}

# cannot just be a block since it caps str and chr
_block_capture_tE9aSVr0 = -> (str, chr) {
    -> (i) { str[i] == chr }
}

indices = -> (str, chr) {
    # _grgr_find_all[(0 ... str.length), -> (i) {str[i] == chr }]
    _grgr_find_all[(0 ... str.length), _block_capture_tE9aSVr0[str, chr]]
}

puts indices["hello", "l"] # => [2, 3]
puts _grgr_find_all[[1, 2, 3], -> (i) { i == 2 }]