# WALK *ExpressionStatement::Expression: ast.Walk mutated *ast.ContextCallExpression(<<<ContextCallExpression>>>) to *ast.IndexExpression(<<<IndexExpression>>>)
# TRANSFORM_PIPELINE *ast.Assignment
__builtin_find_all = -> (range, fun) {out = ([]); i = 0; loop {break if i >= range.size; if fun[i]; out.push(i) end; i = (i + 1)}; out}

# TRANSFORM_PIPELINE *ast.Assignment
__block_Vgm6Ywsp = -> (str, chr) {-> (i) {str[i] == chr}}

# TRANSFORM_PIPELINE *ast.Assignment
__lifted_indices = -> (str, chr) {__builtin_find_all[(0 ... str.length), __block_Vgm6Ywsp[str, chr]]}

# TRANSFORM_PIPELINE *ast.Assignment
__lifted_unwrap = -> (t) {if t.size == 1; t[0]else t end}

# TRANSFORM_PIPELINE *ast.Assignment
__lifted_triangle_from = -> (lines, ptr_inds = :nil) {if !lines.first; raise("no triangle found") end; ptr_inds = (ptr_inds || __lifted_indices[lines.first, $TOP]); row = ""}

# TRANSFORM_PIPELINE *ast.FunctionLiteral
-> (str, chr) {__builtin_find_all[(0 ... str.length), __block_Vgm6Ywsp[str, chr]]}

# TRANSFORM_PIPELINE *ast.ContextCallExpression
print(__lifted_indices["hello", "l"])

# => [2, 3]
# TRANSFORM_PIPELINE *ast.FunctionLiteral
-> (t) {if t.size == 1; t[0]else t end}

# TRANSFORM_PIPELINE *ast.ContextCallExpression
print(__lifted_unwrap[([1, 2, 3])])

# => [1, 2, 3]
# TRANSFORM_PIPELINE *ast.Assignment
$TOP = "^"

# TRANSFORM_PIPELINE *ast.Assignment
$BOTTOM = "-"

# TRANSFORM_PIPELINE *ast.Assignment
$L_SIDE = "/"

# TRANSFORM_PIPELINE *ast.Assignment
$R_SIDE = "\\"

# TRANSFORM_PIPELINE *ast.FunctionLiteral
-> (lines, ptr_inds = :nil) {if !lines.first; raise("no triangle found") end; ptr_inds = (ptr_inds || __lifted_indices[lines.first, $TOP]); row = ""}

