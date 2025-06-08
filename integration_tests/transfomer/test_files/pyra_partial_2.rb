
def indices(str, chr)
    (0 ... str.length).find_all { |i| str[i] == chr }
end

print indices("hello", "l") # => [2, 3]

def unwrap(t)
    t.size == 1 ? t[0] : t
end

print unwrap([1, 2, 3]) # => [1, 2, 3]

$TOP    = "^"
$BOTTOM = "-"
$L_SIDE = "/"
$R_SIDE = "\\"

def triangle_from(lines, ptr_inds = nil)
    raise "no triangle found" if !lines.first
    ptr_inds = ptr_inds || indices(lines.first, $TOP)
    row = ""
end

print triangle_from([ # hello
    " ^ ",
    "/z\\",
    "---",
])
