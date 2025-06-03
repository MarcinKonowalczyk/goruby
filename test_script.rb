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
    # ptr_inds.map { |pt|
    #     x1 = x2 = pt    # left and right sides
    #     y = 0
    #     data = []
    #     loop {
    #         x1 -= 1
    #         x2 += 1
    #         y += 1
    #         row = lines[y]
    #         raise "unexpected triangle end" if !row or x2 > row.size
    #         # are we done?
    #         if row[x1] != $L_SIDE
    #             # look for end
    #             if row[x2] == $R_SIDE # mismatched!
    #                 raise "left side too short"
    #             else
    #                 # both sides are exhausted--look for a bottom
    #                 # p (x1 + 1 .. x2 - 1).map { |e| row[e] }
    #                 # p [x1, x2, pt]
    #                 if (x1 + 1 .. x2 - 1).all? { |x| row[x] == $BOTTOM }
    #                     break
    #                 else
    #                     raise "malformed bottom"
    #                 end
    #             end
    #         elsif row[x2] != $R_SIDE
    #             # look for end
    #             if row[x1] == $L_SIDE # mismatched!
    #                 raise "right side too short"
    #             else
    #                 # both sides are exhausted--look for a bottom
    #                 if (x1 + 1 .. x2 - 1).all? { |x| row[x] == $BOTTOM }
    #                     break
    #                 else
    #                     raise "malformed bottom"
    #                 end
    #             end
    #         # elsif x1 == 0   # we must have found our bottom...
    #         end
    #         #todo: do for other side
    #         # we aren't done.
    #         data.push row[x1 + 1 .. x2 - 1]
    #     }
    #     op = data.join("").gsub(/\s+/, "")
    #     args = []
    #     if row[x1] == $TOP or row[x2] == $TOP
    #         next_inds = [x1, x2].find_all { |x| row[x] == $TOP }
    #         args.push triangle_from(lines[y..-1], next_inds)
    #     end
    #     unwrap [op, *args]
    # }
end
